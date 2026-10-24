# 13: Classic synchronization problems

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** breadth and fluency. Each one is a 15–30 minute kata that drills one specific idea.
Solve each **twice**: once with Mutex/Cond (or semaphores), once with channels.

## Problems (✓ when both versions pass under -race)
| # | Problem | What it drills |
|---|---|---|
| [ ] `[E]` 1 | Print in order (first/second/third) | happens-before signalling |
| [ ] `[E]` 2 | Print FooBar alternately | ping-pong handoff |
| [ ] `[T]` 3 | Zero-Even-Odd | multi-party turn-taking |
| [ ] `[S]` 4 | FizzBuzz multithreaded | predicate-driven waits, `Broadcast` |
| [ ] `[S]` 5 | Building H2O | barrier + grouping (2:1) |
| [ ] `[E]` 6 | Dining philosophers | deadlock: resource ordering, waiter/arbiter, Chandy–Misra |
| [ ] `[E]` 7 | Readers–writers (3 variants) | reader-pref, writer-pref, fair |
| [ ] `[T]` 8 | Sleeping barber | bounded waiting room, rendezvous |
| [ ] `[T]` 9 | Cigarette smokers | why semaphores alone aren't enough (pusher pattern) |
| [ ] `[T]` 10 | Santa Claus problem | priority between two groups |
| [ ] `[T]` 11 | Unisex bathroom / search-insert-delete | categorical mutual exclusion, starvation |
| [ ] `[T]` 12 | Traffic-light intersection (LeetCode 1279) | state machine under a lock |
| [ ] `[E]` 13 | Bounded blocking queue (LeetCode 1188) | see `05` |
| [ ] `[E]` 14 | Web crawler multithreaded (LeetCode 1242) | see `14` |
| [ ] `[S]` 15 | Deadlock-detect: build a lock wrapper that records lock order and detects cycles (a lockdep-lite) | lock-order graphs |
| [ ] `[E]` 16 | Bank transfer between accounts without deadlock | lock ordering by ID |

## Done when
All 16 pass, and for 6, 7, 15, 16 you can explain the deadlock/starvation argument formally (Coffman conditions).
