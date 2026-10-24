# 06: Worker pool / executor

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** a production-quality executor with bounded queueing, rejection policies, panic isolation,
futures, and graceful shutdown.

## API
```go
type Pool struct{ /* ... */ }

func (p *Pool) Submit(ctx context.Context, task func(context.Context)) error // behaviour per RejectPolicy
func (p *Pool) Shutdown(ctx context.Context) error   // stop accepting, drain queued, wait (bounded by ctx)
func (p *Pool) ShutdownNow() []func(context.Context) // stop accepting, cancel running, return queued

// Go 1.27 allows generic methods on concrete types (not in interfaces):
func (p *Pool) Go[T any](ctx context.Context, f func(context.Context) (T, error)) *Future[T]
```

## Build
- [ ] `[E]` `fixed/`: N workers + your `05-bounded-queue`. Rejection policies: `Block`, `Reject` (`ErrQueueFull`), `CallerRuns`, `DropOldest`.
- [ ] `[E]` Panic isolation: a panicking task must not kill its worker or the process. Recover, report through a hook, continue.
- [ ] `[E]` `Future[T]`: `Get(ctx)`, `Done() <-chan struct{}`. Settable exactly once (a `sync.Once` or a closed-channel pattern).
- [ ] `[E]` Shutdown:
  - [ ] graceful (poison pills vs close-the-queue; justify your pick)
  - [ ] immediate (cancel a shared ctx, return what was queued)
- [ ] `[S]` `elastic/`: core + max workers, idle timeout, scale-up when the queue is non-empty (`ThreadPoolExecutor` semantics). Note the counter-intuitive "queue first, then grow" behaviour.
- [ ] `[E]` `boundedgo/`: the idiomatic minimal alternative. A semaphore channel + a goroutine per task (Bryan Mills's argument that "worker pools are often unnecessary in Go"). Benchmark it against `fixed/`.
- [ ] `[S]` Metrics: queued, active, completed, rejected, task latency histogram. Use atomics, without a lock on the hot path.

## Tests
- [ ] `[E]` 100k tasks: each runs exactly once, with no goroutine leaks after Shutdown (goleak).
- [ ] `[E]` Saturation: each reject policy behaves as specified.
- [ ] `[E]` Shutdown while tasks are being submitted concurrently: no panic (send on closed), no lost accepted task.
- [ ] `[E]` A panicking task doesn't reduce the worker count.

## Experiments
- [ ] `[S]` CPU-bound vs I/O-bound tasks: throughput vs worker count. Find the knee, and relate it to `GOMAXPROCS` and Little's Law.
- [ ] `[S]` `fixed/` vs `boundedgo/` vs `errgroup.SetLimit` vs unbounded `go f()` at 1M tiny tasks.

## Drill
A fixed pool with graceful shutdown and Futures, in under 20 minutes.
