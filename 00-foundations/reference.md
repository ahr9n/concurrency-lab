# 00: Foundations reference

## Concepts
- **Concurrency vs parallelism.** Concurrency is about program structure. Parallelism is about simultaneous execution. GOMAXPROCS bounds the parallelism.
- **Goroutines.** M:N scheduled. The stack starts small (about 2 KB; adaptive initial size since 1.19) and grows by copying. Creating one is cheap; *leaking* one is not.
- **The Go memory model (revised 2022):**
  - **happens-before** = sequenced-before ∪ synchronized-before, transitively.
  - **DRF-SC:** a data-race-free program behaves sequentially consistently.
  - A race on a **single-word** value can only observe some written value. A race on **multi-word** values (interfaces, slices, strings, structs) can tear and corrupt memory.
  - `sync/atomic` operations are **sequentially consistent**. This is stronger than C++ `relaxed`/`acq_rel`, and comparable to Java `volatile` for that single variable.
- **Synchronizing edges you must know:**
  - goroutine start (the `go` statement happens before the goroutine runs)
  - channel send → matching receive completes
  - close → receive that sees the close
  - for an unbuffered channel, receive → send completes
  - for a channel with capacity C, the k-th receive → the (k+C)-th send completes. This is the semaphore rule.
  - `Mutex` n-th Unlock → (n+1)-th Lock
  - `Once.Do(f)` return of f → any `Do` return
  - atomics
- **Not** synchronizing: goroutine exit (nothing waits on it unless you add a WaitGroup or channel), `time.Sleep`, or "it worked on my laptop".
- **Atomics:** `atomic.Int32/Int64/Uint64/Bool/Pointer[T]/Value`. The typed ones (1.19+) guarantee 64-bit alignment even on 32-bit platforms. `And`/`Or` were added in 1.23.
- **Read-modify-write:** `x++` is load, add, store. Visibility alone doesn't make it atomic.
- **Cache coherence and false sharing:** cores invalidate each other's whole cache lines, so independent hot variables on the same line ping-pong between cores. Fix: padding (`_ [N]byte`, or `golang.org/x/sys/cpu.CacheLinePad`).

## Hidden knowledge / gotchas
- `fatal error: concurrent map writes` is a *runtime throw*, not a panic. `recover` can't catch it. The detection is best-effort, not guaranteed.
- Since 1.22, loop variables are per-iteration, so the classic `go func(){ use(i) }()` closure bug is gone. Know why it existed.
- A panic in any goroutine kills the whole process. You can't recover it from the parent goroutine.
- `runtime.Gosched()` is not synchronization.
- A benchmark that gets "faster with more threads" on a sharded counter can still be slower than a single-threaded loop. Always compare against `-cpu=1`.
- Apple M1–M4 cores use **128-byte** cache lines. Pad accordingly, or measure.
- There is no Go equivalent of Java's `LongAdder` in std. You build it yourself (sharded and padded).

## Questions you must answer
1. What exactly does DRF-SC promise, and what happens when a program is *not* race-free?
2. Why isn't `x++` atomic, and why doesn't making `x` "visible" fix it?
3. Give 4 different happens-before edges in Go.
4. Why is a race on an interface value more dangerous than a race on an `int`?
5. Atomic vs mutex: when is each faster, and why? What about at 64 cores?
6. What is false sharing? How would you detect it? How would you fix it?
7. Are Go atomics acquire/release or sequentially consistent? Why did Go choose that?

## Reading
- https://go.dev/ref/mem (read it twice)
- Russ Cox, *Hardware Memory Models*, *Programming Language Memory Models*, *Updating the Go Memory Model*: https://research.swtch.com/mm
- *Concurrency in Go*, ch. 1–2
- Go blog: *Share Memory By Communicating*

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**after `racy/` + `publish/`**
- Read `[E]` [Uber: Data Race Patterns in Go](https://www.uber.com/blog/data-race-patterns-in-go/): the races come from Go's convenient idioms (closures, slices, maps); you'll recognise every pattern.

**after the sharded counter**
- Code `[S]` [tailscale `syncs/shardedint.go`](https://github.com/tailscale/tailscale/blob/main/syncs/shardedint.go): per-CPU sharded counter in production: fast writes, slow reads.

**during the GOMAXPROCS experiment**
- Code `[D]` [automaxprocs `maxprocs/maxprocs.go`](https://github.com/uber-go/automaxprocs/blob/master/maxprocs/maxprocs.go): cgroup-aware GOMAXPROCS; mostly built into Go since 1.25.
