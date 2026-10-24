# Bug hunt

Twelve small Go packages, each with one realistic concurrency bug and a test that exposes it.
**Your job:** find the bug, explain it in one sentence, fix it, and get the test green under the race detector.

They're behind the `bughunt` build tag, so normal `go test ./...`, lint, and CI skip them.

```sh
make bug N=04        # run one drill (race detector on)
make bugs            # run all of them (expect failures until you fix them)
make bug-reset       # throw away your fixes and start again
```

**Rules:**
1. Don't open [ANSWERS.md](ANSWERS.md) until the drill passes, or until you've spent 20 minutes on it.
2. Before fixing, try `go vet -tags bughunt ./interview/drills/bughunt/...`: it catches two of the twelve. Write down which tool *should* have caught each bug (`-race`, `go vet`, `synctest`, a goroutine dump, the leak profile). That's half the lesson.
3. Time yourself and log it in [../LOG.md](../LOG.md).

| # | Drill | Symptom you'll see |
|---|---|---|
| 01 | `counter` | race report / fatal error under load |
| 02 | `transfer` | test hangs, then times out |
| 03 | `first-result` | synctest complains about goroutines left behind |
| 04 | `queue` | panic: index out of range |
| 05 | `collect` | panic or missing values |
| 06 | `square-all` | results are missing |
| 07 | `stats` | the count is always 0 |
| 08 | `worker` | `Run` never returns after cancel |
| 09 | `config` | race report |
| 10 | `registry` | the expensive dial runs many times |
| 11 | `cache` | the second call hangs |
| 12 | `inventory` | hangs under mixed reads and writes |
