# 09: Flow control (rate limiting, overload, resilience)

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** protect a service and its dependencies under load. These are concurrent
check-then-act problems with time in the mix.

## Build
- [ ] `[E]` `tokenbucket/`:
  - [ ] mutex version: lazy refill on acquire (no ticker goroutine), monotonic time, overflow-safe math
  - [ ] lock-free version: a single `atomic.Uint64` CAS loop (pack tokens + timestamp, or GCRA "theoretical arrival time")
  - [ ] `Allow()`, `Wait(ctx)`, `Reserve(n)`. Compare the semantics with `x/time/rate`.
- [ ] `[S]` `slidingwindow/`: fixed window, sliding log, sliding window counter. Compare accuracy and memory.
- [ ] `[S]` `keyed/`: a per-key limiter map (user/tenant) with idle eviction. Avoid a global lock on the hot path (sharding).
- [ ] `[S]` `concurrencylimit/`: max in-flight requests (semaphore), plus an *adaptive* limit (AIMD or Vegas-style, like Netflix concurrency-limits).
- [ ] `[S]` `breaker/`: a circuit breaker (closed → open → half-open) with a rolling error-rate window and a limited number of half-open probes. Transitions must be race-free.
- [ ] `[E]` `retry/`: exponential backoff with **full jitter**, a retry budget (a token bucket on retries), and respect for ctx deadlines.
- [ ] `[S]` `hedge/`: hedged requests. After the p95 delay, send a second request; take the first response and cancel the other.
- [ ] `[S]` `shed/`: load shedding. Reject at the door when queue time or in-flight count exceeds a threshold (LIFO-under-overload option).

## Tests (all with `synctest` virtual time)
- [ ] `[E]` The token bucket admits exactly `burst + rate×t` over t, under 64 concurrent callers. No over-admission.
- [ ] `[S]` The breaker opens at the threshold, half-opens after the cooldown, and allows only K probes concurrently.
- [ ] `[E]` Retry: the attempt schedule is within the jitter bounds and stops at the ctx deadline.

## Experiments
- [ ] `[S]` Mutex vs CAS token bucket at 1/8/64 goroutines. Watch the CAS retry rate climb.
- [ ] `[S]` Simulate a slow dependency: with vs without a breaker/budget, measure goroutine count and latency (show the retry storm).

## Drill
Token bucket (mutex + lazy refill) with `Wait(ctx)`, in under 12 minutes.
