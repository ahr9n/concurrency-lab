# 19: Verification and deterministic simulation

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** go beyond "the stress test passed". Use a model checker to explore *every* interleaving of a small design,
and deterministic simulation to replay rare failures exactly. Do this after `05`, `11`, and `12`, because you'll verify your own code.

**Setup:** TLA+ Toolbox or the VS Code TLA+ extension (TLC model checker); Go 1.27.

## Build
- [ ] `[S]` `tla/boundedqueue/`:
  - Following the BlockingQueue tutorial, spec your `05` queue with **one** condition variable and `Signal`.
  - Let TLC find the deadlock trace.
  - Fix it (two conditions, or `Broadcast`) and re-check.
- [ ] `[S]` `tla/groupcommit/`:
  - Model the `11` log writer: appenders, buffer swap, flush, ack.
  - Invariant: *acked ⇒ durable*.
  - Introduce the bug "ack before fsync" and watch TLC produce the counterexample.
- [ ] `[T]` `tla/transfer/`: model `13`'s bank transfer and check deadlock freedom with and without lock ordering. Add the liveness property "every transfer eventually completes".
- [ ] `[S]` `dst/`: a tiny deterministic simulator for your `12` WAL-KV or `17` lease election:
  - [ ] inject all nondeterminism behind interfaces: `Clock`, `Rand`, `Disk` (can crash, tear, or fail fsync), `Network` (drop, delay, reorder, partition)
  - [ ] run "nodes" as state machines stepped by a single-threaded, seed-driven scheduler (no free-running goroutines)
  - [ ] the same seed → the same execution. A failing seed is a bug report you can replay.
  - [ ] fault injection (BUGGIFY-style): rare branches forced on randomly per seed
  - [ ] run 10k seeds in CI-sized batches, and check invariants after every step
- [ ] `[S]` `synctest-limits/`: write down what `testing/synctest` makes deterministic (time, durable blocking) and what it doesn't (scheduling order, mutexes, real I/O). When is it "DST-lite" enough?
- [ ] `[T]` `fuzz/`: Go native fuzzing of operation sequences against a single-threaded model (as in `12`'s model test), to find the shortest failing history.
- [ ] `[T]` Read etcd's raft TLA+ spec and its trace validation: checking traces from the real implementation against the spec.

## Done when
You can explain model checking vs testing vs simulation, show one bug TLC found in *your* design, and replay a simulated crash from a seed.
