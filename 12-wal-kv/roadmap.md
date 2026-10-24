# 12: Capstone, a WAL-backed concurrent KV store

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** combine everything into one system you can draw, defend, and demo. This feeds the system-design
round ("design a durable KV / storage engine") and the behavioural round ("a hard perf bug I found").

## Architecture
```
clients ──► ingest (rate limiter 09, concurrency limit) ──► WAL (11, group commit)
                                                           │ durable
                                                           ▼
                                   memtable (sharded RWMutex map or skiplist, 08/02)
                                                           │ size threshold
                                                           ▼
                          flush to immutable sorted segment files (background, 06)
                                                           │
                                     compaction (background, bounded, 07 scheduler)
reads: memtable → immutable memtables → segments (newest first), with snapshot isolation via sequence numbers
```

## Milestones
- [ ] `[E]` **M1:** `Put/Get/Delete` over a WAL plus an in-memory map. Recovery replays the WAL. Crash test from `11`.
- [ ] `[S]` **M2:** memtable rotation. Active memtable → immutable → flushed to a sorted segment file (SSTable-lite with an index block). Truncate the WAL after the flush is durable. (Get the ordering right: segment fsync → manifest update → WAL delete.)
- [ ] `[S]` **M3:** concurrent reads during flush and compaction. Readers hold a refcounted *version* (list of segments). Files are deleted only when no version references them.
- [ ] `[S]` **M4:** compaction in the background: merge segments, drop tombstones and overwritten keys. Throttle it so foreground p99 stays bounded. Make compaction cancellable.
- [ ] `[S]` **M5:** snapshots / MVCC-lite. Every write gets a seq number. `Snapshot()` pins a seq, and reads see only ≤ seq.
- [ ] `[S]` **M6:** write stalls / backpressure. Too many immutable memtables or L0 files slow down or block writes (RocksDB-style).
- [ ] `[S]` **M7:** observability. Metrics (atomics), pprof endpoints, trace regions, flight recorder on slow ops.
- [ ] `[S]` **M8:** a gRPC or HTTP server front-end with graceful shutdown (04), per-tenant limits (09).
- [ ] `[T]` **Stretch:** an iterator API with range scans (merge iterator over sources); bloom filters per segment; batch writes (atomic multi-key).

## Tests
- [ ] `[S]` Randomised model test: run random ops against both your store and a `map` reference, crash-restart at random points, compare.
- [ ] `[T]` Linearizability of single-key Get/Put under concurrency (porcupine).
- [ ] `[S]` Crash at every interesting point (after WAL write before fsync, after segment write before manifest, mid-compaction) with failpoints.
- [ ] `[S]` Long soak test with `-race` and the goroutine-leak profile.

## Deliverables
- [ ] `[S]` `DESIGN.md`: architecture diagram, invariants, crash-consistency argument per step, trade-offs, and what you'd change at 100× scale.
- [ ] `[S]` A benchmark report (write/read throughput, p99 under compaction, effect of group-commit batch size).
- [ ] `[E]` A 10-minute walkthrough you can give in an interview.
