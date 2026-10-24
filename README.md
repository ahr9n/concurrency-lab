# concurrency-lab

A project-based curriculum for Go concurrency, from the memory model up to a crash-safe
write-ahead-log KV store. Aimed at interview-grade depth (dedicated concurrency / system-programming
rounds) and beyond and at real production judgement.

## How to use this repo

- Start at [ROADMAP.md](ROADMAP.md) — the master sequence, phases, and definition of done.
- Global reading list and sources: [REFERENCE.md](REFERENCE.md).
- Every roadmap item is tagged by tier: `[E]` essential · `[S]` senior · `[D]` Go depth · `[T]` theory. See [ROADMAP.md](ROADMAP.md#tiers).
- Every module folder has:
  - `roadmap.md` — what to build, experiments to run, tests to write, the timed drill, exit criteria.
  - `reference.md` — concepts, hidden gotchas, questions you must be able to answer, reading, and **real-world reading** (big-tech reports + production Go code, mapped to the step where you read them).
  - `NOTES.md` (you write it) — your one-page interview answer sheet for that topic.
- `interview/` — question bank, mock cards, bug-hunt drills, drills log, Java bridge.

## Layout

| Folder | Topic | Track |
|---|---|---|
| `00-foundations` | goroutines, memory model, data races, atomics, false sharing | Core |
| `01-tooling` | `-race`, `synctest`, benchmarks, pprof, trace, leak detection | Core |
| `02-sync-primitives` | Mutex, RWMutex, Cond, Once, WaitGroup, Pool, Map; build your own | Core |
| `03-channels-csp` | channel semantics, select, timers, pipeline patterns | Core |
| `04-context-lifecycle` | cancellation, deadlines, errgroup, leaks, graceful shutdown | Core |
| `05-bounded-queue` | producer–consumer, backpressure (3 implementations) | Interview Tier 1 |
| `06-worker-pool` | executor, rejection policies, shutdown, futures | Interview Tier 1 |
| `07-task-scheduler` | delayed / periodic / priority / DAG-dependency scheduling | Interview Tier 2 |
| `08-cache` | thread-safe LRU + TTL, sharding, singleflight | Interview Tier 1 |
| `09-flow-control` | rate limiters, circuit breaker, load shedding, retries, hedging | Interview Tier 2 |
| `10-lock-free` | CAS, Treiber stack, MS queue, SPSC ring, ABA, copy-on-write | Depth |
| `11-log-writer` | synchronous log writer: naive → group commit + double buffer | Interview flagship |
| `12-wal-kv` | capstone: WAL-backed KV with recovery and compaction | Capstone |
| `13-classic-problems` | philosophers, readers–writers, H2O, barber, print-in-order… | Breadth |
| `14-parallel-algorithms` | parallel sort, map-reduce, crawler, Amdahl/Gustafson | Breadth |
| `15-runtime-internals` | GMP scheduler, hchan, mutex impl, netpoller, GC interaction | Depth |
| `16-network-services` | servers, connection pools, pub/sub broker, deadlines | Breadth |
| `17-distributed-coordination` | leases, fencing tokens, distributed locks, idempotency | Breadth |
| `18-db-concurrency` | isolation anomalies, locking reads, deadlock retry, Spanner transactions | Daily work |
| `19-verification-and-simulation` | TLA+ model checking, deterministic simulation testing | Depth (stretch) |
| `interview/` | question bank, mock cards, bug-hunt drills, Java bridge | Prep |

## Tooling
Toolchain: Go 1.27 (`go.mod`). Lint needs golangci-lint **v2.13+** (`go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0`).

```sh
make test              # race tests, all modules (-count=50); make test M=05 for one module
make bench M=11        # benchmarks for one module
make vet lint          # static checks (CI runs vet, lint, race tests on every push)
make bug N=04          # one bug-hunt drill; make bugs for all; make bug-reset to start over
make progress          # checked roadmap items per tier
```

- Start every module's notes from [templates/NOTES.md](templates/NOTES.md).
- Package names can't start with a digit: in `05-bounded-queue/condq/`, use `package condq`.

## Disclaimer
Personal study notes, not affiliated with or endorsed by any company mentioned. Company and product names are
trademarks of their respective owners and are used only to cite their **publicly published** engineering blogs,
papers, and open-source code. This repo contains no confidential or non-public material, and no interview
questions attributed to any specific employer. Third-party links are provided for reference only. Their content
belongs to their authors. Practice problems are generic, well-known computer-science exercises.
