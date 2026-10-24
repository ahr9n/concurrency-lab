# 09: Flow control reference

## Concepts
- **Token bucket:** a capacity B (burst) refilled at rate r.
  - Lazy refill: `tokens = min(B, tokens + r·(now−last))`, computed on acquire.
  - The **check-then-act must be atomic** (one lock or one CAS).
- **GCRA:** stores one timestamp (the theoretical arrival time) instead of tokens+time, so a single-word CAS is enough. It's mathematically equivalent to a leaky bucket.
- **Leaky bucket** (as a meter or as a queue) vs token bucket: the leaky bucket smooths output, the token bucket allows bursts.
- **Window algorithms:**
  - Fixed window: allows 2× bursts at window edges.
  - Sliding log: exact, but O(n) memory.
  - Sliding window counter: a weighted pair of windows. Cheap and approximately right.
- **Distributed rate limiting:** a central store (Redis + Lua for atomicity), or local buckets with periodically synced quotas. Databricks moved to in-memory token buckets with batch-reported usage to cut tail latency. Read their blog.
- **Concurrency limits vs rate limits:** by Little's Law, in-flight = rate × latency. Limiting concurrency adapts automatically when latency rises. Rate limits don't.
- **Circuit breaker:** stops calling a failing dependency. Half-open lets a few probes through. Use a rolling window (ring of buckets).
- **Retries:**
  - Exponential backoff + **full jitter** (`sleep = rand(0, min(cap, base·2^n))`) prevents synchronized retry waves.
  - Retry budgets (e.g. ≤10% extra load) prevent retry amplification across layers.
  - Never retry non-idempotent operations blindly.
- **Hedging:** cuts tail latency at a small extra load. Requires cancellation and idempotency. See *The Tail at Scale*.
- **Load shedding:** reject early and cheaply. Under overload, prefer LIFO or newest-first processing, because the oldest requests have probably already timed out on the client side.

## Hidden knowledge / gotchas
- Use **monotonic time** for refill. In Go, `time.Since(start)` is monotonic if `start` came from `time.Now()`. Never use a Unix timestamp difference.
- A ticker goroutine refilling tokens wastes wakeups and is less precise than lazy refill.
- A CAS loop under heavy contention can be *slower* than a mutex (retry storms, cache-line ping-pong).
- `x/time/rate.Limiter.Wait` reserves tokens *before* sleeping. If ctx is cancelled it returns them. Read its `Reservation` design.
- Per-key limiter maps leak memory without eviction of idle keys.

## Questions you must answer
1. Implement a thread-safe token bucket. Now make it lock-free.
2. Fixed vs sliding window: show the edge-burst problem.
3. Design rate limiting for 10k tenants across 100 servers.
4. Rate limit vs concurrency limit: which protects better against a slow dependency? Why?
5. Explain backoff with jitter. Why "full" jitter?
6. How do retries cause cascading failures? How do budgets fix that?
7. Circuit breaker states and transitions. How do you avoid a race on half-open?

## Reading

