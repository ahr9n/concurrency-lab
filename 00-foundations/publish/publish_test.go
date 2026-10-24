package publish

import (
	"os"
	"sync"
	"testing"
	"time"
)

// Fill in each constructor as you implement it; nil entries are skipped.
// "plain-bool" is meant to be broken: it is skipped unless RUN_BROKEN=1, so you can watch
// the race detector catch it without turning CI red:
//
//	RUN_BROKEN=1 go test -race -run 'TestReadersSeeCompleteConfig/plain-bool' ./00-foundations/publish
var implementations = []struct {
	name   string
	broken bool
	new    func() Mailbox
}{
	{"plain-bool", true, nil},
	{"atomic-bool", false, nil},
	{"atomic-pointer", false, nil},
	{"channel-close", false, nil},
}

func TestReadersSeeCompleteConfig(t *testing.T) {
	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			if impl.new == nil {
				t.Skip("not implemented yet")
			}
			if impl.broken && os.Getenv("RUN_BROKEN") != "1" {
				t.Skip("intentionally broken; set RUN_BROKEN=1 to run it")
			}
			for range 1_000 {
				publishOnce(t, impl.new())
			}
		})
	}
}

func publishOnce(t *testing.T, m Mailbox) {
	t.Helper()
	const readers = 4
	results := make(chan *Config, readers)

	var wg sync.WaitGroup
	for range readers {
		wg.Go(func() {
			results <- m.Await()
		})
	}
	go m.Publish(&Config{Addr: "db:5432", Timeout: 30, Tags: []string{"a", "b", "c"}})

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("readers still waiting 5s after Publish")
	}
	close(results)

	for cfg := range results {
		complete := cfg != nil && cfg.Addr == "db:5432" && cfg.Timeout == 30 && len(cfg.Tags) == 3
		if !complete {
			t.Fatalf("reader saw an incomplete config: %+v", cfg)
		}
	}
}
