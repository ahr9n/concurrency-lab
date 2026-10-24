# Master roadmap

You aren't on an interview timeline yet, so this plan goes deep before it goes fast:
**understand, build, measure, then explain**. Rough budget: 6–8 h/week. Phases can overlap. Follow the tier paths below.

## Tiers

Every item in every module roadmap is tagged with a tier:

| Tag | Tier | Meaning | Skip it? |
|---|---|---|---|
| `[E]` | **Essential** | Comes up in interviews *and* in day-to-day backend work. Know it cold. | Never |
| `[S]` | **Senior** | Production judgement at scale: overload, durability, tail latency, design trade-offs. What separates senior/staff from mid-level engineers at large-scale companies. | Not if you're aiming for senior |
| `[D]` | **Go depth** | Go-runtime specifics: scheduler, mutex internals, timers, GC. Makes you the Go expert on the team. Rarely asked directly. | Yes, come back later |
| `[T]` | **Theory / rare** | Academic or niche: exotic lock-free queues, classic puzzles, linearizability checkers. Good for curiosity. | Yes |

### Paths
- **Interview-ready (about 8–10 weeks):** all `[E]` items in modules 00–06, 08, 09, 11, and 12 M1.
- **Senior at scale (the target, about 16–18 weeks):** `[E]` + `[S]` everywhere + one *Real-world reading* block per week as the treat.
- **Go master (ongoing, no deadline):** add `[D]`, mostly through the 15-runtime-internals reading sessions.
- **Curiosity:** `[T]`, whenever you feel like it. Never on the critical path.

### Module weight
| Module | Mostly | Notes |
|---|---|---|
| 00 foundations | E | The memory model underpins everything |
| 01 tooling | E/S | Race, synctest, and leaks are E; pprof, trace, and the flight recorder are S |
| 02 sync primitives | E/S | Spinlock is T, starvation mode is D |
| 03 channels | E | Pub/sub and throughput benchmarks are S |
| 04 context / lifecycle | E/S | Re-implementing `context` is D |
| 05 bounded queue | **E** | Interview core |
| 06 worker pool | **E** | Interview core |
| 07 scheduler | S | Job/DAG systems (Airflow / Temporal / Borg-style); cron and timing wheel are T |
| 08 cache | **E** | Interview core; Caffeine-style buffering is S |
| 09 flow control | **S** | Token bucket and retries are E; the rest is senior-at-scale material |
| 10 lock-free | S/T | Copy-on-write pointer is E; MS queue and seqlock are T |
| 11 log writer | **E** | Interview flagship; rotation and fsync failure handling are S |
| 12 WAL-KV | S | M1 is E; M2–M8 are senior systems depth |
| 13 classic problems | E/T | 6 are E (incl. philosophers, readers–writers, bank transfer), the rest T |
| 14 parallel algorithms | S | The crawler is E |
| 15 runtime internals | D | GMP, netpoller, channels, and GC are S (senior engineers are expected to know them) |
| 16 network services | S | Connection pools, coalescing, load balancing: daily life at scale |
| 17 distributed coordination | S | Idempotency is E; vector clocks and Raft are T |
| 18 DB concurrency | **E/S** | Lost update, deadlock retry, and `WithTx` are E; gap locks, SKIP LOCKED, and Spanner are S |
| 19 verification and simulation | S/T | TLA+ on your own queue/log writer and a small DST harness; optional stretch |

### Sustainability rules
- Budget **hours, not scope**: 6–8 h/week. When the hours are used up, stop, even mid-module.
- **Take a break week every 5–6 weeks.** You retain more by coming back to a topic than by grinding.
- Each week, apply one concept to real code you work on, or to an open-source project (workers, errgroup, shutdown). That's mastery practice that costs no extra time.
- Mastery means you can explain it out loud without notes. A checked box doesn't count.

## Phases

| Phase | Weeks | Modules | Milestone you can demonstrate |
|---|---|---|---|
| 1. Ground truth | 1–2 | 00, 01 | Explain the Go memory model and DRF-SC; show a race, the race detector catching it, and a false-sharing benchmark |
| 2. Primitives | 3–5 | 02, 03, 04 | Build a semaphore, barrier, RW lock, and pipeline with cancellation; zero goroutine leaks under test |
| 3. Interview core | 6–9 | 05, 06, 08 | Bounded queue, worker pool, and LRU+TTL cache, each reproducible from a blank file in under 25 min |
| 4. Systems behaviour | 10–12 | 09, 07 | Token bucket (mutex and CAS), circuit breaker, and a delayed/DAG scheduler driven by fake time |
| 5. Flagship | 13–15 | 11 | Group-commit log writer with benchmarks showing the fsync-under-lock regression, plus a crash test |
| 6. Capstone | 16–19 | 12 | WAL-backed KV: recovery, compaction, backpressure, metrics, clean shutdown |
| 7. Depth sweep | continuous | 10, 15 | Walk through `runtime/chan.go`, `internal/sync/mutex.go`, scheduler, and sync.Map internals out loud |
| 8. Interview mode | 20 (+ before any loop) | `interview/` | Timed drills, 3 mock interviews, Java bridge for verbal rounds |

## Definition of done (every module)

1. **At least 2 implementations** where the module asks for them, with trade-offs written down.
2. **Tests** pass `make test M=NN` (race detector, `-count=50`). Anything involving time or blocking also has a `testing/synctest` test.
3. **Benchmark** with `b.RunParallel` / `b.Loop` and recorded numbers (keep them in `NOTES.md`).
4. **`NOTES.md`** in the module, copied from [templates/NOTES.md](templates/NOTES.md) (keep it to one page):
   - the invariants
   - the happens-before edges that make it correct
   - the trade-offs
   - 5 likely interviewer follow-ups, with your answers
5. **Timed drill**: rebuild the core from a blank file without notes, log the time in `interview/drills/LOG.md`.
6. **Explain it out loud** in 5 minutes. If you can't explain it, it isn't done.

## Working rhythm

- **Session start (10 min):** re-derive yesterday's invariant from memory.
- **Build (60–90 min):** code first, then tests that try to break it.
- **Break it (20 min):** add a deliberate bug (`if` instead of `for`, a missing unlock, a lock held across I/O). Check that your tests catch it. If they don't, the tests are weak.
- **Weekly:** one drill, one source-reading hour (`15`) or one *Real-world reading* block, update `NOTES.md`.

## Rules of the lab

- Before using a std / x/ package for something, build a version from lower-level parts. Then use the real one and compare.
- Never hold a lock across I/O, channel sends, or callbacks into user code, unless `NOTES.md` explains why.
- Every goroutine you start has a documented owner and a documented way to stop.
- Measure before you optimise. Record before and after numbers.
- Comments explain *why*, never *how*.
