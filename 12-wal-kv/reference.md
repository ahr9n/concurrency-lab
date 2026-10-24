# 12: WAL KV reference

## Concepts
- **Write path:** WAL append (durable) → memtable insert → ack. Recovery = replay the WAL into the memtable.
- **LSM tree:** memtable → immutable memtable → L0 SSTables → leveled/tiered compaction. Write-optimised. Read amplification is paid down by compaction and bloom filters.
- **Amplification trade-offs:** write vs read vs space amplification (the RUM conjecture). Leveled compaction has low read/space and high write amplification. Tiered is the opposite.
- **Manifest / version set:** the authoritative list of live files. Update it atomically: write a new manifest, fsync, rename, fsync the dir. `rename` is atomic on POSIX filesystems.
- **Crash consistency ordering:** data durable *before* metadata references it, and metadata durable *before* the old data is deleted.
- **Concurrency structure:**
  - single writer for the WAL (group commit)
  - a concurrent memtable (sharded locks, or a lock-free skiplist as in RocksDB/Pebble)
  - immutable segments (no locks needed to read them)
  - refcounted versions for safe file deletion
  - background goroutines for flush and compaction, with bounded concurrency
- **MVCC / snapshots:** keys carry seq numbers. Readers pick a seq. Compaction must keep the versions visible to live snapshots.
- **Write stalls:** backpressure when background work falls behind. Without it, memory grows or reads degrade without limit.

## Real systems to read (Go first)
- **Badger** (Dgraph, Go, WiscKey-style key-value separation)
- **bbolt** (B+tree, single writer, mmap): a contrast to LSM
- **etcd** WAL + snapshot + raft storage
- LevelDB / RocksDB wiki (compaction, write stalls)
- Bitcask paper (the simplest log-structured hash table)

## Questions you must answer
1. Walk me through a write, from API call to durable ack.
2. What happens on a crash during memtable flush? During compaction?
3. How do reads proceed while compaction deletes files?
4. LSM vs B-tree: when would you choose each?
5. How do you implement snapshots?
6. Where are the locks in your system, and why is each one there? What is lock-free and why?
7. How would you scale writes beyond one WAL? What breaks (cross-shard atomicity)?

## Reading
- O'Neil et al., *The Log-Structured Merge-Tree* (1996)
- Lu et al., *WiscKey* (FAST 2016)
- Sheehy & Smith, *Bitcask* (2010)
- Pebble docs: `docs/` folder
- *Database Internals*, part I; DDIA ch. 3

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**before M1**
- Read `[S]` [CockroachDB: Introducing Pebble](https://www.cockroachlabs.com/blog/pebble-rocksdb-kv-store/): rewrote RocksDB in Go to escape cgo, at the cost of owning a storage engine.

**after M1 (recovery)**
- Code `[S]` [loki `pkg/ingester/wal.go`](https://github.com/grafana/loki/blob/main/pkg/ingester/wal.go): WAL for an ingester.
- Code `[S]` [loki `pkg/ingester/replay_controller.go`](https://github.com/grafana/loki/blob/main/pkg/ingester/replay_controller.go): replay with memory backpressure.

**after M3 (versions)**
- Code `[S]` [pebble `version_set.go`](https://github.com/cockroachdb/pebble/blob/master/version_set.go): refcounted immutable versions + manifest (your M3).

**after M5 (snapshots / MVCC)**
- Code `[S]` [badger `txn.go`](https://github.com/dgraph-io/badger/blob/main/txn.go): MVCC oracle with commit-time conflict detection (your M5).

**contrast: B-tree, single writer**
- Code `[S]` [bbolt `db.go`](https://github.com/etcd-io/bbolt/blob/main/db.go): single writer, many mmap readers; long readers block page reuse.
- Code `[S]` [bbolt `tx.go`](https://github.com/etcd-io/bbolt/blob/main/tx.go): transaction lifecycle in a B+tree store.

**after M6 (backpressure)**
- Read `[S]` [AWS: DynamoDB (ATC 2022)](https://www.usenix.org/conference/atc22/presentation/vig): per-partition and global admission control for predictable latency.
