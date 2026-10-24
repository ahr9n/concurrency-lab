# 01: Tooling reference

## Concepts
- **Race detector (`-race`)** is ThreadSanitizer-based. It is dynamic: it only reports races that actually *happen* during the run. Overhead is roughly 2–20× CPU and 5–10× memory. It is not a proof of correctness. Combine it with stress and `-count`.
- **`testing/synctest`** (GA in 1.25; `synctest.Sleep` added in 1.27):
  - `synctest.Test(t, f)` runs `f` in a *bubble* with a fake clock.
  - Time advances only when every goroutine in the bubble is *durably blocked*.
  - `synctest.Wait()` blocks until the others are durably blocked.
  - Ideal for timeouts, TTLs, rate limiters, tickers, backoff.
  - Know what counts as "durably blocked": channel operations on bubble channels, `sync.Cond`, `sync.WaitGroup`, `time.Sleep`. Mutexes and I/O do **not** count.
- **Benchmarks:**
  - `b.Loop()` (1.24+) prevents dead-code elimination of the loop body.
  - `b.RunParallel` for contention.
  - `-cpu=1,2,4,8` for scaling.
  - `benchstat` for statistics.
  - `-benchmem` for allocations.
- **pprof profiles:**
  - `cpu`, `heap`, `allocs`
  - `goroutine`: stacks; `debug=2` is the panic-style dump
  - `mutex`: contention on Unlock; needs `SetMutexProfileFraction`
  - `block`: time blocked on channels, select, and Cond; needs `SetBlockProfileRate`
  - `goroutineleak`: GC-reachability-based detection of permanently blocked goroutines. Experimental in 1.26, default in 1.27.
- **Execution tracer:** `go tool trace` shows per-goroutine timelines, scheduler latency, GC pauses, syscalls, and network blocking. The **flight recorder** (1.25) keeps a rolling in-memory window for production.
- **goleak** (`go.uber.org/goleak`): asserts that no unexpected goroutines remain at the end of a test.
- **Deadlock detection:** the runtime only reports `all goroutines are asleep` when *every* goroutine is blocked. A partial deadlock (some goroutines stuck while others keep running) is invisible to it, and in practice that's the common case.

## Hidden knowledge / gotchas
- `-race` can't see races in code paths your test doesn't execute concurrently. Design tests to *maximise overlap*: use a start gate, many iterations, and `-count`.
- Leak profile limitation: a goroutine blocked on a primitive reachable from a **global variable** can't be flagged, because the profile relies on GC reachability.
- `time.Now()` inside a synctest bubble returns fake time starting at 2000-01-01 UTC.
- The mutex profile attributes the cost to the *unlocker* (the holder that delayed others), not the waiter. Read it accordingly.
- A goroutine dump labels each goroutine with its wait reason (`chan receive`, `semacquire`, `select`, `sync.Mutex.Lock`, `IO wait`) and how long it's been blocked. Learn these.

## Questions you must answer
1. How would you debug a deadlock in a production Go service? (Name the tools in order.)
2. Why can't the race detector prove your program correct?
3. How would you test a TTL cache without real sleeps?
4. Your p99 doubled after a deploy and CPU is flat. What do you look at?
5. How do you detect goroutine leaks in tests and in production?

## Reading
- https://go.dev/doc/articles/race_detector
- https://pkg.go.dev/testing/synctest and the Go blog post *Testing Time (and other asynchronicities)* (Aug 2025)
- https://go.dev/doc/diagnostics
- Go blog: *More predictable benchmarking with testing.B.Loop*; *Flight recorder in Go 1.25*
- Go 1.26 / 1.27 release notes, goroutine leak profile sections: https://go.dev/doc/go1.26

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**after the race experiment**
- Read `[S]` [Uber: Dynamic Data Race Detection in Go Code](https://www.uber.com/blog/dynamic-data-race-detection-in-go-code/): `-race` across a monorepo found ~2,000 races; overhead and triage were the hard part.

**after the leak demo**
- Read `[S]` [Uber: Goroutine Leaks in Enterprise Microservices (CGO 2024)](https://arxiv.org/abs/2312.12002): goleak in tests found 857 leaks; prod profiling cut memory up to 9.2×.
- Read `[S]` [Uber: LeakProf](https://www.uber.com/blog/leakprof-featherlight-in-production-goroutine-leak-detection/): accepts heuristic false positives in exchange for near-zero overhead.
- Code `[E]` [goleak `leaks.go`](https://github.com/uber-go/goleak/blob/master/leaks.go): finds leaks by parsing goroutine stacks; heuristic, needs ignore lists.
