# 07: Task scheduler reference

## Concepts
- **Single-timer scheduler:** keep a min-heap by deadline. One goroutine waits on `select { timer.C, newTaskCh, ctx.Done }`. Reset the timer whenever the head changes. This uses O(1) goroutines instead of one goroutine or timer per task.
- **Timing wheels** (Varghese & Lauck): O(1) insert and cancel, with a bounded tick granularity. Used by Kafka, Netty, and the Linux kernel. Go's runtime uses per-P 4-ary heaps for timers.
- **Fixed-rate vs fixed-delay:** rate schedules from the planned start (it can catch up, or burst). Delay schedules from the end of the previous run.
- **DAG execution:**
  - Topological order via in-degrees.
  - Readiness is event-driven: when a task completes, decrement its children and enqueue the ones that reach zero.
  - Concurrency is bounded by a semaphore.
  - Cancellation is shared ctx + errgroup.
- **Priority + starvation:** strict priority starves low classes. Use aging, weighted fair queueing, or separate pools per class (bulkheads).
- **Idempotency and retries:** in a scheduler, retries assume the tasks are idempotent. If they aren't, you need dedupe keys (see `17`).

## Hidden knowledge / gotchas
- Never `time.AfterFunc` per task when you have millions. Each one is a runtime timer, and the callback runs in its own goroutine.
- Timer reset semantics changed in 1.23 (no stale values after Reset). Code written for older Go often has a drain-then-reset dance that's now unnecessary. Know both versions.
- Wall clock vs monotonic: `time.Time` from `time.Now()` carries a monotonic reading. Durations between two `Now()`s are immune to wall-clock jumps. Deadlines persisted to disk aren't.
- A DAG scheduler that holds its mutex while *launching* tasks can deadlock if a task completion callback needs that same mutex.

## Questions you must answer
1. Design a delayed job scheduler that handles 1M pending jobs. Heap or wheel? How do you cancel?
2. How do you schedule a DAG of tasks with max parallelism and fail-fast?
3. A periodic job takes longer than its period. What should happen?
4. How do you prevent starvation in a priority queue?
5. How do you make the scheduler survive a restart? (Persist, WAL, replay; see `12`.)

## Reading
- Varghese & Lauck, *Hashed and Hierarchical Timing Wheels* (1987)
- `$(go env GOROOT)/src/runtime/time.go` (per-P timer heaps)
- Java `ScheduledThreadPoolExecutor` docs (fixed-rate vs fixed-delay semantics)
- Kahn's algorithm (any algorithms text)

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**after `delayed/`**
- Code `[S]` [client-go `util/workqueue/delaying_queue.go`](https://github.com/kubernetes/client-go/blob/master/util/workqueue/delaying_queue.go): heap + a single timer goroutine (your `delayed/`).
- Code `[S]` [client-go `util/workqueue/rate_limiting_queue.go`](https://github.com/kubernetes/client-go/blob/master/util/workqueue/rate_limiting_queue.go): composes the delaying queue with rate limiters.

**after `priority/`**
- Code `[S]` [kubernetes `pkg/scheduler/backend/queue/scheduling_queue.go`](https://github.com/kubernetes/kubernetes/blob/master/pkg/scheduler/backend/queue/scheduling_queue.go): active heap + backoff + unschedulable pools, requeued on events.

**after `dag/`**
- Read `[T]` [Google: Borg (EuroSys 2015)](https://research.google/pubs/large-scale-cluster-management-at-google-with-borg/): overcommit + admission control: utilization vs interference.
- Code `[S]` [temporal `service/matching/`](https://github.com/temporalio/temporal/tree/main/service/matching): sync-match tasks to long-polling workers, falling back to a DB backlog.
- Code `[S]` [cadence `service/matching/tasklist/`](https://github.com/uber/cadence/tree/master/service/matching/tasklist): the same design before the Temporal fork; diff the two.
