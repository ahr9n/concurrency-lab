# 13: Classic problems reference

## Concepts
- **Coffman conditions** for deadlock (all four are required): mutual exclusion, hold-and-wait, no preemption, circular wait. Breaking any one prevents deadlock. Lock ordering breaks circular wait. Try-lock + backoff breaks hold-and-wait (and risks livelock).
- **Deadlock vs livelock vs starvation:**
  - Deadlock: nobody progresses and everyone is blocked.
  - Livelock: everyone is busy but nobody progresses.
  - Starvation: some goroutine never progresses.
- **Handling strategies:** prevention (ordering), avoidance (Banker's algorithm), detection + recovery (wait-for graph cycles, timeouts), ignoring it (the ostrich approach).
- **Dining philosophers solutions:** global fork ordering (asymmetric), arbiter/waiter (a semaphore of N−1 seats), Chandy–Misra (dirty/clean forks, distributed).
- **Readers–writers:** reader preference starves writers, writer preference starves readers, and the fair variant uses a turnstile / FIFO queue.
- **Turnstile / rendezvous / barrier / multiplex** are the reusable semaphore idioms from the Little Book.

## Questions you must answer
1. State the Coffman conditions and show which one your solution breaks.
2. How would you detect a deadlock at runtime in a Go service?
3. Transfer money between accounts A and B concurrently without deadlock.
4. Dining philosophers: give 3 solutions and their trade-offs.
5. Livelock: give an example and a fix.

## Reading
- *The Little Book of Semaphores* (all of ch. 3–5; ch. 6–7 for more)
- LeetCode Concurrency problem set (1114–1279). Solve them in Go locally, since LeetCode's concurrency section doesn't offer Go for all of them.
- Coffman, Elphick & Shoshani, *System Deadlocks* (1971)
- Linux `lockdep` design docs (for #15)

## Real-world reading
None. The classics stand on their own: use *The Little Book of Semaphores*.
