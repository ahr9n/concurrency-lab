//go:build bughunt

package transfer

import (
	"sync"
	"testing"
	"time"
)

func TestOppositeTransfers(t *testing.T) {
	a := NewAccount(1, 1_000_000)
	b := NewAccount(2, 1_000_000)

	done := make(chan struct{})
	go func() {
		var wg sync.WaitGroup
		wg.Go(func() {
			for range 100_000 {
				_ = Transfer(a, b, 1)
			}
		})
		wg.Go(func() {
			for range 100_000 {
				_ = Transfer(b, a, 1)
			}
		})
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("transfers did not finish within 10s")
	}

	total := a.Balance() + b.Balance()
	if total != 2_000_000 {
		t.Fatalf("total = %d, want 2000000", total)
	}
}
