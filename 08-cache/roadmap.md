# 08: Thread-safe cache (LRU + TTL)

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** a Tier-1 interview problem. Get it correct first, then scale it, and know exactly why each step is needed.

> Common interview variant: **sharded locks + LRU + WAL persistence**. TTL is a common follow-up.
> After `12-wal-kv` M1, add a `persistent/` variant: every Put/Delete goes through the `11` writer, and recovery rebuilds the cache.

## Build (progressively)
- [ ] `[E]` `lru/`: `map[K]*node` + an intrusive doubly-linked list with sentinel head/tail. Single goroutine, O(1) Get/Put/evict.
- [ ] `[E]` `safelru/`: wrap it with one `sync.Mutex`. **Don't** use an `RWMutex` read lock for Get, because Get mutates recency. Write that down.
- [ ] `[E]` `ttl/`: per-entry `expiresAt`.
  - [ ] lazy expiry on Get
  - [ ] a background janitor with bounded work per tick (sample like Redis, or keep an expiry heap)
  - [ ] janitor stop via ctx
- [ ] `[E]` `sharded/`: N shards by `maphash` of the key. Per-shard LRU means *approximate* global LRU. Document that trade-off.
- [ ] `[E]` `loader/`: `GetOrLoad(ctx, key, loadFn)` with `singleflight` (build your own first): one load per key under a stampede, waiters share the result, and a caller's ctx cancel doesn't cancel the shared load.
- [ ] `[S]` `bufferedlru/` (stretch): Caffeine-style. Reads record accesses into a lossy striped buffer, and a single drainer applies the recency updates under the lock.
- [ ] `[S]` Stats: hits, misses, evictions, loads via atomics.

## Tests
- [ ] `[E]` Single-thread eviction order correctness (table tests).
- [ ] `[E]` Concurrent stress: `Len() ≤ capacity` always, no panics, `-race` clean.
- [ ] `[E]` TTL with `synctest`: no stale read after expiry, and the janitor actually reclaims entries.
- [ ] `[E]` Stampede: 1000 concurrent `GetOrLoad` for one key → `loadFn` called once.

## Experiments
- [ ] `[S]` Throughput: single-lock vs RWMutex-but-wrong (show the race) vs sharded (8/64/256) vs buffered, at a 90/10 read/write mix with Zipf keys (`math/rand/v2` Zipf).
- [ ] `[S]` Hit-ratio cost of sharding (approximate LRU) vs exact LRU.

## Drill
Thread-safe LRU + TTL with lazy expiry, in under 20 minutes.
