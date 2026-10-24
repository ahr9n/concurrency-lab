# Mock interview cards

Eight company-neutral problems, written the way a concurrency round tends to run: a small start, then constraints that
escalate every 10–15 minutes. Use them for timed self-mocks, or paste a card to Claude with:
*"Run this card as the interviewer. Reveal one constraint at a time, push on my invariants, and don't accept hand-waving."*

Each card: **Prompt** → **Constraints** (in order) → **Probes** (questions to expect at any point) → **Module**.

---

### 1. Bounded blocking queue (45 min) → 05
- **Prompt:** implement `Put`/`Take` for a fixed-capacity queue shared by many goroutines.
- **Constraints:**
  1. Add `Close()`: blocked producers fail, consumers drain the rest.
  2. Add a `ctx` / timeout to both calls.
  3. Make `TakeBatch(n)` return up to n items in one call.
- **Probes:** why `for` and not `if`? Why two conditions? Signal or Broadcast? How do you test it?

### 2. Executor with shutdown (50 min) → 06
- **Prompt:** a fixed pool of N workers with `Submit(task)`.
- **Constraints:**
  1. Bounded queue + a rejection policy.
  2. `Shutdown()` (drain) vs `ShutdownNow()` (cancel, return pending).
  3. `Submit` returns a `Future` with `Get(ctx)`.
  4. A task panics.
- **Probes:** the Submit-vs-Shutdown race; how to size N for CPU-bound vs I/O-bound; is a pool even needed in Go?

### 3. Delayed task scheduler (50 min) → 07
- **Prompt:** `Schedule(at, fn)` runs fn at time `at`.
- **Constraints:**
  1. `Cancel(id)`.
  2. A new earlier task arrives while you're sleeping.
  3. Periodic tasks (fixed-rate vs fixed-delay).
  4. 1M pending tasks.
- **Probes:** one timer vs many; heap vs timing wheel; what if fn is slow? How would you test this without real sleeps?

### 4. Concurrent web crawler (40 min) → 14
- **Prompt:** crawl from a seed URL with a `fetch(url) []url` function, visiting each URL once.
- **Constraints:**
  1. Max N concurrent fetches.
  2. Depth limit.
  3. Per-host limit.
  4. Cancel everything on the first fatal error.
- **Probes:** how do you know you're done? Where is the dedupe check-and-insert atomic? Is it leak-free?

### 5. Thread-safe rate limiter (35 min) → 09
- **Prompt:** `Allow()` admits at most R requests/sec with burst B.
- **Constraints:**
  1. `Wait(ctx)`.
  2. Per-key limits for 1M keys, with idle eviction.
  3. Make it lock-free.
  4. Many servers share one global limit.
- **Probes:** the check-then-act race; monotonic time; when does CAS lose to a mutex?

### 6. Thread-safe LRU cache (45 min) → 08
- **Prompt:** `Get`/`Put` with capacity C, safe for concurrent use.
- **Constraints:**
  1. TTL per entry.
  2. `GetOrLoad` with a stampede guarantee (one load per key).
  3. Scale to 64 cores.
  4. Survive a restart (persist it).
- **Probes:** why not an RWMutex? What does sharding give up? Where's the race in TTL expiry?

### 7. Durable log writer (60 min) → 11
- **Prompt:** many goroutines call `Append(rec)`. When it returns, the record must survive a crash. One file.
- **Constraints:**
  1. Throughput: no fsync per record.
  2. `AppendAsync` + `WaitDurable(lsn)`.
  3. Rotate files without stalling writers.
  4. fsync returns an error.
  5. Prove nothing acknowledged is lost.
- **Probes:** what's under the lock and what isn't? The happens-before edges across the buffer swap; page cache vs disk; recovery and torn tails.

### 8. Read-mostly config / routing table (30 min) → 10, 02
- **Prompt:** thousands of goroutines read a routing table per request; it updates every few seconds.
- **Constraints:**
  1. Readers must never block.
  2. Updates are partial (one route at a time).
  3. A reader needs a consistent snapshot across two lookups.
- **Probes:** RWMutex vs `atomic.Pointer` copy-on-write; cost of a write; memory reclamation (why GC makes this easy in Go).

---

**After each mock, in `drills/LOG.md`:** time, which constraint you got stuck on, which probe you answered weakly.
Then re-read that module's `NOTES.md`.
