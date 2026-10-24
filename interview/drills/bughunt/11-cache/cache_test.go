//go:build bughunt

package cache

import (
	"errors"
	"testing"
	"time"
)

func TestMissThenHit(t *testing.T) {
	c := New()
	c.Set("a", "1")

	done := make(chan struct{})
	go func() {
		defer close(done)

		_, err := c.Get("missing")
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("Get(missing) err = %v, want ErrNotFound", err)
		}
		v, err := c.Get("a")
		if err != nil || v != "1" {
			t.Errorf("Get(a) = %q, %v; want 1, nil", v, err)
		}
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Get did not return within 2s")
	}
}
