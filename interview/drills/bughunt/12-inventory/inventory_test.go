//go:build bughunt

package inventory

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestReadsAndWritesTogether(t *testing.T) {
	inv := New()
	names := make([]string, 100)
	for i := range names {
		names[i] = fmt.Sprintf("sku-%d", i)
		inv.Set(names[i], 1)
	}

	done := make(chan struct{})
	go func() {
		var wg sync.WaitGroup
		for range 4 {
			wg.Go(func() {
				for range 2_000 {
					inv.Total(names...)
				}
			})
		}
		wg.Go(func() {
			for i := range 20_000 {
				inv.Set(names[i%len(names)], i)
			}
		})
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("readers and writers did not finish within 10s")
	}
}
