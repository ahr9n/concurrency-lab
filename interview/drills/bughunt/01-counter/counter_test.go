//go:build bughunt

package counter

import (
	"sync"
	"testing"
)

func TestConcurrentInc(t *testing.T) {
	c := New()
	var wg sync.WaitGroup
	for range 50 {
		wg.Go(func() {
			for range 1000 {
				c.Inc("page")
			}
		})
	}
	wg.Wait()

	got := c.Get("page")
	if got != 50_000 {
		t.Fatalf("Get = %d, want 50000", got)
	}
}
