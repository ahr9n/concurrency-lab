# 05: Bounded queue reference

## Invariants
- `0 ≤ count ≤ capacity` at every unlock.
- Every successful `Put` is matched by exactly one successful `Take` (no loss, no duplication).
- After `Close`, no `Put` succeeds. `Take` succeeds until the buffer is drained.
- Waiters always re-check the predicate after waking (`for`, never `if`).

## Key design points
- **Two condition variables:** producers wait on `notFull`, consumers on `notEmpty`. Each side signals the *other* side's condition, so you can use `Signal` (wake one) instead of `Broadcast` (thundering herd).
- **One condition variable + `Signal` is a bug:** a producer can wake another producer instead of a consumer, and everyone sleeps. With one Cond you must use `Broadcast`.
- **Spurious / stolen wakeups:** between `Signal` and the woken goroutine reacquiring the lock, another goroutine can consume the slot. That's why the `for` loop exists.
- **Close semantics:** `Broadcast` both conds so every waiter re-checks the `closed` flag.
- **Backpressure:** a bounded queue *is* the backpressure mechanism. Blocking pushes pressure upstream to producers. Alternatives: reject (`TryPut`), drop oldest, drop newest. The right choice depends on the domain (telemetry tolerates drops; payments don't).
- **Cancellation:** `Cond` can't select. Options:
  - (a) `context.AfterFunc(ctx, func(){ mu.Lock(); cond.Broadcast(); mu.Unlock() })`
  - (b) a channel-based implementation
  - (c) timed polling (bad)

## Production equivalents
- Go: a buffered channel (most cases); a `chan` + `ctx` select; `container/ring` + a mutex for drop-oldest.
- Java: `ArrayBlockingQueue` (one lock, two conditions), `LinkedBlockingQueue` (two locks), `SynchronousQueue` (zero capacity, like an unbuffered chan), `LinkedTransferQueue`, Disruptor (ring buffer, sequence barriers).

## Follow-ups interviewers ask
1. Why `while` and not `if`?
2. Why two conditions? What breaks with one?
3. Signal vs Broadcast: when does each one break?
4. How do you shut down cleanly? What about producers blocked in Put?
5. How do you add a timeout to Take?
6. How would you reduce contention? (Two-lock queue, batching, sharding, lock-free MPSC.)
7. Unbounded queue: what goes wrong in production? (Memory blowup, latency hiding, no backpressure.)
8. Fairness: is your queue FIFO for *waiters*? Does it matter?
9. How would you make it lock-free? (See `10-lock-free`.)
10. What is Little's Law, and how does it size the queue? (L = λW)

## Reading
- *The Little Book of Semaphores*, §4.1 (producer–consumer)
- Java `ArrayBlockingQueue` / `LinkedBlockingQueue` source. Read both; they're short.
- Michael & Scott 1996 (two-lock queue section)
- Bryan Mills, *Rethinking Classical Concurrency Patterns*: the "condition variables" section

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**after your 3 implementations**
- Read `[E]` [Slack: Scaling Slack's Job Queue](https://slack.engineering/scaling-slacks-job-queue/): Redis ran out of memory and the queue locked up; fixed with a Kafka front buffer. Textbook backpressure failure.
- Code `[E]` [client-go `util/workqueue/queue.go`](https://github.com/kubernetes/client-go/blob/master/util/workqueue/queue.go): dirty + processing sets: dedupe, and a key is never processed concurrently. Read this first.

**when thinking about backpressure policy**
- Read `[S]` [AWS: Avoiding insurmountable queue backlogs](https://aws.amazon.com/builders-library/avoiding-insurmountable-queue-backlogs/): FIFO backlogs make recovery slow; per-tenant queues, LIFO, and sidelining old work.
