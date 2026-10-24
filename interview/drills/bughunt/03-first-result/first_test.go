//go:build bughunt

package first

import (
	"context"
	"testing"
	"testing/synctest"
	"time"
)

func replica(name string, latency time.Duration) func(context.Context) string {
	return func(ctx context.Context) string {
		select {
		case <-time.After(latency):
		case <-ctx.Done():
		}
		return name
	}
}

// synctest fails the test if any goroutine is still blocked when the test returns.
func TestFirstLeavesNoGoroutinesBehind(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		replicas := []func(context.Context) string{
			replica("eu", 30*time.Millisecond),
			replica("us", 10*time.Millisecond),
			replica("ap", 50*time.Millisecond),
		}

		got := First(t.Context(), replicas)
		if got != "us" {
			t.Fatalf("First = %q, want %q", got, "us")
		}

		time.Sleep(time.Second)
	})
}
