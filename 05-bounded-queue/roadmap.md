# 05: Bounded blocking queue (producer–consumer)

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** the canonical interview primitive. You should be able to write it in your sleep, in three styles,
and explain every line.

## API (same for all implementations)
```go
type Queue[T any] interface {
    Put(ctx context.Context, v T) error   // blocks when full; returns ctx.Err() or ErrClosed
    Take(ctx context.Context) (T, error)  // blocks when empty; ErrClosed once closed AND drained
    TryPut(v T) bool
    TryTake() (T, bool)
    Close()                               // idempotent; wakes all waiters
    Len() int
}
```

## Build
- [ ] `[E]` `condq/`: `sync.Mutex` + **two** `sync.Cond` (notFull, notEmpty) + ring buffer. This is the "use locks and condition variables" version interviewers ask for. Make ctx work using `context.AfterFunc` + `Broadcast`.
- [ ] `[E]` `chanq/`: a buffered channel inside, plus a `closed` channel for shutdown. Explain why you can't just `close(ch)` while producers may be sending.
- [ ] `[T]` `semq/`: two semaphores (slots, items) + mutex around the slice. This is the textbook Dijkstra version.
- [ ] `[S]` `twolockq/`: linked-list queue with separate head/tail locks (the Michael–Scott blocking queue, like Java `LinkedBlockingQueue`). Explain the dummy node.
- [ ] `[S]` Optional: `PutBatch` / `TakeBatch` (drain up to N). This is what a log writer actually needs.

## Tests (shared suite run against all implementations)
- [ ] `[E]` M producers × N consumers, K items each. Every item is received exactly once (checksum + count).
- [ ] `[E]` Per-producer FIFO is preserved (sequence numbers).
- [ ] `[E]` `Put` blocks when full, `Take` blocks when empty. Verify with `synctest.Wait()`.
- [ ] `[E]` ctx cancellation unblocks waiters with no leaked goroutine.
- [ ] `[E]` Close: blocked producers get `ErrClosed`; consumers drain the remainder, then get `ErrClosed`.
- [ ] `[S]` Mutation tests: change `for` to `if` in a wait loop, or `Broadcast` to `Signal` in Close. Your suite must fail.

## Experiments
- [ ] `[S]` Throughput vs implementation vs capacity (1, 16, 1024) vs producers/consumers (1:1, 8:8, 64:1).
- [ ] `[E]` One `Cond` for both directions with `Signal`: reproduce the lost-wakeup / wrong-wakeup hang.

## Drill
`condq` from a blank file in **under 15 min**, including Close and ctx.

## Done when
You can answer every follow-up in `reference.md` without hesitation.
