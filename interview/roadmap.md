# Interview mode

Start this phase after phase 6 of [ROADMAP.md](../ROADMAP.md), and repeat it in compressed form before any loop.

## 1. Timed drills (log every attempt in `drills/LOG.md`)
Rebuild from a blank file with no notes, explaining aloud as if an interviewer were present.

| Drill | Target time | Source module |
|---|---|---|
| Bounded queue (Mutex + 2 Cond, Close, ctx) | 15 min | 05 |
| Worker pool + graceful shutdown + Future | 20 min | 06 |
| Thread-safe LRU (+ TTL, + sharding talk) | 20 min | 08 |
| Token bucket (mutex, lazy refill, Wait(ctx)) | 12 min | 09 |
| **Durable log writer, group commit / double buffer** | **25 min** | 11 |
| Writer-preferring RW lock | 12 min | 02 |
| DAG executor with parallelism limit | 25 min | 07 |
| Concurrent web crawler with dedupe | 15 min | 14 |
| Pseudocode-only version of the log writer (no compiler) | 15 min | 11 |

Practise the pseudocode version too: some rounds are done in pseudocode or are mostly verbal.

**Bug hunt:** the 12 broken programs in [drills/bughunt/](drills/bughunt/README.md) (`make bug N=01`). Target: under 10 min each, naming the tool that should have caught it. Debugging concurrent code is an interview format of its own.

## 2. Mock interviews (minimum 3)
Use the cards in [mocks.md](mocks.md): 8 problems with escalating constraints and probes. Do at least one
log-writer card (#7), one verbal-only session (memory model, CAS/ABA, map and mutex internals), and one scale-it card (#4 or #6).

## 3. The answer protocol (use it in every drill)
1. Restate the requirements, and ask about durability, ordering guarantee, capacity, and failure behaviour.
2. State the invariants *before* coding.
3. Write the simplest correct version (a mutex), then explain where it's slow.
4. Improve it (batching, sharding, swapping outside the lock) and say why it's still correct (happens-before).
5. Explain the tests: stress, crash, what you'd inject.
6. Name the production equivalent and when you'd use it instead.

## 4. Java bridge (2 sessions)
Many interviewers think in Java or C++, so verbal probes are often Java-flavoured. See the Java section in `reference.md`.

## 5. Before a real loop
- [ ] Confirm with the recruiter that Go is allowed in the coding and concurrency rounds, and whether the round uses pseudocode, a runnable IDE, or Google Docs. **Not confirmed by any source.**
- [ ] Re-run every drill once in the final week.
- [ ] Re-read every `NOTES.md` (one page each).
