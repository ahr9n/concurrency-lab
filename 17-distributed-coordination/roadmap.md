# 17: Distributed coordination (concurrency beyond one process)

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** the same problems (mutual exclusion, ordering, exactly-once) across machines. Depth here pays off
in system-design rounds.

## Build (single-binary simulations with in-memory "nodes" and a fault-injecting network)
- [ ] `[S]` `lease/`: a lease-based leader election against a fake etcd-like store (CAS with TTL). Renew loop, lose the lease on a GC pause, and step down safely.
- [ ] `[S]` `fencing/`: fencing tokens. A storage service that rejects writes carrying an older token. Show the GC-pause/split-brain bug without fencing, then fix it.
- [ ] `[E]` `idempotency/`: an idempotency-key middleware with a state machine (in-progress / done / failed), concurrent duplicate requests, and TTL.
- [ ] `[S]` `outbox/`: the transactional outbox pattern with a concurrent relay (at-least-once publishing with dedupe downstream).
- [ ] `[T]` `lamport/` + `vectorclock/`: order events across simulated nodes and detect concurrent updates.
- [ ] `[T]` `raft-lite/` (stretch, big): leader election and log replication only, with randomised timeouts. Test it with a partitioning network. Or work through the MIT 6.5840 labs in Go.

## Done when
You can explain why "distributed lock with Redis" is not mutual exclusion without fencing, and design exactly-once *effects* on top of at-least-once delivery.
