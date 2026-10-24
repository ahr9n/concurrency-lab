# Bug hunt: answers (spoilers)

Each entry gives the bug, why it happens, the fix, the tool that should catch it, and the production lesson.

## 01 counter: unsynchronized map
- **Bug:** concurrent `c.hits[key]++` on a plain map.
- **Why:** map writes aren't safe for concurrent use. The runtime may `fatal error: concurrent map writes` (and you can't recover from that), and `++` is a read-modify-write besides.
- **Fix:** guard with a `sync.Mutex`, or use a sharded map / atomic counters per key.
- **Tool:** `-race`.
- **Lesson:** "it's just a counter" is the most common race in real code (see Uber's *Data Race Patterns in Go*, module 00).

## 02 transfer: lock-order deadlock
- **Bug:** `Transfer(a,b)` locks a→b while `Transfer(b,a)` locks b→a.
- **Why:** circular wait (a Coffman condition). Each goroutine holds one lock and waits for the other.
- **Fix:** always lock in a global order, e.g. by account ID.
- **Tool:** a goroutine dump (two goroutines in `sync.Mutex.Lock`, each holding the other's lock). The runtime's deadlock detector stays silent because other goroutines (the test timer) are alive.
- **Lesson:** this is the classic interview "bank transfer" (module 13, #16) and a real production outage pattern.

## 03 first-result: goroutine leak
- **Bug:** unbuffered `results` channel. Only the first sender is received, so the slower replicas block forever on send.
- **Fix:** buffer the channel with `len(replicas)`, *or* cancel a derived ctx after the first result and `select` on send vs `ctx.Done()`. Best is both: cancel to save work, buffer to guarantee no leak.
- **Tool:** `synctest` (blocked goroutines at bubble exit), goleak, the `goroutineleak` profile.
- **Lesson:** most real Go concurrency bugs are blocked channel operations (Tu et al., ASPLOS 2019).

## 04 queue: `if` instead of `for` around `Cond.Wait`
- **Bug:** `Broadcast` wakes both consumers. The first takes the item, and the second continues past the `if` into an empty slice.
- **Fix:** `for len(q.items) == 0 { q.cond.Wait() }`.
- **Tool:** `synctest` makes the interleaving deterministic, and the panic does the rest.
- **Lesson:** always re-check the predicate after waking (stolen and spurious wakeups). This is the #1 bounded-queue follow-up (module 05).

## 05 collect: producer closes a shared channel
- **Bug:** every producer calls `close(out)`. The first to finish closes it while the others are still sending, so you get `send on closed channel`, `close of closed channel`, or missing values.
- **Fix:** a `WaitGroup` across producers, plus one separate goroutine that waits and then closes.
- **Tool:** `-race` (it reports the close/send race), plus the panic.
- **Lesson:** the sender closes, and with multiple senders a *coordinator* closes (module 03 ownership rule).

## 06 square-all: `wg.Add` inside the goroutine
- **Bug:** `Add(1)` runs inside the new goroutine, so `Wait` can observe zero before any `Add` happened.
- **Fix:** call `wg.Add(1)` before the `go` statement, or use `wg.Go(func(){...})` (Go 1.25+).
- **Tool:** **`go vet`** (the `waitgroup` analyzer, Go 1.25+) reports "WaitGroup.Add called from inside new goroutine".
- **Lesson:** `Add` must *happen before* `Wait`. This is a happens-before question in disguise (module 00 / 02).

## 07 stats: mutex copied by a value receiver
- **Bug:** `func (s Stats) Record()` copies the struct, *including the mutex and the counter*, so it increments a copy.
- **Fix:** a pointer receiver, `func (s *Stats) Record()`.
- **Tool:** **`go vet` (copylocks)** flags it immediately. This is why CI runs vet.
- **Lesson:** never copy a struct that contains a `sync` type after first use. Keep receivers consistent.

## 08 worker: cancellation ignored
- **Bug:** `for job := range jobs` only exits when `jobs` is closed. It never looks at `ctx`.
- **Fix:** `select` on `ctx.Done()` and `jobs` (handling `ok == false`).
- **Tool:** `synctest` (Run is still blocked after cancel), goleak.
- **Lesson:** context cancellation is cooperative. Every blocking point needs a ctx case (module 04).

## 09 config: unsynchronized publication
- **Bug:** `ready = true` and `current = …` are plain writes read by another goroutine. There's no happens-before edge, so the reader can see `ready` without seeing `current`, or spin forever.
- **Fix:** publish with `atomic.Pointer[Config]` (or a channel close, or `sync.Once`).
- **Tool:** `-race`.
- **Lesson:** a flag variable is not synchronization. This is the Go version of "double-checked locking needs volatile" (module 00 `publish/`).

## 10 registry: check-then-act on `sync.Map`
- **Bug:** `Load` miss → `dial` → `Store`. Fifty callers all miss at once and all dial. Each `sync.Map` call is atomic, but the sequence isn't.
- **Fix:** `LoadOrStore(name, sync.OnceValue(dial))` and call the stored func (one dial, and every caller shares the result). Or use `singleflight` plus a cache. Plain `LoadOrStore(name, dial())` is still wrong, because it dials before storing.
- **Tool:** a stress test with a start gate. `-race` won't catch it, because there's no data race, only a logic race.
- **Lesson:** atomic operations don't compose. This is a cache stampede in miniature (module 08 `loader/`).

## 11 cache: early return without unlock
- **Bug:** the not-found path returns while still holding `c.mu`, so the next call blocks forever.
- **Fix:** `defer c.mu.Unlock()` right after `Lock()`.
- **Tool:** a goroutine dump shows the goroutine stuck in `sync.Mutex.Lock` with no holder running.
- **Lesson:** lock and defer the unlock on adjacent lines. Hand-placed unlocks rot as code paths are added.

## 12 inventory: recursive `RLock`
- **Bug:** `Total` holds `RLock` and calls `Count`, which takes `RLock` again. If a writer calls `Lock` in between, the second `RLock` blocks behind the writer, and the writer waits for the first `RLock`. That's a deadlock.
- **Fix:** have the inner read use an unexported helper that assumes the lock is held (`inv.items[name]`).
- **Tool:** a goroutine dump (readers in `RLock`, a writer in `Lock`). It only shows up under write load, which makes it a classic flaky-test bug.
- **Lesson:** `sync.RWMutex` is writer-preferring, so recursive read locking is forbidden (module 02). Use the "locked helper" convention: `fooLocked()`.
