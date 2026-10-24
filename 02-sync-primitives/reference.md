# 02: Sync primitives reference

## Contracts
- **`sync.Mutex`:**
  - Not reentrant. Not owned: any goroutine may Unlock, which is legal but almost always a bug.
  - The zero value is usable. Must not be copied after first use (`go vet` copylocks catches this).
  - `TryLock` exists, but the docs say correct uses are rare.
- **Mutex modes** (`src/internal/sync/mutex.go`):
  - *Normal*: a newly arriving goroutine can barge ahead of a woken waiter. Higher throughput.
  - *Starvation*: a waiter that waits more than **1 ms** (`starvationThresholdNs = 1e6`) switches the mutex to direct FIFO handoff.
  - The fast path is one CAS. The slow path spins briefly on multicore, then parks on a semaphore.
- **`sync.RWMutex`:**
  - Writer-preferring: once a writer is waiting, new `RLock`s block.
  - So **recursive read locking deadlocks** if a writer arrives in between.
  - No upgrade from read to write, no downgrade.
  - Use it only when reads dominate *and* critical sections are non-trivial, because `RLock` itself bounces a shared counter's cache line.
- **`sync.Cond`:**
  - `Wait` atomically unlocks `L`, parks, and relocks.
  - **Always call it in a `for !predicate {}` loop**, because of spurious and stolen wakeups.
  - `Signal` wakes one waiter, `Broadcast` wakes all.
  - **There's no timeout and no context.** You can't `select` on it. That's why channels usually win in Go (see Bryan Mills's talk).
- **`sync.WaitGroup`:**
  - `Add` must *happen before* `Wait` can observe zero, so call it before the `go` statement.
  - `WaitGroup.Go(f)` (1.25) does Add+Done for you.
  - Reusing a WaitGroup before the previous `Wait` has returned panics or misbehaves.
- **`sync.Once`:**
  - `Do` guarantees *f has completed* before any `Do` returns, which is why it needs a mutex.
  - A CAS-only Once lets the second caller return before f finishes.
  - If f panics, Once counts as done.
- **`OnceFunc` / `OnceValue` / `OnceValues`** (1.21): if f panics, every call re-panics with the same value.
- **`sync.Pool`:**
  - Per-P local caches with a victim cache (1.13+). Objects survive about one GC cycle.
  - Not a free-list or connection pool, and you get no guarantee of reuse.
  - Don't pool objects of wildly varying size (memory bloat). The standard example is `fmt`'s printers.
- **`sync.Map`:**
  - Built for (1) write-once/read-many keys and (2) disjoint key sets per goroutine.
  - **Since 1.24 it's built on a concurrent hash-trie** (`internal/sync/hashtriemap.go`). The older read-map/dirty-map design is gone, so know both for interviews.
  - Not type-safe. `Range` is not a consistent snapshot.
- **`atomic.Value`** vs **`atomic.Pointer[T]`**: prefer the typed pointer. `Value` panics if you store inconsistent types.

## Building blocks equivalence
| Want | Mutex+Cond | Channel |
|---|---|---|
| mutex | `Mutex` | `make(chan struct{},1)` |
| semaphore(n) | count + Cond | `make(chan struct{}, n)` |
| event / latch | bool + Cond.Broadcast | `close(ch)` |
| cancellable wait | **hard** (needs extra goroutine or Cond hack) | `select { case <-ch: case <-ctx.Done(): }` |
| fair FIFO | explicit queue of waiters | not FIFO under contention |

## Hidden knowledge / gotchas
- Lock ordering is the only general deadlock-avoidance strategy for multiple locks. Order by a global rank or by address (that's what `select` does internally for channels).
- Never call unknown code (callbacks, interfaces, `log` hooks) while holding a lock. It can re-enter or block.
- `defer mu.Unlock()` is fine since 1.14 (open-coded defers, near-zero cost). Unlock early only for a reason you can name.
- Embedding `sync.Mutex` in an exported struct exports `Lock`/`Unlock`. Use a named unexported field.
- `x/sync/semaphore.Weighted` is FIFO. A large request blocks smaller ones behind it, which prevents starvation of large requests.
- Cyclic barrier bug: waiters must wait on the *generation* they arrived in, otherwise a fast thread from round k+1 releases waiters of round k.

## Questions you must answer
1. Why isn't Go's mutex reentrant? How would you design around that?
2. Explain mutex normal vs starvation mode and the trade-off.
3. Why must `Cond.Wait` be in a loop? What's a lost wakeup vs a spurious wakeup?
4. When is `RWMutex` *slower* than `Mutex`?
5. Why does `sync.Once` need a mutex? Why doesn't a CAS suffice?
6. When would you use `sync.Map`? When is it worse than a sharded map?
7. How would you implement a semaphore? A barrier? A latch?
8. How do you make a lock acquisition cancellable?

## Reading
- https://pkg.go.dev/sync and https://pkg.go.dev/sync/atomic
- Source: `$(go env GOROOT)/src/internal/sync/mutex.go`, `src/sync/rwmutex.go`, `src/sync/once.go`, `src/sync/pool.go`, `src/internal/sync/hashtriemap.go`
- Bryan C. Mills, *Rethinking Classical Concurrency Patterns* (GopherCon 2018), with slides
- *The Little Book of Semaphores*, ch. 3–4
- *Concurrency in Go*, ch. 3

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**after `semaphore/`**
- Code `[E]` [sync `semaphore/semaphore.go`](https://github.com/golang/sync/blob/master/semaphore/semaphore.go): weighted FIFO semaphore: fair, but a big waiter blocks the small ones behind it.

**after `syncmap/`**
- Code `[S]` [tailscale `syncs/shardedmap.go`](https://github.com/tailscale/tailscale/blob/main/syncs/shardedmap.go): sharded map; compare with your `sync.Map` benchmark.