- https://github.com/Netflix/concurrency-limits (README)
- `golang.org/x/time/rate` source

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**after `tokenbucket/`**
- Read `[E]` [Databricks: High Performance Ratelimiting](https://www.databricks.com/blog/high-performance-ratelimiting-databricks): Redis → in-memory sharded token buckets with batched reporting, ~10× better tail latency.
- Read `[E]` [Stripe: Scaling your API with rate limiters](https://stripe.com/blog/rate-limiters): four limiter types; low-priority traffic is shed first.
- Code `[E]` [time `rate/rate.go`](https://github.com/golang/time/blob/master/rate/rate.go): lazy token bucket with a Reservation model, no goroutine.
- Code `[S]` [ratelimit `limiter_atomic_int64.go`](https://github.com/uber-go/ratelimit/blob/main/limiter_atomic_int64.go): CAS leaky bucket; read it next to the mutex version.
- Code `[S]` [ratelimit `limiter_mutexbased.go`](https://github.com/uber-go/ratelimit/blob/main/limiter_mutexbased.go): the mutex version of the same limiter.

**after `slidingwindow/` + `keyed/`**
- Read `[S]` [Cloudflare: rate limiting for millions of domains](https://blog.cloudflare.com/counting-things-a-lot-of-different-things/): approximate, PoP-local sliding-window counters instead of global accuracy.
- Read `[S]` [AWS: Shuffle-sharding](https://aws.amazon.com/builders-library/workload-isolation-using-shuffle-sharding/): spend extra capacity so one bad tenant can only hurt a few others.
- Code `[S]` [dskit `limiter/rate_limiter.go`](https://github.com/grafana/dskit/blob/main/limiter/rate_limiter.go): a global limit split across ring members (the multi-server tenant limit).

**after `concurrencylimit/`**
- Read `[S]` [Netflix: Performance Under Load](https://netflixtechblog.medium.com/performance-under-load-3e6fa9a60581): adaptive concurrency limits modelled on TCP congestion control.
- Read `[S]` [Uber: Cinnamon load shedder](https://www.uber.com/blog/cinnamon-using-century-old-tech-to-build-a-mean-load-shedder/): TCP-Vegas + PID control: load shedding with no per-service tuning.
- Code `[S]` [VictoriaMetrics `lib/writeconcurrencylimiter/concurrencylimiter.go`](https://github.com/VictoriaMetrics/VictoriaMetrics/blob/master/lib/writeconcurrencylimiter/concurrencylimiter.go): bounded ingest concurrency in a high-throughput TSDB.

**after `breaker/`**
- Code `[S]` [datadog-agent `comp/forwarder/defaultforwarder/impl/blocked_endpoints.go`](https://github.com/DataDog/datadog-agent/blob/main/comp/forwarder/defaultforwarder/impl/blocked_endpoints.go): per-endpoint circuit breaker in Datadog's forwarder.
- Code `[S]` [mimir `pkg/ingester/circuitbreaker.go`](https://github.com/grafana/mimir/blob/main/pkg/ingester/circuitbreaker.go): circuit breaker on ingester push/read paths.

**after `retry/`**
- Read `[E]` [AWS: Timeouts, retries and backoff with jitter](https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/): retries amplify overload; jitter and budgets limit the damage.
- Code `[E]` [client-go `util/workqueue/default_rate_limiters.go`](https://github.com/kubernetes/client-go/blob/master/util/workqueue/default_rate_limiters.go): per-item exponential backoff max'd with a global bucket; you must call `Forget`.
- Code `[S]` [client-go `util/flowcontrol/backoff.go`](https://github.com/kubernetes/client-go/blob/master/util/flowcontrol/backoff.go): backoff bookkeeping per key, with garbage collection.

**after `hedge/`**
- Read `[S]` [Google: The Tail at Scale](https://research.google/pubs/pub40801/): hedged requests: a little extra load for a much better p99.
- Code `[S]` [dskit `ring/replication_set.go`](https://github.com/grafana/dskit/blob/main/ring/replication_set.go) (`DoUntilQuorum`): hedged quorum fan-out that cancels the rest once quorum is reached.

**after `shed/`**
- Read `[S]` [Netflix: Prioritized Load Shedding](https://netflixtechblog.com/keeping-netflix-reliable-using-prioritized-load-shedding-6cc827b02f94): shed by priority so playback survives while extras degrade.
- Read `[S]` [LinkedIn: Hodor overload detection](https://engineering.linkedin.com/blog/2022/hodor--detecting-and-addressing-overload-in-linkedin-microservic): threshold-free overload detection that drops just enough traffic.
- Read `[S]` [Facebook: Fail at Scale](https://cacm.acm.org/practice/fail-at-scale/): CoDel + adaptive LIFO: serve the newest requests first under backlog.
- Read `[S]` [AWS: Using load shedding to avoid overload](https://aws.amazon.com/builders-library/using-load-shedding-to-avoid-overload/): reject early and cheaply; work that later times out is pure waste.
- Read `[S]` [Google SRE: Handling Overload](https://sre.google/sre-book/handling-overload/): client-side throttling and per-request criticality.
- Read `[S]` [Google SRE: Addressing Cascading Failures](https://sre.google/sre-book/addressing-cascading-failures/): how overload spreads, and the levers that stop it.

**last: the big one**
- Read `[S]` [CockroachDB: Admission Control](https://www.cockroachlabs.com/blog/admission-control-in-cockroachdb/): Go's scheduler can't prioritise goroutines, so they built token/slot queues with epoch-LIFO.
- Code `[S]` [cockroach `pkg/util/admission/`](https://github.com/cockroachdb/cockroach/tree/master/pkg/util/admission): the best real-world admission/load-shedding code in Go (`work_queue.go`, `granter.go`).
