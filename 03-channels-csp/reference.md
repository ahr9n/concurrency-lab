# 03: Channels reference

## Channel state table (memorise)
| Op | nil | open | closed |
|---|---|---|---|
| send | blocks forever | blocks until receiver or buffer space | **panic** |
| receive | blocks forever | blocks until value | zero value, `ok=false` (after buffer drained) |
| close | **panic** | ok | **panic** |
| len/cap | 0 | current/capacity | buffered remainder / capacity |

## Concepts
- **Ownership rule:** the goroutine that *writes* to a channel owns it, closes it, and is the only closer. Receivers never close.
- **Close is a broadcast:** every receiver wakes. It's the cheapest one-to-many signal in Go (`done` channels, `ctx.Done()`).
- **Unbuffered = synchronization point:** the receive happens before the send completes. **Buffered(C) = semaphore:** the k-th receive happens before the (k+C)-th send completes.
- **select:**
  - Picks uniformly at random among ready cases, so there's no priority.
  - A `nil` channel case never fires, so you can disable a case by setting its channel to nil.
  - `default` makes the select non-blocking.
- **Internals** (detail in `15-runtime-internals`): `hchan` = lock + ring buffer (`buf`, `sendx`, `recvx`) + `sendq`/`recvq` wait queues of `sudog`s.
  - A send to a *waiting receiver* copies directly into the receiver's stack (no buffer hop).
  - `select` locks all involved channels **in address order** to avoid deadlock.
- **Timers (1.23 rework):**
  - Timer and Ticker channels are effectively unbuffered.
  - After `Stop`/`Reset`, no stale value is received.
  - Unreferenced, unstopped timers are GC-eligible, so `time.After` in a loop no longer leaks until it fires.
  - 1.27 removed the `asynctimerchan` GODEBUG escape hatch.
  - Still prefer one `time.NewTimer` + `Reset` in hot loops to avoid allocation.
- **Range-over-func iterators (1.23):** `iter.Pull` is implemented with runtime coroutines (a direct switch, not scheduler handoff). A goroutine+channel iterator is much slower and can leak if the consumer stops early.

## Patterns vocabulary
generator · pipeline · fan-out · fan-in · ordered fan-out · or-done · tee · bridge · heartbeat ·
replicated requests · batching · debounce · throttle · pub/sub · error propagation (`errgroup`) · bounded parallelism (semaphore).

## Hidden knowledge / gotchas
- **Most Go concurrency bugs in the wild are message-passing bugs, not shared-memory bugs** (Tu et al., ASPLOS 2019): blocking forever on a send nobody receives.
- A sender that might outlive its receiver needs `select { case ch <- v: case <-ctx.Done(): return }`. A bare send leaks.
- `for v := range ch` only ends when the channel is closed. Forgetting to close it is a leak.
- The `tee` pattern blocks on the slowest consumer. You need per-consumer buffering or a drop policy.
- `len(ch)` is a racy snapshot. Never base correctness on it.
- Channels of large structs copy the whole value. Send pointers, or keep values small.
- A buffered channel used as a queue is *not* a bounded queue with cancellation, timeouts, and metrics. Know the limits (see `05-bounded-queue`).
- Closing a channel to signal "done" while other goroutines may still send means you need a separate `done` channel. Never close from the receiver side.

## Questions you must answer
1. What happens when you send to or close a closed channel? Receive from a nil one?
2. Who should close a channel when there are multiple senders?
3. How would you implement a priority select?
4. How does `select` avoid deadlock when it locks multiple channels?
5. Unbuffered vs buffered: what are the semantic differences, not just the performance ones?
6. Design a batcher that flushes on size or time. What's the bug with `time.After` in the loop?
7. Fan-out while preserving output order: how?
8. Channels vs mutex: when do you choose which? (Rule of thumb: channels for transferring ownership or signalling, mutexes for protecting state.)

## Reading
- Go blog: *Go Concurrency Patterns: Pipelines and cancellation*; *Go Concurrency Patterns: Timing out, moving on*
- Talks: Pike *Go Concurrency Patterns*, Ajmani *Advanced Go Concurrency Patterns*, Kavya Joshi *Understanding Channels*
- Source: `$(go env GOROOT)/src/runtime/chan.go`, `src/runtime/select.go`
- https://go.dev/wiki/Go123Timer
- *Concurrency in Go*, ch. 4

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**before you start**
- Read `[S]` [Real-World Concurrency Bugs in Go (ASPLOS 2019)](https://songlh.github.io/paper/go-study.pdf): 171 bugs from Docker/k8s/etcd: message passing causes as many bugs as shared memory.

**after ordered fan-out**
- Code `[S]` [conc `stream/stream.go`](https://github.com/sourcegraph/conc/blob/main/stream/stream.go): concurrent work, callbacks in submit order: ordered fan-out with head-of-line blocking.

**after `pubsub/` (slow subscribers)**
- Code `[S]` [nats-server `server/client.go`](https://github.com/nats-io/nats-server/blob/main/server/client.go) (search `flushOutbound`): per-client outbound buffers; slow consumers are dropped at the write deadline.
