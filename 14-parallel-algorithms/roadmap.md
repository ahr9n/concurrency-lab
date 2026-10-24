# 14: Parallel algorithms and data processing

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** CPU parallelism. Divide and conquer, data parallelism, and knowing when parallelism doesn't help.
This is relevant to a data-platform company (Spark-like thinking at goroutine scale).

## Build
- [ ] `[S]` `mergesort/`: parallel merge sort with a depth/size cutoff to sequential. Measure speedup vs `slices.Sort`.
- [ ] `[T]` `prefixsum/`: parallel scan (Blelloch up-sweep/down-sweep or a two-pass block scan).
- [ ] `[S]` `mapreduce/`: word count over a directory of files.
  - [ ] map stage parallel per file (bounded)
  - [ ] **shuffle** by hash(key) into R reducers via channels
  - [ ] reduce stage
  - [ ] combiner optimisation
  - [ ] compare with a sharded concurrent map
- [ ] `[E]` `crawler/`: concurrent web crawler (Go tour exercise++): depth limit, global dedupe (sync.Map or mutex-set), per-host concurrency limit, politeness delay, ctx cancellation, no leaks. Test against `httptest.Server`.
- [ ] `[T]` `matmul/`: blocked matrix multiply. Show cache effects dominating parallelism.
- [ ] `[S]` `parallelfor/`: a generic `ParallelFor(n, grain, fn)` helper with static vs dynamic (work-queue) chunking. Measure it on skewed workloads.
- [ ] `[S]` `externalsort/`: sort a file 10× bigger than your memory budget. Parallel run generation + k-way merge with a heap.
- [ ] `[S]` `topk/`: top-K frequent items over a stream: per-shard heaps merged at the end; Count-Min Sketch for approximate.

## Experiments
- [ ] `[S]` Speedup curve at `GOMAXPROCS` = 1..N for each. Compare with Amdahl's prediction given the measured serial fraction.
- [ ] `[S]` Too fine-grained tasks: show goroutine overhead beating the gains (no cutoff vs a cutoff).

## Done when
You can predict, before measuring, roughly whether parallelising a given loop will help, and explain why.
