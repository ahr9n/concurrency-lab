# 04: Context, cancellation, goroutine lifecycle

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** every goroutine has an owner, a way to stop, and a way to report errors.
Shutdown is ordered and loses no data.

## Build
- [ ] `[D]` `minictx/`: re-implement the core of `context` (`Background`, `WithCancel`, `WithDeadline`, `Value`). Use parent→child propagation via a children map under a mutex, with a lazily created `done` channel.
- [ ] `[E]` `errgroup/`: re-implement `errgroup.Group` with `Go`, `Wait`, `WithContext`, `SetLimit`, and `TryGo`. Compare with `x/sync/errgroup`.
- [ ] `[S]` `afterfunc/`: use `context.AfterFunc` to make a `sync.Cond` wait cancellable. This is the documented trick.
- [ ] `[S]` `supervisor/`: a long-running service with:
  - N background loops (ticker-driven)
  - restart-on-panic with backoff
  - stop via ctx
  - `Stop()` that waits for the loops to drain
- [ ] `[E]` `shutdown/`: an HTTP server plus background workers.
  - Use `signal.NotifyContext`, then:
    1. stop accepting
    2. `srv.Shutdown(ctx)` with a deadline
    3. drain queues
    4. flush state
    5. exit
  - Write the ordering down.
- [ ] `[S]` `causes/`: use `WithCancelCause`, `WithTimeoutCause`, `context.Cause`, `WithoutCancel` (for an audit log that must outlive the request).

## Experiments
- [ ] `[D]` Measure the overhead of a 10-deep context chain on `ctx.Done()` and `ctx.Value()` (Value is a linked-list walk).
- [ ] `[E]` Panic inside an errgroup goroutine. What happens? (Check your `x/sync` version.)

## Drill
A cancellable worker group with first-error-cancels-all and bounded parallelism, in under 15 minutes.

## Done when
You can design graceful shutdown for a service with HTTP, consumers, and background flushers, and defend the order.
