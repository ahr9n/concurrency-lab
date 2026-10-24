# 03: Channels and CSP patterns

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** know every corner of channel semantics, and use the pipeline pattern vocabulary fluently,
including cancellation and avoiding leaks.

## Build
- [ ] `[E]` `semantics_test.go`: a table-driven test that *demonstrates* each cell of the channel-state table (nil / open / closed × send / recv / close / len / cap).
- [ ] `[E]` `patterns/` (each one leak-free and cancellable via `ctx`):
  - [ ] generator
  - [ ] pipeline (stage chain)
  - [ ] fan-out (N workers on one input)
  - [ ] fan-in / merge (WaitGroup + closer goroutine)
  - [ ] **ordered** fan-out: parallel processing that keeps the input order (sequence numbers + reorder buffer, or channel-of-channels)
  - [ ] or-done
  - [ ] tee (with the blocking hazard explained)
  - [ ] bridge (channel of channels)
  - [ ] heartbeat
  - [ ] first-response-wins (replicated requests)
  - [ ] batcher: emit when size ≥ N **or** after T elapses. Use a timer, reset correctly.
  - [ ] debouncer and throttler
- [ ] `[S]` `pubsub/`: an in-process broker with subscribe/unsubscribe. Handle slow subscribers with one of: drop / block / per-subscriber buffer with eviction. Pick one and justify it.
- [ ] `[E]` `select_tricks/`:
  - nil-channel case disabling
  - priority select (two-level select)
  - non-blocking send/recv with `default`
  - timeouts with `time.After` vs `time.NewTimer`

## Experiments
- [ ] `[S]` Measure unbuffered vs buffered (1, 64, 1024) channel throughput against `Mutex` + slice queue.
- [ ] `[D]` Show that `select` picks uniformly among ready cases: count picks over 1M iterations.
- [ ] `[E]` Leak hunt: write 3 pipelines that leak (abandoned consumer, missing close, early return) and catch each with goleak or the leak profile.

## Drill
From a blank file in under 15 minutes: ordered parallel map over a channel with ctx cancellation, no leaks.

## Done when
You can say which goroutine closes each channel and why, for every pattern, and prove there are no leaks.
