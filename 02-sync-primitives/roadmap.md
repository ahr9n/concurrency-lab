# 02: Sync primitives (use them, then build them)

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** know every `sync` type's exact contract, and be able to build the classic primitives
from Mutex+Cond and from channels.

## Build (from scratch, each with tests and a benchmark against the std version)
- [ ] `[T]` `spinlock/`: a CAS-based spinlock with `runtime.Gosched` backoff. Benchmark it against `sync.Mutex` under low and high contention.
- [ ] `[S]` `chanmutex/`: a mutex from a `chan struct{}` with capacity 1. Add `TryLock` and `LockContext(ctx)`.
- [ ] `[E]` `semaphore/`: a counting semaphore, in two versions:
  - (a) a buffered channel
  - (b) Mutex+Cond with weighted acquire. Compare to `x/sync/semaphore.Weighted`, including its FIFO fairness.
- [ ] `[S]` `latch/`: a CountDownLatch (`CountDown`, `Wait`, `WaitContext`).
- [ ] `[S]` `barrier/`: a cyclic barrier for N parties, reusable across generations. (Watch the generation-count trap.)
- [ ] `[E]` `rwlock/` (`08-cache` depends on this), in two versions:
  - [ ] reader-preferring
  - [ ] writer-preferring (mirroring `sync.RWMutex` semantics)
- [ ] `[S]` `once/`: re-implement `sync.Once` (fast-path atomic + slow-path mutex). Write down why a pure CAS version is **wrong**.
- [ ] `[E]` `oncevalue/`: use `sync.OnceValue` / `OnceValues` / `OnceFunc`. Note the panic semantics.
- [ ] `[S]` `pool/`: a buffer pool with `sync.Pool`. Benchmark allocations with and without it, and show that pooled objects get dropped across GC cycles.
- [ ] `[S]` `syncmap/`: benchmark `sync.Map` vs `map+Mutex` vs `map+RWMutex` vs a sharded map, for (1) write-once/read-many and (2) disjoint-key workloads.

## Experiments
- [ ] `[D]` Force mutex **starvation mode**: one goroutine re-locks in a tight loop while others wait. Measure waiter latency and compare to the 1 ms threshold.
- [ ] `[E]` Recursive `RLock` deadlock: reader A holds RLock, writer W calls Lock, then A calls RLock again. Reproduce it and explain it.
- [ ] `[E]` Copy a struct containing a `sync.Mutex` and see what `go vet` (copylocks) says.

## Drill
In under 15 minutes from a blank file: a writer-preferring RW lock and a cyclic barrier.

## Done when
You can state the exact contract of every `sync` type and the bug each one is designed to prevent.
