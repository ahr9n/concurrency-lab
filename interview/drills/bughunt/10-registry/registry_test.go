//go:build bughunt

package registry

import (
	"sync"
	"testing"
)

func TestOneDialPerName(t *testing.T) {
	var r Registry
	start := make(chan struct{})
	clients := make([]*Client, 50)

	var wg sync.WaitGroup
	for i := range clients {
		wg.Go(func() {
			<-start
			clients[i] = r.Get("payments-db")
		})
	}
	close(start)
	wg.Wait()

	if r.Dials() != 1 {
		t.Fatalf("Dials = %d, want 1", r.Dials())
	}
	for _, c := range clients {
		if c != clients[0] {
			t.Fatal("callers received different clients for the same name")
		}
	}
}
