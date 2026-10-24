# 10: Lock-free and wait-free structures

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** understand CAS-based design, its hazards, and when it's worth it (rarely), well enough to
defend or reject it in an interview.

## Build
- [ ] `[S]` `treiber/`: a lock-free stack with `atomic.Pointer[node]`. Discuss why Go's GC makes ABA a non-issue *for this structure*, and construct an ABA scenario that does bite (index-based / pooled nodes).
- [ ] `[T]` `msqueue/`: the Michael–Scott lock-free MPMC queue (dummy node, helping tail advance).
- [ ] `[S]` `spsc/`: a single-producer/single-consumer ring buffer with only atomic loads/stores of head/tail (no CAS). Pad head/tail to separate cache lines. Benchmark it against a channel.
- [ ] `[T]` `mpsc/`: a Vyukov-style intrusive MPSC queue (used by actor mailboxes and the Go runtime).
- [ ] `[E]` `cow/`: a copy-on-write config/routing table. Readers `Load()` an `atomic.Pointer[T]` with no locks, writers copy+swap under a mutex. **This is the one you'll actually use in production.**
- [ ] `[T]` `seqlock/`: a sequence lock for read-mostly small structs. Note why it's awkward in Go: torn reads are races under the memory model, so reads need atomics.
- [ ] `[T]` `stamped/`: a version-stamped pointer to defeat ABA (a pointer to an immutable {ptr, version} box).

## Tests
- [ ] `[T]` Linearizability check: record histories in stress tests and verify them with `github.com/anishathalye/porcupine`.
- [ ] `[S]` -race clean, and N×M push/pop with no loss/duplication.

## Experiments
- [ ] `[S]` Treiber vs `Mutex`+slice stack at 1/4/16/64 goroutines. Expect the mutex to win often. Explain why.
- [ ] `[S]` SPSC ring vs buffered channel. Expect the ring to win by about 5–10×. Explain what the channel pays for.

## Done when
You can explain lock-free vs wait-free vs obstruction-free, ABA, linearizability, and *why you probably shouldn't* ship your own.
