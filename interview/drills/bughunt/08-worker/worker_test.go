//go:build bughunt

package worker

import (
	"context"
	"testing"
	"testing/synctest"
)

func TestRunStopsOnCancel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		jobs := make(chan int)
		done := make(chan struct{})

		handled := 0
		go func() {
			Run(ctx, jobs, func(int) { handled++ })
			close(done)
		}()

		jobs <- 1
		jobs <- 2
		cancel()
		synctest.Wait()

		select {
		case <-done:
		default:
			t.Fatal("Run did not return after ctx was cancelled")
		}
		if handled != 2 {
			t.Fatalf("handled = %d, want 2", handled)
		}
	})
}
