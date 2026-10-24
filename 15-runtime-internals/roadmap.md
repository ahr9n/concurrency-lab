# 15: Go runtime internals (source-reading track)

> **Tiers:** `[E]` essential (interview + daily work) · `[S]` senior (production at scale) · `[D]` Go depth (runtime mastery) · `[T]` theory / rare. See [ROADMAP.md](../ROADMAP.md#tiers).

**Goal:** answer "how does Go actually do this?" at the level of the runtime source. About 1 h/week, alongside the other modules.

## Reading sessions (open the file and trace the path; write ≤ 1 page of notes per session in NOTES.md)
- [ ] `[S]` **S1: GMP model.** `runtime/runtime2.go` (`g`, `m`, `p`, `schedt`), `proc.go` `schedule()`, `findRunnable()`. Draw the diagram.
- [ ] `[D]` **S2: Run queues and work stealing.** `runqput`, `runqget`, `runqsteal`, the `runnext` slot, global run queue fairness (every 61 schedticks).
- [ ] `[D]` **S3: Preemption.** Cooperative (stack-check prologue) vs async (signal-based, 1.14+). `sysmon`. Why a tight loop no longer freezes GC.
- [ ] `[S]` **S4: Syscalls and netpoller.** `entersyscall`/`exitsyscall`, P handoff, `netpoll` (epoll/kqueue). Why blocking I/O doesn't block other goroutines.
- [ ] `[S]` **S5: Channels.** `chan.go` `chansend`, `chanrecv`, `closechan`; `select.go` `selectgo` (lock ordering, random poll order, `sudog`s).
- [ ] `[D]` **S6: Mutex and semaphores.** `internal/sync/mutex.go` (states, spinning, starvation), `runtime/sema.go` (the treap of waiters).
- [ ] `[D]` **S7: sync.Map / HashTrieMap, sync.Pool** (`poolqueue.go`, per-P, victim cache).
- [ ] `[D]` **S8: Timers.** `runtime/time.go` after the 1.23 rework.
- [ ] `[S]` **S9: GC interaction.** Write barriers, STW phases, GC assist (how allocation-heavy goroutines get taxed), Green Tea (1.26 default). Why the GC makes lock-free in Go easier.
- [ ] `[D]` **S10: Stacks.** Growth by copying, `morestack`. Why you can't take stable pointers to stack memory across growth (and why the GC handles it).
- [ ] `[D]` **S11: Coroutines.** `runtime/coro.go` behind `iter.Pull`.
- [ ] `[D]` **S12: Leak detection.** How the `goroutineleak` profile uses GC reachability.

## Experiments
- [ ] `[D]` `GODEBUG=schedtrace=1000,scheddetail=1` while running your worker pool. Read the P and M states.
- [ ] `[D]` Pin a goroutine in a tight loop (no function calls) with `GOMAXPROCS=1` and observe async preemption. Disable it with `GODEBUG=asyncpreemptoff=1`.
- [ ] `[S]` `go tool trace` your log writer and find the syscall (fsync) blocks and the P handoff.
