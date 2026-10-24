# 18: Database concurrency control

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** most concurrency bugs in backend services live in the database, not in goroutines. Learn to see
anomalies, choose the right isolation and locking, and retry correctly, on MySQL/InnoDB and Spanner.

**Setup:** MySQL 8.4 in Docker (`docker run -e MYSQL_ROOT_PASSWORD=dev -p 3306:3306 mysql:8.4`), the
Spanner emulator (`gcr.io/cloud-spanner-emulator/emulator`), and optionally Postgres for SSI. Every exercise is a Go
test that drives two or more `*sql.Tx` from separate goroutines, stepped in a fixed order with channels.

## Build
- [ ] `[E]` `anomalies/`: reproduce dirty read, non-repeatable read, phantom, lost update, and write skew at READ COMMITTED / REPEATABLE READ / SERIALIZABLE. Produce a table of which level prevents what *in MySQL*, and compare it with Hermitage.
- [ ] `[E]` `lostupdate/`: an inventory decrement, four ways:
  - (a) read-modify-write (broken)
  - (b) atomic `UPDATE … SET qty = qty - 1 WHERE id = ? AND qty >= 1`
  - (c) `SELECT … FOR UPDATE`
  - (d) a version column + optimistic retry
  - Benchmark them under contention.
- [ ] `[E]` `deadlock/`:
  - Provoke an InnoDB deadlock: two transactions updating rows in opposite order.
  - Read `SHOW ENGINE INNODB STATUS`.
  - Handle error 1213 with retry + jitter.
  - Then remove the deadlock with consistent ordering.
- [ ] `[E]` `withtx/`: a `WithTx(ctx, db, fn)` helper: begin, run, commit; retry only retryable errors (deadlock 1213, lock-wait timeout 1205, serialization failures); bounded attempts with jitter; ctx-aware. Compare it with `crdb.ExecuteTx`.
- [ ] `[S]` `gaplocks/`: show next-key/gap locks under REPEATABLE READ blocking an insert into a range, and the difference at READ COMMITTED.
- [ ] `[S]` `skiplocked/`: a job queue on a table. N workers claim with `SELECT … FOR UPDATE SKIP LOCKED`, vs naive polling (duplicates, contention). This is the database version of `05`/`06`.
- [ ] `[S]` `spanner/`: on the emulator, run concurrent read-write transactions on one row. Observe aborts and automatic retries. Then prove the gotcha: **the transaction function can run more than once**, so a side effect inside it (publishing an event, sending an email) is a bug.
- [ ] `[S]` `longtx/`: one long-running transaction:
  - It holds connections, so `database/sql` pool exhaustion makes every request hang.
  - Its old snapshot grows the InnoDB history list.
  - Fix: ctx timeouts and short transactions.
- [ ] `[T]` `ssi/`: Postgres SERIALIZABLE (SSI). Show write skew aborting with a serialization failure where MySQL REPEATABLE READ lets it through.

## Tests
- [ ] Every anomaly test is deterministic: goroutines step through barriers (channels), never sleeps.
- [ ] `lostupdate`: 50 concurrent decrements on qty = 20 → exactly 20 succeed and qty ends at 0, for (b), (c), and (d).
- [ ] `withtx`: inject a deadlock error N times → succeeds on attempt N+1; a non-retryable error → no retry.

## Drill
Explain, then code, the "decrement stock without overselling" problem in 3 correct ways in under 20 minutes, with the trade-offs of each.

## Done when
You can say which anomalies your production database's default isolation allows, and you never write read-then-write without a lock, a condition, or a version.
