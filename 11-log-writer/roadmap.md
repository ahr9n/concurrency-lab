# 11: Synchronous log writer (the flagship problem)

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Prompt to solve:** "Many application goroutines call `Log(record)`. Write all records to a single file.
When `Log` returns, the record must survive a crash. Maximise throughput: don't fsync per record."
Common variants: an event writer shared by thousands of threads; `push(bytes)` + async push + wait-for-durability;
rotating files without stalling writers; proving with a crash test that no acknowledged record is lost.
Expect constraints to be added mid-problem.

## API
```go
type Writer interface {
    Append(ctx context.Context, rec []byte) (lsn uint64, err error) // returns only after durable
    AppendAsync(rec []byte) (lsn uint64, err error)                 // durable later; caller may WaitDurable
    WaitDurable(ctx context.Context, lsn uint64) error
    Close() error                                                   // flush everything, then stop
}
```

## Build in versions (keep every version; the benchmarks tell the story)
- [ ] `[E]` **v0 `naive/`:** one mutex → `write` → `fsync` → unlock. Correct and slow. This is your baseline.
- [ ] `[E]` **v1 `writerloop/`:** callers send `{rec, done chan error}` on a channel. A single writer goroutine owns the file, writes, fsyncs, and replies. (Channel ownership instead of a lock.)
- [ ] `[E]` **v2 `groupcommit/`:** the writer goroutine **drains everything queued** (batch) → one write → one fsync → completes the whole batch. Throughput now scales with concurrency.
- [ ] `[E]` **v3 `doublebuffer/`:** a mutex-protected active buffer.
  - Appenders copy into it, take an LSN, then wait on a `Cond` or a per-LSN notification for `durableLSN ≥ mine`.
  - A flusher swaps the active and standby buffers *under the lock*, then writes and fsyncs *outside* the lock, advances `durableLSN`, and broadcasts.
  - Appends proceed during fsync. **This is the answer interviewers want to hear explained.**
- [ ] `[S]` **v4 `leader/`:** leader-follower group commit, with no dedicated flusher. The first waiter becomes the leader and flushes for everyone (Postgres / MySQL binlog style). Add an optional `commitDelay` that waits up to D µs if ≥ K other appenders are active (Postgres `commit_delay`/`commit_siblings`).
- [ ] `[E]` Record framing: `[len u32][crc32c u32][payload]`. On open: scan, verify CRC, **truncate the torn tail**, resume LSN.
- [ ] `[E]` Backpressure: a bounded buffer. When full, `Append` blocks (ctx-aware) instead of growing without limit.
- [ ] `[S]` Error semantics: if fsync fails, the writer becomes *poisoned* and every waiter gets the error. Never retry an fsync on the same fd (fsyncgate).
- [ ] `[S]` Segment rotation without stalling appenders: pre-create the next segment in the background, swap it at a batch boundary, fsync the new file **and its directory**.

## Tests
- [ ] `[E]` Every `Append` that returned nil is present after reopening. Use LSNs to check gaps and order.
- [ ] `[E]` **Crash test:** run the writer in a subprocess (`exec.Command(os.Args[0])` with a test env var), `SIGKILL` it at a random time, reopen, verify that all acknowledged LSNs exist and the tail is clean.
- [ ] `[S]` Torn-write test: truncate the file mid-record, or flip a byte, and check that recovery stops at the last valid record.
- [ ] `[S]` fsync failure injection behind a `File` interface. Every pending and future appender gets the error.
- [ ] `[E]` `-race`, goleak on Close.

## Benchmarks (put the table in NOTES.md)
- [ ] `[E]` Records/sec and p50/p99 Append latency for v0–v4 at 1/8/64/256 concurrent appenders and 100 B / 4 KB records.
- [ ] `[E]` Count fsyncs per second and records per fsync (batch size). The plot of v0 vs v3 is your story.
- [ ] `[S]` Measure the fsync cost on your Mac (APFS `F_FULLFSYNC`) vs a Linux VM/container (`fdatasync`).

## Drill
v3 from a blank file (no framing, no recovery) in **under 25 minutes**, while explaining the happens-before edges out loud.
