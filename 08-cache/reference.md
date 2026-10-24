# 08: Cache reference

## Key facts
- **Get mutates state in LRU** (move-to-front), so a reader lock is *wrong*. This is the classic probe. Options:
  - an exclusive lock on Get
  - approximate recency (CLOCK: set an access bit atomically under a read lock, then sweep)
  - buffering accesses (Caffeine)
- **TTL races:** check-expiry-then-return must happen under the same lock, or you return an entry that the janitor is deleting.
- **Janitor cost:** a full scan under the lock causes latency spikes. Bound the work per tick (sampling), or use an expiry min-heap / timing wheel.
- **Sharding:** removes single-lock contention, but gives up global LRU order. Each shard needs capacity/N, and hot keys still collide.
- **Stampede / thundering herd:** on a miss, N goroutines all load. `singleflight` dedupes in-flight loads. Also consider stale-while-revalidate and jittered TTLs to avoid synchronized expiry.
- **singleflight gotchas:**
  - One caller's ctx must not cancel everyone's load. Run the load with `context.WithoutCancel` plus its own timeout.
  - Errors are shared too.
  - Use `Forget` to drop a key so the next call starts a new load.
- **Eviction policies:** LRU, LFU, ARC, 2Q, W-TinyLFU (Caffeine/Ristretto, best general hit ratio), CLOCK (approximate LRU, lock-friendlier).
- **Memory:** `map` buckets never shrink. Long-lived caches with churn should periodically rebuild their map. `weak.Pointer` (1.24) + `runtime.AddCleanup` enable caches that let GC reclaim values.

## Production equivalents
- Go: `github.com/hashicorp/golang-lru/v2` (single lock + 2Q/ARC), `github.com/dgraph-io/ristretto` (TinyLFU, buffered), `github.com/maypok86/otter`, `golang.org/x/sync/singleflight`, `sync.Map` (grow-only caches).
- Java: Caffeine, Guava Cache, `ConcurrentHashMap.computeIfAbsent`. Know that CHM's `computeIfAbsent` holds the bin lock during compute, so a slow loader blocks the bin, and a recursive compute on the same map can deadlock or throw.

## Questions you must answer
1. Implement a thread-safe LRU. Why not an RWMutex?
2. Add TTL. Lazy vs active expiry: what are the trade-offs?
3. How would you scale to 64 cores? What do you lose?
4. How do you prevent a cache stampede?
5. How does Caffeine get near-lock-free reads with LRU semantics?
6. How would you bound memory instead of entry count?
7. Write-through vs write-back vs write-around. How do you keep the cache and DB consistent under concurrent writes (delete-after-write, versioning, leases)?

## Reading
- Caffeine design: https://github.com/ben-manes/caffeine/wiki/Design
- Einziger, Friedman & Manes, *TinyLFU: A Highly Efficient Cache Admission Policy* (2017)
- Dgraph blog: *Introducing Ristretto*
- https://pkg.go.dev/golang.org/x/sync/singleflight and its source (~150 lines)

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**after `safelru/` + `ttl/`**
- Read `[S]` [Discord: Why Discord is switching from Go to Rust](https://discord.com/blog/why-discord-is-switching-from-go-to-rust): GC spikes on an LRU-heavy service led to a rewrite; ask what `GOMEMLIMIT` would have changed.
- Code `[E]` [golang-lru `lru.go`](https://github.com/hashicorp/golang-lru/blob/main/lru.go): one mutex over a non-thread-safe core; compare with your `safelru/`.
- Code `[E]` [golang-lru `expirable/expirable_lru.go`](https://github.com/hashicorp/golang-lru/blob/main/expirable/expirable_lru.go): the same with TTL; compare with your `ttl/`.

**after `loader/` (singleflight)**
- Read `[S]` [Facebook: Scaling Memcache (leases)](https://www.usenix.org/conference/nsdi13/technical-sessions/presentation/nishtala): leases fix stale sets and thundering herds; eventual consistency is accepted.
- Code `[E]` [groupcache `singleflight/singleflight.go`](https://github.com/golang/groupcache/blob/master/singleflight/singleflight.go): the original singleflight from Google's download servers, ~60 lines.
- Code `[S]` [groupcache `groupcache.go`](https://github.com/golang/groupcache/blob/master/groupcache.go): peer-owned keys + hot-key cache; no updates, deletes, or TTL, by design.

**after `bufferedlru/`**
- Code `[S]` [ristretto `ring.go`](https://github.com/dgraph-io/ristretto/blob/main/ring.go): lossy striped buffers absorb reads; Caffeine's idea in Go.
- Code `[S]` [ristretto `policy.go`](https://github.com/dgraph-io/ristretto/blob/main/policy.go): TinyLFU admission: Sets are async and can be dropped.

**for consistency follow-ups**
- Read `[S]` [Meta: Cache made consistent](https://engineering.fb.com/2022/06/08/core-infra/cache-made-consistent/): cache consistency from 6 to 10 nines, by measuring inconsistency first.
