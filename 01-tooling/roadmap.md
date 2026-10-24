# 01: Tooling (test, measure, debug)

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** be able to *prove* concurrent code correct and *show* where it's slow or stuck.
Every later module depends on this.

## Build
- [ ] `[E]` `stress/`: a reusable helper `Hammer(t, goroutines, iters, fn)` that starts every goroutine at the same moment behind a shared start gate (a closed channel).
- [ ] `[E]` A test using `testing/synctest.Test`: something that sleeps or times out (e.g. a retry with backoff) runs in virtual time, so it finishes instantly and the same way every time.
- [ ] `[E]` A deliberately leaky function, plus:
  - (a) a `goleak.VerifyNone(t)` test that fails
  - (b) a `synctest` bubble that reports the deadlock
  - (c) the `goroutineleak` pprof profile showing the leak
- [ ] `[E]` Deliberate deadlock: two mutexes locked in opposite order. Capture a goroutine dump (`SIGQUIT` / `debug.Stack` / `pprof goroutine?debug=2`) and read it.
- [ ] `[S]` Contention demo: a hot mutex. Turn on `runtime.SetMutexProfileFraction` and `SetBlockProfileRate`, then find the hotspot with `go tool pprof`.
- [ ] `[S]` Trace demo: `go test -trace=trace.out`, then `go tool trace`. Find the goroutine-blocking and scheduler-latency views.
- [ ] `[S]` Flight recorder: keep a `runtime/trace.FlightRecorder` running and dump it when a request is slow.

## Experiments
- [ ] `[E]` Run the racy counter 100 times without `-race` and count how many runs are wrong. Then run it with `-race` and look at the overhead (time and memory).
- [ ] `[S]` `benchstat` two variants with `-count=10` and read the p-value.
- [ ] `[D]` `GODEBUG=schedtrace=1000` on a CPU-saturating program, and interpret the output.

## Done when
For any later module you can, in under 5 minutes:
- show a race
- show a leak
- show the top contended lock
- show a slow-path trace
