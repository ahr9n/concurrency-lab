# 18: Database concurrency reference

## Concepts
- **Anomalies:**
  - Dirty read
  - Non-repeatable read
  - Phantom
  - **Lost update:** two read-modify-writes, and one is overwritten.
  - **Write skew:** two transactions read overlapping data and write disjoint rows, breaking an invariant, e.g. "at least one doctor on call".
  - Read skew
- **Isolation levels are not standard across engines.** ANSI's definitions are ambiguous (Berenson et al.).
  - **MySQL InnoDB REPEATABLE READ** (the default): consistent-snapshot reads, but locking reads/writes see the *latest* data. Lost updates are possible with plain `SELECT` + `UPDATE`.
  - **Postgres REPEATABLE READ** errors on concurrent update instead.
  - **SERIALIZABLE:** MySQL turns plain SELECTs into locking reads. Postgres uses SSI (optimistic: abort on dangerous structures).
- **MVCC vs locking:**
  - MVCC gives readers snapshots without blocking writers.
  - Locking reads (`FOR UPDATE` / `FOR SHARE`) take row locks.
  - 2PL holds locks until commit.
- **InnoDB specifics:**
  - record, gap, and next-key locks (they prevent phantoms at RR)
  - insert-intention locks
  - automatic deadlock detection, which rolls back a victim (error 1213)
  - `innodb_lock_wait_timeout` (error 1205)
  - Locks are taken on *index records*, so a missing index means many more rows get locked.
- **Spanner:**
  - Read-write transactions use pessimistic locking with wound-wait deadlock avoidance, so a younger transaction may be aborted.
  - Read-only transactions and stale reads take no locks (snapshot at a timestamp).
  - TrueTime gives external consistency.
  - Aborted transactions are retried by the client library, **re-running your closure**.
- **Application-level patterns:**
  - atomic conditional `UPDATE`
  - optimistic concurrency (version / `updated_at` + `WHERE version = ?`, retry on 0 rows)
  - pessimistic `SELECT … FOR UPDATE`
  - `SKIP LOCKED` / `NOWAIT` for queues
  - unique constraints as the ultimate dedupe (idempotency keys, module 17)
  - advisory locks (`GET_LOCK`), used carefully
- **Retries:** retry the *whole transaction*, only on retryable errors, with jittered backoff and an attempt cap. Make the body idempotent, and keep side effects outside it (outbox, module 17).

## Hidden knowledge / gotchas
- ORMs love read-then-write. `obj = load(); obj.qty -= 1; save(obj)` is a lost update at every default isolation level that matters.
- `SELECT … FOR UPDATE` on a non-unique or unindexed predicate can lock far more than you think (gap locks), causing surprise deadlocks.
- Deadlocks are *normal* under concurrency. The bug is not retrying them, or not ordering lock acquisition (sort IDs before multi-row updates).
- A transaction open across a network call (an HTTP request to another service) holds locks and a connection for that whole time. This is the most common cause of DB-induced outages.
- `database/sql`: a `*sql.Tx` is bound to one connection. Forgetting `Rollback` on an error path leaks that connection until GC/timeout. Use `defer tx.Rollback()` (a no-op after Commit).
- Spanner: never do non-idempotent work inside `ReadWriteTransaction`'s function. Read-only workloads should use read-only or stale-read transactions to avoid lock contention entirely.
- Long snapshots in InnoDB stop purge from running, so the history list length grows and every query slows down.

## Questions you must answer
1. What anomalies does MySQL REPEATABLE READ allow? Postgres? Spanner?
2. Prevent overselling stock under concurrency: give three designs and their trade-offs.
3. Explain write skew with an example. Which isolation level prevents it, and how?
4. How do you handle deadlocks in application code? How do you reduce them?
5. Optimistic vs pessimistic concurrency control: when do you choose each?
6. Why can a Spanner transaction function run twice, and what does that imply?
7. Design a job queue on a relational table with many workers.
8. A long transaction is slowing the whole database. Why? How do you find it?

## Reading
- Berenson et al., *A Critique of ANSI SQL Isolation Levels* (1995): https://www.microsoft.com/en-us/research/publication/a-critique-of-ansi-sql-isolation-levels/
- Hermitage (isolation behaviour per database, tested): https://github.com/ept/hermitage
- MySQL 8.4 manual: [InnoDB Locking](https://dev.mysql.com/doc/refman/8.4/en/innodb-locking.html), [Isolation Levels](https://dev.mysql.com/doc/refman/8.4/en/innodb-transaction-isolation-levels.html), [Deadlocks](https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlocks.html)
- PostgreSQL: [Transaction Isolation](https://www.postgresql.org/docs/current/transaction-iso.html); Ports & Grittner, *Serializable Snapshot Isolation in PostgreSQL* (VLDB 2012): https://arxiv.org/abs/1208.4179
- Spanner docs: [Transactions](https://docs.cloud.google.com/spanner/docs/transactions), [TrueTime and external consistency](https://docs.cloud.google.com/spanner/docs/true-time-external-consistency)
- Corbett et al., *Spanner: Google's Globally-Distributed Database* (OSDI 2012): https://research.google/pubs/spanner-googles-globally-distributed-database-2/
- *Designing Data-Intensive Applications*, ch. 7 (Transactions)
- Jepsen analyses (how real databases violate their claims): https://jepsen.io/analyses

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**after `anomalies/`**
- Read `[S]` [Cockroach Labs: What write skew looks like](https://www.cockroachlabs.com/blog/what-write-skew-looks-like/): concrete write-skew cases, and why serializable isolation is worth its cost.

**after `lostupdate/` (c)**
- Read `[S]` [Cockroach Labs: When and why to use SELECT FOR UPDATE](https://www.cockroachlabs.com/blog/when-and-why-to-use-select-for-update-in-cockroachdb/): pessimistic locking to cut retry storms under contention, at the price of blocking.

**after `withtx/`**
- Code `[S]` [cockroach-go `crdb/tx.go`](https://github.com/cockroachdb/cockroach-go/blob/master/crdb/tx.go): `ExecuteTx`'s retry loop for serialization failures. Note how it keeps the body retry-safe.

**after `spanner/`**
- Code `[S]` [google-cloud-go `spanner/client.go`](https://github.com/googleapis/google-cloud-go/blob/main/spanner/client.go): `ReadWriteTransaction` wraps your function in an abort-retry loop. This is why the function must be idempotent.
