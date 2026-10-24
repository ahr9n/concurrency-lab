# 14: Parallel algorithms reference

## Concepts
- **Amdahl's law:** speedup ≤ 1 / (s + (1−s)/N). The serial fraction s caps everything.
- **Gustafson's law:** with more cores you solve *bigger* problems, and scaled speedup ≈ N − s(N−1).
- **Work and span** (CLRS ch. 27): parallelism = work / span. Pick a grain size so per-task overhead stays far below the task's work.
- **Data parallelism vs task parallelism.** Static vs dynamic scheduling. Work stealing for irregular workloads.
- **Memory-bound vs compute-bound:** if you're bandwidth-bound, more cores won't help. Cache blocking and data layout (struct-of-arrays) matter more than goroutines.
- **Map-reduce:** map → (combine) → shuffle/partition by key → sort/group → reduce. Skewed keys create stragglers, so you need salting or splitting of hot keys. It maps directly onto Spark stages.
- **Deterministic output:** parallel reductions with floating point are non-deterministic unless you fix the reduction order.

## Hidden knowledge / gotchas
- Goroutine creation is about 1 µs. A task under about 10 µs needs batching.
- `GOMAXPROCS` defaults to the container CPU limit since 1.25. On a laptop it's all cores, including efficiency cores on Apple Silicon, which skews benchmarks.
- False sharing in per-worker result slots: pad them, or have each worker write to a local variable and assign it once at the end.
- Global dedupe in a crawler: check-and-insert must be atomic (`LoadOrStore`), or you crawl duplicates.

## Questions you must answer
1. Parallelise merge sort. Where do you stop splitting?
2. Implement word count map-reduce with N mappers and R reducers. How do you handle skew?
3. Build a concurrent web crawler with dedupe, depth limits, and per-host limits.
4. Why might your parallel version be slower than the sequential one?
5. Amdahl vs Gustafson: explain them with numbers.
6. Top-K over a 1 TB stream on one machine.

## Reading
- Dean & Ghemawat, *MapReduce* (OSDI 2004); Zaharia et al., *Resilient Distributed Datasets* (NSDI 2012)
- CLRS ch. 27 (multithreaded algorithms)
- Go tour: Exercise *Web Crawler*
- McKenney, *perfbook*, ch. 5–6 (partitioning, counting)

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**after `mapreduce/` (sharded workers)**
- Code `[S]` [datadog-agent `pkg/aggregator/time_sampler_worker.go`](https://github.com/DataDog/datadog-agent/blob/main/pkg/aggregator/time_sampler_worker.go): sharded workers fed over channels (a shuffle stage).
