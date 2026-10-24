# 11: Log writer reference

## The core insight
**Don't hold the lock across I/O.** The lock protects only the in-memory buffer and the LSN counters.
Swap the buffers under the lock, then do the slow write+fsync without it. One fsync makes a whole batch
durable (**group commit**), so throughput grows with concurrency instead of collapsing.

## Invariants (say these out loud)
- LSNs are assigned in append order under the lock. Bytes land in the file in LSN order.
- `durableLSN` only increases. It's updated only after fsync returns successfully.
- `Append` returns success ⇔ `durableLSN ≥ myLSN`.
- At most one flush is in flight at a time (or flushes are ordered), so the file is written sequentially.
- After a fsync error the writer is permanently failed (poisoned). No later success may be reported.

## Happens-before edges in v3
1. The appender's copy into `active` happens before its Unlock. Unlock happens before the flusher's Lock. So the flusher sees the bytes.
2. The flusher's `durableLSN = x` (under the lock or atomic) happens before Broadcast/Unlock. Unlock happens before the waiter's re-Lock. So the waiter sees `durableLSN`.
3. The swapped-out buffer is touched only by the flusher until it's handed back. Ownership transfers at the swap, under the lock.

## Durability facts (where people get it wrong)
- `write()` → OS page cache only. `bufio.Writer.Flush` → only into the kernel. **Neither is durable.**
- `(*os.File).Sync()`: on Linux it calls `fsync`; **on macOS Go uses `fcntl(F_FULLFSYNC)`** (verified in `src/internal/poll/fd_fsync_darwin.go`). That's why fsync on a Mac is slow and honest.
- **A new file needs its directory fsynced too**, or the file's directory entry can vanish after a crash.
- `fdatasync` skips the metadata flush when size is unchanged. Preallocating (`fallocate`) + `fdatasync` is the classic WAL optimisation.
- **fsyncgate (PostgreSQL 2018):** on Linux, after a failed fsync the dirty pages may be dropped and marked clean, so a *retry can succeed while the data is lost*. The only safe response is to crash and recover from the log.
- O_DIRECT / O_DSYNC are alternatives with their own trade-offs (alignment, per-write sync).
- A torn write at the tail after a crash is normal. Framing + CRC + truncating to the last valid record is the standard recovery.

## Design alternatives to compare
| Design | Lock across fsync? | Batching | Notes |
|---|---|---|---|
| mutex + fsync | yes | none | correct, fsync-bound (~1/fsync_latency ops/s) |
| single writer goroutine | n/a (owner) | optional | simplest correct concurrent design in Go |
| group commit (drain queue) | no | natural | the standard answer |
| double buffering | no | natural | appends overlap with the fsync |
| leader/follower | no | natural | no dedicated thread; Postgres/MySQL style |
| async + periodic fsync | no | time-based | *not* synchronous durability; bounded loss window |

## Production references
- PostgreSQL WAL: `commit_delay` / `commit_siblings`, `wal_sync_method`
- MySQL InnoDB / binlog group commit (`binlog_group_commit_sync_delay`)
- LevelDB/RocksDB log format (32 KB blocks, record types FULL/FIRST/MIDDLE/LAST)
- Kafka: relies on page cache + replication instead of fsync per message. A different durability model.
- Log4j2 AsyncLogger / LMAX Disruptor (async, not durable-on-return)

## Follow-ups interviewers ask
1. Why not fsync under the lock? Quantify it: fsync ≈ 0.1–10 ms, so ≤ 10k ops/s total regardless of thread count.
2. How do you guarantee ordering across the buffer swap?
3. What if fsync fails? What if it hangs?
4. What exactly survives a power loss vs a process crash? (Process crash: the page cache survives. Power loss: only fsynced data.)
5. How do you recover? How do you detect a partially written record?
6. How do you bound memory if the disk is slow? (Backpressure.)
7. Latency vs throughput: does group commit increase latency for a single caller? (Slightly, if there's a commit delay.)
8. How would you shard the log for more throughput, and what ordering do you lose?
9. How do readers read the log concurrently with writers (tailing up to durableLSN)?
10. Per-thread ordering vs global total ordering: which did you guarantee, and how?

## Reading
- Dan Luu, *Files are hard*: https://danluu.com/file-consistency/
- Pillai et al., *All File Systems Are Not Created Equal* (OSDI 2014)
- PostgreSQL docs: *WAL Configuration* (commit_delay section)
- LevelDB `doc/log_format.md`
- *Database Internals* ch. 5 (recovery); DDIA ch. 3 and 7

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**after v2/v3**
- Read `[S]` [LinkedIn (Jay Kreps): The Log](https://engineering.linkedin.com/distributed-systems/log-what-every-software-engineer-should-know-about-real-time-datas-unifying): the append-only log as the source of truth; the philosophy behind 11 and 12.
- Code `[E]` [pebble `record/log_writer.go`](https://github.com/cockroachdb/pebble/blob/master/record/log_writer.go): flusher goroutine + group commit + pending-sync queue: the interview problem in production.

**after v4 (leader/follower)**
- Code `[S]` [pebble `commit.go`](https://github.com/cockroachdb/pebble/blob/master/commit.go): commit pipeline: concurrent apply, in-order publish; the top doc comment is gold.

**after record framing + fsync errors**
- Read `[E]` [LWN: PostgreSQL's fsync() surprise](https://lwn.net/Articles/752063/): after a failed fsync a retry can "succeed" while the data is lost. Never retry fsync.
- Read `[S]` [Can Applications Recover from fsync Failures? (ATC 2020)](https://www.usenix.org/conference/atc20/presentation/rebello): PostgreSQL, LevelDB, LMDB, SQLite, and Redis all mishandle fsync errors.
- Code `[E]` [etcd `server/storage/wal/wal.go`](https://github.com/etcd-io/etcd/blob/main/server/storage/wal/wal.go): CRC-framed WAL; recovery and torn-tail handling.
- Code `[E]` [etcd `server/storage/wal/encoder.go`](https://github.com/etcd-io/etcd/blob/main/server/storage/wal/encoder.go): the record framing itself.

**after segment rotation**
- Code `[S]` [etcd `server/storage/wal/file_pipeline.go`](https://github.com/etcd-io/etcd/blob/main/server/storage/wal/file_pipeline.go): pre-allocates the next segment in the background, so rotation never stalls.

**for the "tail the log" follow-up**
- Code `[S]` [prometheus `tsdb/wlog/watcher.go`](https://github.com/prometheus/prometheus/blob/main/tsdb/wlog/watcher.go): a tailing reader that follows a live WAL.

**contrast: the opposite (lossy) choice**
- Code `[S]` [zap `zapcore/buffered_write_syncer.go`](https://github.com/uber-go/zap/blob/master/zapcore/buffered_write_syncer.go): timer-flushed buffering: the opposite, lossy-on-crash choice.
