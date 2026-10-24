# 15: Runtime internals reference

## Must-know model
- **G** = goroutine (stack, status, `sched` context). **M** = OS thread. **P** = processor (scheduling context: local run queue up to 256 Gs + a `runnext` slot, mcache). `GOMAXPROCS` = number of Ps.
- **Scheduling:** an M must hold a P to run Go code. `findRunnable` checks, in order:
  1. `runnext`
  2. the local queue
  3. the global queue (checked every 61 ticks for fairness)
  4. netpoll
  5. stealing half of another P's queue
- **Blocking syscall:** the M leaves with its G, and the P is handed to another M (sysmon retakes Ps stuck in a syscall longer than about 20 µs). **Network I/O** goes through the netpoller instead, so the G parks and the M stays free.
- **Preemption:** synchronous at function prologues (stack bound check) plus **asynchronous via signals (SIGURG) since 1.14**, so tight loops can be preempted.
- **Channels:** `hchan{lock, buf, sendx, recvx, qcount, dataqsiz, recvq, sendq, closed}`. A blocked goroutine waits as a `sudog` on `recvq`/`sendq`. A direct handoff copies to the receiver's stack. `select` shuffles the poll order (fairness) and locks channels in address order (no deadlock).
- **Mutex:** a state word (locked, woken, starving, waiter count) + a runtime semaphore. Spins only on multicore with other running Ps and an empty local runq.
- **Semaphores** (`sema.go`): waiters live in a balanced tree keyed by address, plus FIFO lists per address.
- **GC:** concurrent tri-colour mark-sweep with a hybrid write barrier and short STW phases. Mutator **assists** make allocating goroutines pay for marking. Green Tea (1.26) improves the locality of marking small objects.
- **Stacks:** start small and grow by copying (so pointers into stacks get adjusted). This is one reason Go forbids pointer arithmetic.

## Questions you must answer
1. Explain GMP. Why does P exist? (Scalability: per-P caches and queues; decoupling from blocked Ms.)
2. What happens when a goroutine makes a blocking syscall? A network read?
3. How does the scheduler achieve fairness? What is work stealing?
4. How is a goroutine preempted? What changed in 1.14?
5. Walk through `ch <- v` when (a) a receiver is waiting, (b) the buffer has space, (c) the buffer is full.
6. How does `select` choose? How does it avoid deadlock?
7. How many OS threads does a Go program use? What bounds it? (`debug.SetMaxThreads`, default 10k)
8. Goroutines vs OS threads vs Java virtual threads: compare stack, scheduling, and pinning.

## Reading
- Kavya Joshi, *The Scheduler Saga* (GopherCon 2018) and *Understanding Channels* (2017)
- Austin Clements, *Proposal: Non-cooperative goroutine preemption* (#24543)
- Source: `$(go env GOROOT)/src/runtime/{proc.go,runtime2.go,chan.go,select.go,sema.go,time.go,coro.go,mgc.go}`

## Real-world reading
Do each block when you reach that step. For each one, note the trade-off the authors accepted.

**with S1–S2 (scheduler)**
- Read `[D]` [Vyukov: Scalable Go Scheduler Design Doc](https://go.dev/s/go11sched): replacing the global scheduler lock with per-P work-stealing run queues.

**with S9 (GC)**
- Read `[D]` [Go team: Getting to Go (GC history)](https://go.dev/blog/ismmkeynote): non-moving, non-generational, concurrent: latency chosen over throughput.
- Read `[S]` [Go team: A Guide to the Go GC](https://go.dev/doc/gc-guide): GOGC trades CPU for memory; a `GOMEMLIMIT` set too low thrashes. Required reading.

**after S9: GC at scale**
- Read `[S]` [Uber: How We Saved 70K Cores (GOGC)](https://www.uber.com/blog/how-we-saved-70k-cores-across-30-mission-critical-services/): dynamic GOGC against the memory limit saved huge CPU, with some OOM risk.
- Read `[S]` [Twitch: Go memory ballast](https://blog.twitch.tv/en/2019/04/10/go-memory-ballast-how-i-learnt-to-stop-worrying-and-love-the-heap/): a 10 GiB ballast cut p99 by ~45%; a hack that `GOMEMLIMIT` replaced. Know why.
- Read `[S]` [Datadog: tracking down a Go 1.24 memory regression](https://www.datadoghq.com/blog/engineering/go-memory-regression/): RSS +20% while heap profiles looked clean; the cause was a runtime allocator change.
