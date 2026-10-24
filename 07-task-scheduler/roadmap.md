# 07: Task scheduler (delayed, periodic, priority, DAG)

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** scheduling problems built on time, heaps, and dependency graphs. This covers the "job scheduler /
dependency-aware execution" family (Airflow / Temporal-style), plus classic delayed-execution questions.

## Build
- [ ] `[S]` `delayed/`: `Schedule(at time.Time, f)` / `Cancel(id)`. **One** goroutine plus a min-heap plus **one** `time.Timer` that is reset to the earliest deadline. Wake early when a new earliest task arrives. Cancellation is lazy (tombstones) or eager (heap index).
- [ ] `[S]` `periodic/`: fixed-rate vs fixed-delay semantics. Decide what happens when a run overruns the period (skip / queue / run concurrently).
- [ ] `[S]` `priority/`: a priority worker pool. Avoid starvation of low priority with aging or weighted fair queueing.
- [ ] `[S]` `dag/`: execute a DAG of tasks with dependencies:
  - [ ] cycle detection (Kahn / DFS colours)
  - [ ] max parallelism N; a task starts as soon as all its deps succeed (in-degree counters + ready queue)
  - [ ] failure policy: fail-fast (cancel the rest) vs continue-independent-branches
  - [ ] per-task retries with backoff, and timeouts
  - [ ] deterministic tests with `synctest`
- [ ] `[T]` `cron/`: a minimal cron-expression scheduler (optional; a parsing exercise).
- [ ] `[T]` `timerwheel/`: a hierarchical timing wheel for millions of timeouts. Compare with the heap at 1M timers.

## Tests
- [ ] `[S]` Delayed tasks fire in deadline order, within tolerance, using virtual time (`synctest`).
- [ ] `[S]` Cancel before fire → never runs. Cancel during run → defined behaviour.
- [ ] `[S]` DAG: dependency order respected under random task latencies, parallelism never exceeds N, a cycle returns an error.

## Drill
The DAG executor with a parallelism limit and fail-fast, in under 25 minutes.
