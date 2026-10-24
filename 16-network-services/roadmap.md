# 16: Concurrency in network services

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** the concurrency you meet in real backends: per-connection goroutines, pools, deadlines,
streaming, and coordination across requests.

## Build
- [ ] `[E]` `tcpecho/`: a TCP server, one goroutine per connection, with read/write deadlines, a max-connections semaphore, and graceful shutdown (stop the listener, drain the conns, force-close after a deadline).
- [ ] `[S]` `connpool/`: a generic connection pool. Max open/idle, idle timeout, health check on borrow, ctx-aware `Get`, and a fair waiter queue. Compare its design with `database/sql`'s `DB` pool (read `sql.go` `conn()`/`putConn`).
- [ ] `[S]` `broker/`: a TCP pub/sub broker (topics, fan-out, slow-consumer policy, per-connection write goroutine plus a bounded outbox).
- [ ] `[S]` `coalesce/`: request coalescing and batching in an HTTP handler (a dataloader pattern). Many requests per 5 ms window become one downstream batch call, and each caller gets its own slice of the result.
- [ ] `[S]` `streaming/`: a server-sent-events / gRPC streaming fan-out with backpressure. Detect client disconnect via `r.Context()`.
- [ ] `[S]` `lb/`: a client-side load balancer (round-robin, least-in-flight, power-of-two-choices) with atomic in-flight counters.

## Tests
- [ ] `[S]` Load with `httptest` / a custom client at 10k concurrent connections. No goroutine growth after the clients disconnect.
- [ ] `[S]` Pool: under-limit, at-limit (wait), ctx cancel while waiting, broken connection eviction.

## Experiments
- [ ] `[S]` Slowloris: a client that trickles bytes. Show the server running out of resources without read deadlines, then fix it.
- [ ] `[S]` P2C vs round-robin p99 when one backend is slow.
