//go:build bughunt

package stats

import (
	"sync"
	"testing"
)

func TestRecord(t *testing.T) {
	var s Stats
	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			for range 100 {
				s.Record()
			}
		})
	}
	wg.Wait()

	got := s.Requests()
	if got != 1000 {
		t.Fatalf("Requests = %d, want 1000", got)
	}
}
