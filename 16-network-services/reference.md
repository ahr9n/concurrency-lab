# 16: Network services reference

## Concepts
- **Goroutine-per-connection** works in Go because of the netpoller (see `15`). Bound the *total* with a semaphore, and set deadlines on every conn.
- **Deadlines vs ctx:** `net.Conn` doesn't take a ctx. Use `SetDeadline`, or close the conn from a ctx watcher (`context.AfterFunc(ctx, func(){ conn.Close() })`).
- **`http.Server`:** each request runs in its own goroutine. `r.Context()` is cancelled when the client disconnects. `Shutdown` drains, `Close` is abrupt. `ReadHeaderTimeout` defends against slowloris.
- **`database/sql` pool:** `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`/`IdleTime`. Waiters queue in `connRequests`. A `Rows` you forget to close leaks a connection, and once the pool is exhausted every request hangs.
- **HTTP client:** `Transport` pools connections per host (`MaxIdleConnsPerHost` defaults to 2, a common throughput trap). You must read the body to EOF and close it, or the connection isn't reused.
- **Power of two choices:** pick 2 random backends and send to the one with less load. Exponentially better than random, and it avoids herding.
- **Batching/coalescing:** trades a small added latency for far fewer downstream calls. Results must route back to each caller (per-request result channels).
- **Head-of-line blocking:** one slow consumer on a shared writer stalls everyone. Give each connection its own writer goroutine and a bounded outbox.

## Questions you must answer
1. How does Go handle 100k concurrent connections without 100k OS threads?
2. Design a connection pool. How do waiters get connections fairly? What happens on ctx cancel?
3. How do you cancel a blocking network read?
4. Your service leaks goroutines slowly in production. How do you find it?
5. Implement request coalescing for a downstream batch API.
6. Pub/sub broker: what do you do with a slow subscriber?

## Reading
- `$(go env GOROOT)/src/database/sql/sql.go` (the pool)
- `$(go env GOROOT)/src/net/http/server.go` (`Serve`, `Shutdown`)
- Mitzenmacher, *The Power of Two Choices in Randomized Load Balancing*
- Graphql dataloader (the coalescing pattern)

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**after `tcpecho/`**
- Read `[S]` [Cloudflare: Go net/http timeouts guide](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/): which server/client timeout covers which phase; partly dated, but the diagram is the best one.

**after `connpool/`**
- Code `[S]` [puddle `pool.go`](https://github.com/jackc/puddle/blob/master/pool.go): generic resource pool behind pgxpool; read it before Vitess.
- Code `[S]` [vitess `go/pools/smartconnpool/pool.go`](https://github.com/vitessio/vitess/blob/main/go/pools/smartconnpool/pool.go): YouTube-scale MySQL pooling with lock-free per-setting stacks.
- Code `[S]` [vitess `go/pools/smartconnpool/waitlist.go`](https://github.com/vitessio/vitess/blob/main/go/pools/smartconnpool/waitlist.go): waitlist with direct connection handoff to waiters.

**after `coalesce/`**
- Read `[S]` [Discord: Storing Trillions of Messages (coalescing)](https://discord.com/blog/how-discord-stores-trillions-of-messages): request coalescing in a data-service layer protects hot partitions.
- Code `[S]` [kafka-go `writer.go`](https://github.com/segmentio/kafka-go/blob/main/writer.go): per-partition batching by size/time; Async mode hides errors.
