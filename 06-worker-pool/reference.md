# 06: Worker pool reference

## Concepts
- **Why pool in Go at all?** Goroutines are cheap. You pool to **bound a resource** (CPU, DB connections, memory, downstream QPS), not to save goroutine creation. Often a semaphore around `go f()` is simpler and better.
- **Sizing:**
  - CPU-bound: about `GOMAXPROCS`.
  - I/O-bound: from Little's Law, `workers ≈ throughput × latency`, capped by downstream limits.
- **Rejection policies:** block (backpressure), reject (fail fast, the caller retries or sheds), caller-runs (natural throttling, but it blocks the submitter's goroutine), drop-oldest/newest (lossy).
- **Shutdown:**
  - Graceful: stop intake, finish queued work, then exit.
  - Immediate: stop intake, cancel running work, hand back queued work.
  - Poison pills need one per worker, and fail with multiple producers unless intake is stopped first.
  - Closing the intake channel only works if *nobody* can send afterwards, so guard Submit with a RWMutex or state flag.
- **Submit/Shutdown race:** the classic bug is `Submit` checking "not closed" then sending while `Shutdown` closes the channel (panic: send on closed channel). Fix:
  - hold `RLock` during Submit's send and `Lock` while closing, or
  - never close the work channel; signal via a separate `done` channel and `select`.
- **Panics:** an unrecovered panic in a worker crashes the process. Recover inside the worker loop, per task.
- **Futures/promises:** a result slot + a completion channel closed exactly once. `Get` selects on done vs ctx.

## Production equivalents
- Go: `errgroup.SetLimit`, `x/sync/semaphore`, `github.com/sourcegraph/conc` (pool, panics), `github.com/alitto/pond`.
- Java: `ThreadPoolExecutor` (core/max/queue/handler), `ForkJoinPool` (work stealing), `CompletableFuture`, virtual threads (`newVirtualThreadPerTaskExecutor`: don't pool them, bound with semaphores).

## Hidden knowledge / gotchas
- `ThreadPoolExecutor` grows beyond core size only when the queue is **full**. With an unbounded queue it never grows. Classic interview trivia, and a real outage cause.
- Passing `ctx` from Submit into the task: the submitter's request ctx may be cancelled before the task runs. Decide whether tasks get the pool's lifecycle ctx or the caller's, and document it.
- Work stealing (per-worker deques) beats a single shared queue at high task rates. The Go scheduler does exactly this for goroutines (see `15`).
- A task that blocks on another task in the same bounded pool can deadlock the pool (thread starvation deadlock).

## Questions you must answer
1. Implement a thread pool. Now add shutdown. Now add futures.
2. What happens when the queue is full? Give 4 policies and when to use each.
3. How do you avoid the send-on-closed-channel panic during shutdown?
4. Is a worker pool even necessary in Go? Argue both sides.
5. How do you size the pool for CPU-bound vs I/O-bound work?
6. A task submits a subtask to the same pool and waits on it. What can happen?
7. How would you implement work stealing?

## Reading
- Bryan Mills, *Rethinking Classical Concurrency Patterns*, worker-pool section
- Java `ThreadPoolExecutor` Javadoc (the sizing and queueing sections)
- *Concurrency in Go*, ch. 4 (bounded parallelism) and ch. 6 (work stealing)

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**after `fixed/` + `boundedgo/`**
- Code `[S]` [conc `pool/`](https://github.com/sourcegraph/conc/tree/main/pool): typed bounded pools that re-panic in the caller: safer defaults, less flexibility.

**after shutdown + metrics**
- Read `[S]` [Dropbox: ATF async task framework](https://dropbox.tech/infrastructure/asynchronous-task-scheduling-at-dropbox): at-least-once over SQS means every callback must be idempotent.
- Read `[S]` [Meta: FOQS distributed priority queue](https://engineering.fb.com/2021/02/22/production-engineering/foqs-scaling-a-distributed-priority-queue/): a priority queue on sharded MySQL: durability and multi-tenancy over latency.
