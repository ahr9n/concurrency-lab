# Interview reference

Company-neutral. Formats and problems below are common patterns in concurrency interviews in general.

## Most common problems
| Problem | Module |
|---|---|
| Bounded blocking queue / producer–consumer | 05 |
| Thread pool / executor with shutdown | 06 |
| Delayed / scheduled task executor | 07 |
| Multithreaded web crawler | 14 |
| Thread-safe rate limiter (token bucket) | 09 |
| Readers–writer lock | 02 |
| Concurrent cache / map (LRU, hashmap) | 08 |
| Thread-safe buffered / durable log writer | 11 |
| Debugging racy code (producer–consumer, banking) | 13, drills |
| Primitives from scratch (lock, barrier, CAS, lock-free queue) | 02, 10 |

## Round formats you may meet
- **A dedicated concurrency round** (45–75 min): one problem that grows ("now add shutdown", "now make it durable", "now remove the lock").
- **A follow-up inside a coding round:** "how would you make this thread-safe?"
- **Debug concurrent code:** find the race or deadlock in a given codebase.
- **Verbal deep-dive:** internals (maps, mutexes, memory model, CAS), with little code.
- **Pseudocode-only:** no compiler. Your invariants have to carry the correctness argument.

## Common verbal probes
- What "durable" means on return (fsync vs page cache); amortising fsync across callers
- Ordering guarantees across threads; async submit + wait-for-completion
- Backpressure policy (block / drop / reject / caller-runs)
- Deadlock: causes, lock ordering, how you'd debug it in production
- Lock vs CAS: when lock-free wins and when it loses; ABA
- Visibility and happens-before (`volatile` / atomics); why `x++` isn't atomic
- Condition variables vs semaphores; spurious wakeups
- Why the plain hashmap isn't thread-safe, and how the concurrent one works
- "Remove the lock: does it still work?" Argue it from your invariants

---

## Master question bank (answer each in 2 minutes, out loud)

- **Memory model and fundamentals:** see `00`. DRF-SC; happens-before edges; why `x++` races; atomics vs mutex; false sharing; Go atomics vs Java volatile vs C++ memory_order.
- **Primitives:** see `02`. Mutex modes; reentrancy; Cond loop; RWMutex recursion deadlock; Once and CAS; sync.Map design (old and 1.24+); semaphore, barrier, latch.
- **Channels:** see `03`. State table; ownership and closing; select fairness and lock ordering; ordered fan-out; leaks.
- **Lifecycle:** see `04`. Cancellation propagation; graceful shutdown order; cancelling blocked I/O.
- **Problems:** see `05`–`14`. Every "Questions you must answer" section.
- **Runtime:** see `15`. GMP; syscalls vs netpoller; preemption; channel send walkthrough.
- **Design:** see `11`, `12`, `17`. Durability, group commit, crash consistency; distributed locks and fencing.
- **Debugging:** see `01`. Deadlock in prod; goroutine leak; contention; p99 regression.

---

## Java bridge (many interviewers think in Java or C++)
- **JMM:** happens-before.
  - A `volatile` write → subsequent read of that variable is an edge. It gives visibility and ordering, but no atomicity for compound operations.
  - `synchronized` exit → the next entry is an edge.
  - `final` fields get safe publication after the constructor.
  - Thread start and join are edges.
- **Double-checked locking** needs `volatile`. Without it, another thread can see a partially constructed object (reordering).
- **`HashMap` isn't thread-safe:**
  - Lost updates on concurrent put.
  - Before Java 8, a concurrent resize could form a **cycle in a bucket list, causing an infinite loop**.
  - Treeified bins can also corrupt.
  - With no happens-before, readers can see stale or partial state.
  - Go's equivalent is the runtime-detected `fatal error: concurrent map writes`.
- **`ConcurrentHashMap`:**
  - **Java 7:** segments, i.e. 16 `ReentrantLock` stripes.
  - **Java 8+:**
    - CAS to insert into an empty bin (`casTabAt`).
    - `synchronized` on the bin's head node for non-empty bins.
    - Lock-free reads through `volatile` `Node.val` / `next` and `tabAt` (acquire reads via Unsafe/VarHandle).
    - **Cooperative resize:** `ForwardingNode`; threads help transfer bins.
    - `size()` uses `CounterCell`s (LongAdder-style striping).
    - Bins treeify at 8 entries.
  - `computeIfAbsent` holds the bin lock during compute, so don't do slow or recursive work in it.
  - Compare with Go: `sync.Map` (hash-trie since 1.24) and sharded maps.
- **CAS / atomics:** `AtomicInteger.compareAndSet`; ABA and `AtomicStampedReference`; `LongAdder` vs `AtomicLong` under contention; `VarHandle` access modes (plain / opaque / acquire-release / volatile).
- **Locks:**
  - `ReentrantLock` (fair/unfair, `tryLock(timeout)`, `lockInterruptibly`, multiple `Condition`s) is built on **AQS**: a CLH-variant queue of waiters plus a state int.
  - Also `ReentrantReadWriteLock`, and `StampedLock` (optimistic reads, not reentrant).
- **Executors:**
  - `ThreadPoolExecutor` (core/max/queue/handler). It grows past core only when the queue is full.
  - `ScheduledThreadPoolExecutor` (the "scheduled executor" problem), `ForkJoinPool` (work stealing), `CompletableFuture`.
- **Virtual threads (JDK 21):**
  - Good for throughput on I/O-bound work. Don't pool them.
  - Before JDK 24, `synchronized` pins the carrier thread (JEP 491 fixed this in 24).
  - Compared with goroutines, both are M:N. Goroutines have growable stacks and a preemptive scheduler; virtual threads mount and unmount on carrier threads.
- **Interview line:** "In Go I'd use X. The Java equivalent is Y, and its internals are Z."
