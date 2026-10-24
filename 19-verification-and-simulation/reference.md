# 19: Verification and simulation reference

## Concepts
- **Testing samples interleavings. Model checking enumerates them** (for a bounded model). Concurrency bugs hide in the interleavings that stress tests rarely hit.
- **TLA+ / PlusCal:**
  - You describe states and allowed steps.
  - **Invariants** are safety properties: nothing bad happens.
  - **Temporal properties** cover liveness: something good eventually happens. Fairness assumptions are required for liveness.
  - TLC checks them by exhaustive state exploration, which suffers from state explosion, so keep models small (2–3 processes, small buffers).
- **What to model:** the *protocol* (locks, buffers, messages, crash points), not the code. Bugs in the design are cheapest to find before code exists.
- **Trace validation:** record events from the real implementation and check that they're a valid behaviour of the spec. This closes the spec/code gap.
- **Deterministic simulation testing (DST):**
  - Run the whole system single-threaded under a seeded scheduler.
  - Every source of nondeterminism (time, randomness, thread scheduling, disk, network) is simulated and injected.
  - Faults are injected aggressively.
  - Any failure replays exactly from its seed.
  - This is what lets FoundationDB and TigerBeetle find bugs in hours of simulated time instead of months in production.
- **Why Go makes DST harder:** the goroutine scheduler and `select` are nondeterministic, and `GOMAXPROCS=1` doesn't make them deterministic. Options:
  - structure the code as state machines driven by your own scheduler (like `etcd-io/raft`)
  - run under WASM with a controlled runtime (the Polar Signals approach)
  - use a deterministic hypervisor (Antithesis)
- **Linearizability checking** (porcupine, modules 10/12) verifies *histories* from real runs. It's complementary to model checking.

## Hidden knowledge / gotchas
- A model that passes is only as good as its invariants. Write the invariant first ("acked ⇒ durable"), then the spec.
- Model the crash as a step that can happen *anywhere*. That's where WAL bugs live (module 11).
- Liveness checks need fairness (weak/strong). Without them, TLC reports "a process can just stop", which is a trivial counterexample.
- Simulation is only as deterministic as your discipline: one call to `time.Now()` or one free goroutine breaks replay. Ban them in simulated packages (lint rule, or interfaces only).
- `testing/synctest` already gives you deterministic *time*. For many timeout/TTL/limiter bugs it's enough. Know where it stops.

## Questions you must answer
1. Model checking vs stress testing vs simulation: what does each find that the others miss?
2. How would you gain confidence that your log writer never acknowledges an undurable record?
3. What is deterministic simulation testing? Why is it hard in Go, and how do teams do it anyway?
4. Safety vs liveness: give an example of each for a bounded queue.
5. How do you keep a spec and its implementation from drifting apart?

## Reading
- TLA+ (Lamport): https://lamport.azurewebsites.net/tla/tla.html, the TLA+ Foundation: https://foundation.tlapl.us/
- Hillel Wayne, *Learn TLA+*: https://learntla.com/
- TLA+ tools: https://github.com/tlaplus/tlaplus · Examples: https://github.com/tlaplus/Examples
- Markus Kuppe, *BlockingQueue* step-by-step tutorial: https://github.com/lemmy/BlockingQueue (start here for `tla/boundedqueue/`)
- Will Wilson, *Testing Distributed Systems w/ Deterministic Simulation* (Strange Loop 2014): https://www.youtube.com/watch?v=4fFDFbi3toc
- *FoundationDB: A Distributed Unbundled Transactional Key Value Store* (SIGMOD 2021): https://www.foundationdb.org/files/fdb-paper.pdf

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**before `tla/boundedqueue/`**
- Read `[S]` [Newcombe et al.: How Amazon Web Services Uses Formal Methods (CACM 2015)](https://www.amazon.science/publications/how-amazon-web-services-uses-formal-methods): TLA+ found subtle bugs in DynamoDB/S3 designs; an engineering-time investment that paid off.

**after `tla/groupcommit/`**
- Read `[S]` [MongoDB: eXtreme Modelling in Practice (VLDB 2020)](https://arxiv.org/abs/2006.00915): trace-checking and model-based test-case generation against a real codebase; what worked and what was too costly.
- Code `[T]` [Azure/azure-cosmos-tla](https://github.com/Azure/azure-cosmos-tla): TLA+ specs of Cosmos DB's consistency levels, published as the contract.

**with the etcd trace-validation item**
- Code `[T]` [etcd-io/raft `tla/`](https://github.com/etcd-io/raft/tree/main/tla): the raft spec plus trace validation against the Go implementation.

**before `dst/`**
- Read `[S]` [Antithesis: Deterministic simulation testing](https://antithesis.com/resources/deterministic_simulation_testing/): how DST works and when it's worth it.
- Read `[S]` [TigerBeetle: Safety (the VOPR)](https://docs.tigerbeetle.com/concepts/safety/): simulating years of faults in hours. This is DST as a core design constraint.
- Read `[S]` [Polar Signals: (Mostly) Deterministic Simulation Testing in Go](https://www.polarsignals.com/blog/posts/2024/05/28/mostly-dst-in-go): DST for a Go database via WASM and a modified runtime, with honest limits.
