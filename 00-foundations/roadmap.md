# 00: Foundations

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** a correct mental model of goroutines, the Go memory model, data races, and atomics,
before you touch any higher-level primitive.

## Build
- [ ] `[E]` `racy/`: a counter incremented by N goroutines with plain `++`. Observe lost updates. Run it under `-race`.
- [ ] `[E]` `mapcrash/`: concurrent writes to a plain `map`. Observe `fatal error: concurrent map writes`, and confirm `recover` can't catch it.
- [ ] `[E]` `counter/` with 5 implementations behind one interface `Inc()` / `Load()`:
  - [ ] `sync.Mutex`
  - [ ] `atomic.Int64`
  - [ ] a CAS loop (`CompareAndSwap`) written by hand
  - [ ] a sharded counter (per-shard `atomic.Int64`), **no padding**
  - [ ] a sharded counter **with padding** to a cache line (128 B on Apple M-series, 64 B on most x86)
- [ ] `[E]` `publish/`: publish a struct from one goroutine to another:
  - (a) a plain bool flag: broken, explain why
  - (b) `atomic.Bool`
  - (c) `atomic.Pointer[T]`
  - (d) a channel close
- [ ] `[D]` `tearing/`: race on an `interface{}` or slice header and show a torn read. A multi-word race can corrupt memory.

## Experiments
- [ ] `[S]` Benchmark all counters at `-cpu=1,2,4,8,16` and record the numbers. Explain the curve, and where padding starts to pay off.
- [ ] `[S]` Print `runtime.NumGoroutine()` and `GOMAXPROCS` inside a container with a CPU limit. Check what changed in 1.25.
- [ ] `[D]` Start 1M goroutines that each sleep. Measure RSS. Explain the stack size and growth model.

## Tests
- Counter correctness: `N` goroutines × `M` increments == `N*M`, under `-race -count=50`.

## Drill
Without notes, write the CAS-loop increment and the padded sharded counter in under 10 minutes.

## Done when
You can explain all of these out loud:
- DRF-SC
- why `volatile`-style visibility isn't enough for `x++`
- what a happens-before edge is, with four concrete Go examples
