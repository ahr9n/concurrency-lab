# Global reference

Cross-cutting sources only. Anything tied to one topic (papers, engineering reports, production code, tools) lives in that module's `reference.md`, under *Reading* and *Real-world reading*.

## Primary sources (read these, not summaries of them)
- The Go Memory Model: https://go.dev/ref/mem
- Russ Cox, *Memory Models* series (hardware, programming-language, Go updates): https://research.swtch.com/mm
- Go spec, sections on goroutines, channels, select: https://go.dev/ref/spec
- Package docs: `sync`, `sync/atomic`, `context`, `testing/synctest`, `runtime/pprof`, `runtime/trace`
  - https://pkg.go.dev/sync and the others
- `golang.org/x/sync` (errgroup, semaphore, singleflight): https://pkg.go.dev/golang.org/x/sync
- `golang.org/x/time/rate`: https://pkg.go.dev/golang.org/x/time/rate
- Release notes (concurrency-relevant):
  - 1.14: async preemption
  - 1.19: atomic types, memory model revision
  - 1.22: per-iteration loop vars
  - 1.23: timers rework
  - 1.24: sync.Map on HashTrieMap, swiss maps, `weak`, `runtime.AddCleanup`, `b.Loop`
  - 1.25: `synctest` GA, `WaitGroup.Go`, container-aware GOMAXPROCS, flight recorder
  - 1.26: Green Tea GC default, goroutine-leak profile (experimental)
  - 1.27: goroutine-leak profile GA, `synctest.Sleep`
  - All notes: https://go.dev/doc/devel/release
- Go runtime source (your local copy): `$(go env GOROOT)/src/runtime`, `src/internal/sync`, `src/sync`

## Books
- Katherine Cox-Buday, *Concurrency in Go* (O'Reilly). The Go-specific pattern catalogue.
- Teiva Harsanyi, *100 Go Mistakes and How to Avoid Them*, chapters 8–9 (concurrency).
- Herlihy & Shavit, *The Art of Multiprocessor Programming*. Theory: linearizability, locks, lock-free structures.
- Allen Downey, *The Little Book of Semaphores* (free): https://greenteapress.com/wp/semaphores/
- Paul McKenney, *Is Parallel Programming Hard, And, If So, What Can You Do About It?* (free): https://mirrors.edge.kernel.org/pub/linux/kernel/people/paulmck/perfbook/perfbook.html
- Martin Kleppmann, *Designing Data-Intensive Applications*. WAL, replication, consistency, distributed locks.
- Alex Petrov, *Database Internals*. Storage engines, WAL, LSM trees.

## Talks
- Rob Pike, *Concurrency is not Parallelism* (2012)
- Dmitry Vyukov, *Go scheduler: Implementing language with lightweight concurrency* (2019)
- Michael Knyszek / Michael Pratt, runtime talks (GopherCon 2023–2025)
