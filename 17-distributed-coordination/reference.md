# 17: Distributed coordination reference

## Concepts
- **No shared clock, no shared memory:** ordering comes from messages (Lamport happens-before, the same idea as `00` but across machines).
- **Leases:** time-bounded locks. They're only safe if clock drift is bounded and the holder checks its validity before acting. A GC pause or VM stall can outlive the lease.
- **Fencing tokens:** a monotonically increasing token issued with every lease grant. The resource rejects stale tokens. This is what makes a lease *safe*.
- **Redlock debate:** Kleppmann's critique (no fencing, timing assumptions) vs antirez's reply. Know both positions.
- **Exactly-once:** delivery is at-least-once. You get exactly-once *effects* from idempotent operations, dedupe keys, or transactions (outbox, inbox tables).
- **Consensus:** Raft and Paxos for a replicated log. Used by etcd and ZooKeeper for locks, leader election, and config.
- **CAP / PACELC, linearizability vs serializability**, read-your-writes, monotonic reads.
- **Optimistic concurrency control:** version columns / ETags + CAS-on-write (`UPDATE … WHERE version = ?`), vs pessimistic locks (`SELECT … FOR UPDATE`). At scale, prefer OCC with retries.

## Questions you must answer
1. Design a distributed lock. What breaks under GC pauses and clock skew? How do fencing tokens fix it?
2. How do you guarantee exactly-once processing of payments from a queue?
3. Two users edit the same record concurrently. Compare optimistic and pessimistic approaches.
4. Leader election: how does a deposed leader learn it's no longer leader before doing damage?
5. Lamport clocks vs vector clocks: what can each tell you?

## Reading
- antirez, *Is Redlock safe?* (the reply to Kleppmann below)
- DDIA ch. 5, 7, 8, 9
- Ongaro & Ousterhout, *In Search of an Understandable Consensus Algorithm* (Raft)
- MIT 6.5840 Distributed Systems labs (Go): https://pdos.csail.mit.edu/6.824/
- Lamport, *Time, Clocks, and the Ordering of Events in a Distributed System* (1978)
- Chris Richardson, *Transactional outbox* (microservices.io)

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**after `lease/`**
- Read `[S]` [Google: Chubby lock service](https://research.google/pubs/pub27897/): coarse-grained advisory locks; availability over performance.
- Code `[S]` [client-go `tools/leaderelection/leaderelection.go`](https://github.com/kubernetes/client-go/blob/master/tools/leaderelection/leaderelection.go): lease renewal with optimistic concurrency; clock skew can briefly yield two leaders.
- Code `[S]` [etcd `client/v3/concurrency/mutex.go`](https://github.com/etcd-io/etcd/blob/main/client/v3/concurrency/mutex.go): each waiter watches only its predecessor key, so there's no thundering herd.
- Code `[S]` [etcd `client/v3/concurrency/election.go`](https://github.com/etcd-io/etcd/blob/main/client/v3/concurrency/election.go): leader election on the same primitives.

**after `fencing/`**
- Read `[S]` [Kleppmann: How to do distributed locking](https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html): Redlock is unsafe under GC pauses and clock drift; fencing tokens are required.
- Code `[S]` [temporal `service/history/shard/ownership.go`](https://github.com/temporalio/temporal/blob/main/service/history/shard/ownership.go): rangeID fencing on every persistence write: fencing tokens at Uber/Temporal scale.

**after `idempotency/`**
- Read `[E]` [Stripe: APIs with idempotency](https://stripe.com/blog/idempotency): idempotency keys + retries with jitter; directly relevant to payments.
- Read `[E]` [AWS: Making retries safe with idempotent APIs](https://aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/): client request tokens, because two identical requests may both be intended.
- Read `[S]` [Airbnb: Avoiding Double Payments](https://medium.com/airbnb-engineering/avoiding-double-payments-in-a-distributed-payments-system-2981f6b070bb): idempotency split into pre-RPC / RPC / post-RPC phases.

**after `outbox/`**
- Read `[S]` [Confluent: Exactly-once in Kafka](https://www.confluent.io/blog/exactly-once-semantics-are-possible-heres-how-apache-kafka-does-it/): idempotent producers + transactions give exactly-once only *inside* Kafka.

**with `raft-lite/`**
- Code `[S]` [raft `raft.go`](https://github.com/hashicorp/raft/blob/main/raft.go): Raft as a run-loop actor per role.
- Code `[S]` [raft `rawnode.go`](https://github.com/etcd-io/raft/blob/main/rawnode.go): Raft as a pure state machine where the caller owns I/O; the opposite design.

**optional: gossip**
- Code `[T]` [memberlist `state.go`](https://github.com/hashicorp/memberlist/blob/master/state.go): SWIM gossip failure detection.
