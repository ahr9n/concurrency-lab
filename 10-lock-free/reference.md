# 10: Lock-free reference

## Concepts
- **Progress guarantees:**
  - blocking (locks)
  - obstruction-free (progress in isolation)
  - lock-free (*some* thread always progresses)
  - wait-free (*every* thread progresses in bounded steps)
- **CAS loop:** read → compute → CAS, retry on failure. Correctness rests on the CAS being the *linearization point*.
- **Linearizability:** each operation appears to take effect atomically at some point between its call and return. It's the correctness criterion for concurrent objects (Herlihy & Wing).
- **ABA:** a thread reads A, others change it to B and back to A, and the CAS succeeds wrongly. Classic with manual memory reuse. **In Go, the GC means a node can't be freed and reused while you hold a pointer to it, so pointer-based Treiber/MS structures are ABA-safe.** ABA can still bite with indices, pooled nodes (`sync.Pool`), or value-encoded state. Fixes: version stamps, immutable boxes.
- **Memory reclamation:** in C/C++/Rust you need hazard pointers, epoch-based reclamation, or RCU. In Go the GC does it, which is a big simplification and worth saying in interviews.
- **False sharing:** the head/tail indices of a ring buffer must sit on separate cache lines.
- **Copy-on-write with `atomic.Pointer`:** the Go-idiomatic "RCU-lite". Readers are wait-free. Writers serialize on a mutex and allocate a copy each time.
- **Contention collapse:** under heavy contention a CAS loop spins and bounces the cache line. Locks park waiters. That's why lock-free isn't automatically faster. Backoff helps.

## Hidden knowledge / gotchas
- Go has no `CompareAndSwap` on a two-word struct (no DWCAS). Pack the values into 64 bits, or swap a pointer to an immutable struct.
- Go atomics are sequentially consistent. There are no relaxed/acquire variants, so you can't hand-tune barriers the way you would in C++. That limits the micro-optimisations available, and it keeps you safe.
- The Go runtime itself uses lock-free structures: per-P run queues (`runqput`/`runqget`), `sync.Pool`'s `poolDequeue`, and the `sudog` lists. Read `runtime/proc.go` `runqput`.
- `sync.Map` (1.24+) is a lock-free-read hash-trie. A good real-world example.
- Testing lock-free code: `-race` catches non-atomic access but not logic bugs. Use linearizability checkers and very long stress runs.

## Questions you must answer
1. Implement a lock-free stack. What is ABA, and does it apply in Go?
2. When is lock-free *slower* than a mutex?
3. Define linearizability. How would you test it?
4. Lock-free vs wait-free: give an example of each.
5. Why does an SPSC queue need no CAS?
6. How does ConcurrentHashMap (Java) / sync.Map (Go) achieve lock-free reads?

## Reading
- Herlihy & Shavit, *The Art of Multiprocessor Programming*, ch. 3 (linearizability), 10 (queues), 11 (stacks)
- Michael & Scott 1996; Treiber 1986 (IBM tech report)
- Dmitry Vyukov, 1024cores.net (MPSC/SPMC queue designs)
- LMAX Disruptor paper / tech notes
- https://github.com/anishathalye/porcupine
- `$(go env GOROOT)/src/runtime/proc.go` (runq*)

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**after `cow/`**
- Code `[S]` [grpc-go `balancer/roundrobin/roundrobin.go`](https://github.com/grpc/grpc-go/blob/master/balancer/roundrobin/roundrobin.go): immutable picker snapshots swapped atomically: copy-on-write in production.

**after `spsc/`**
- Code `[D]` [go `src/sync/poolqueue.go`](https://github.com/golang/go/blob/master/src/sync/poolqueue.go): the lock-free SPMC deque behind `sync.Pool`.

**after `seqlock/`**
- Code `[T]` [gvisor `pkg/sync/seqcount.go`](https://github.com/google/gvisor/blob/master/pkg/sync/seqcount.go): a seqlock in Go; see how they work within the memory model.
