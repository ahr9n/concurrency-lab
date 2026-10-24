# 04: Context and lifecycle reference

## Concepts
- **`context.Context`** carries cancellation, a deadline, and request-scoped values across API boundaries. It is the first parameter, never stored in structs (except in rare lifecycle-owning types), and never nil.
- **Cancellation is cooperative:** context only *signals*. Your code must `select` on `ctx.Done()` or pass ctx down to a blocking call that honours it.
- **Propagation:** cancelling a parent cancels all its children. Children never cancel parents. Each `WithCancel` registers with the nearest cancelable ancestor, and you must `cancel()` to release that registration (`go vet` lostcancel).
- **Newer APIs:**
  - 1.20: `WithCancelCause`, `Cause`
  - 1.21: `WithDeadlineCause`, `WithTimeoutCause`, `AfterFunc`, `WithoutCancel`
  - 1.16: `signal.NotifyContext`
- **errgroup:**
  - The first error cancels the derived ctx.
  - `Wait` returns the first error.
  - `SetLimit(n)` bounds concurrency (it's a built-in semaphore). `TryGo` is the non-blocking variant.
- **Goroutine lifecycle rules:**
  1. Know when it will exit.
  2. Know who waits for it.
  3. Know how errors and panics surface.
  4. Never start one in a library without giving the caller a way to stop it.

## Hidden knowledge / gotchas
- `ctx.Value` lookup is O(depth). Use it only for request-scoped *metadata* (trace IDs, auth), never for dependencies or domain data.
- `context.WithoutCancel(ctx)` keeps the values but drops cancellation and deadline. It's right for "finish this audit write even if the client hung up".
- `context.AfterFunc(ctx, f)` runs f in its own goroutine when ctx is done. Its return value `stop()` tells you whether f was prevented from running.
- `http.Server.Shutdown` doesn't wait for hijacked connections (WebSockets). Track those yourself.
- A goroutine blocked in a syscall or non-ctx-aware I/O can't be cancelled by ctx. You need deadlines on the conn/file, or closing the fd.
- In errgroup, a goroutine that ignores ctx keeps `Wait` blocked even after the first error.
- `time.Sleep` doesn't honour ctx. Use `select { case <-time.After(d): case <-ctx.Done(): }` or a timer.

## Questions you must answer
1. How does cancellation propagate through a context tree? What leaks if you don't call cancel?
2. Design graceful shutdown for: HTTP + Pub/Sub consumer + batch flusher. In what order, and why?
3. How do you cancel a goroutine blocked on a network read?
4. When is `WithoutCancel` correct? When is it a bug?
5. How do you make `sync.Cond.Wait` respect a context?
6. What belongs in context values and what doesn't?

## Reading
- Go blog: *Go Concurrency Patterns: Context*; *Contexts and structs*
- https://pkg.go.dev/context (read every example) and `$(go env GOROOT)/src/context/context.go`
- *100 Go Mistakes*, #60–#62 (context mistakes)

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**after `errgroup/`**
- Code `[E]` [sync `errgroup/errgroup.go`](https://github.com/golang/sync/blob/master/errgroup/errgroup.go): first error cancels the shared ctx; `SetLimit` bounds parallelism.

**after `supervisor/` + `shutdown/`**
- Code `[S]` [cockroach `pkg/util/stop/stopper.go`](https://github.com/cockroachdb/cockroach/blob/master/pkg/util/stop/stopper.go): every async task registered, then quiesce → stop: reliable shutdown.
- Code `[S]` [tableflip `upgrader.go`](https://github.com/cloudflare/tableflip/blob/master/upgrader.go): zero-downtime restart by handing listener fds to a re-exec'd child.
