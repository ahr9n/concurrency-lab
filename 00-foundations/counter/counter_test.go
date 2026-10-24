package counter

import (
	"sync"
	"testing"
)

// Fill in each constructor as you implement it; nil entries are skipped.
var implementations = []struct {
	name string
	new  func() Counter
}{
	{"mutex", nil},
	{"atomic", nil},
	{"cas", nil},
	{"sharded", nil},
	{"sharded-padded", nil},
}

func eachImpl(t *testing.T, test func(t *testing.T, newCounter func() Counter)) {
	t.Helper()
	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			if impl.new == nil {
				t.Skip("not implemented yet")
			}
			test(t, impl.new)
		})
	}
}

func TestCountsEveryIncrement(t *testing.T) {
	eachImpl(t, func(t *testing.T, newCounter func() Counter) {
		const goroutines, perGoroutine = 32, 2_000
		c := newCounter()

		// A start gate maximises overlap; without it early goroutines finish before late ones start.
		start := make(chan struct{})
		var wg sync.WaitGroup
		for range goroutines {
			wg.Go(func() {
				<-start
				for range perGoroutine {
					c.Inc()
				}
			})
		}
		close(start)
		wg.Wait()

		got := c.Load()
		if got != goroutines*perGoroutine {
			t.Fatalf("Load = %d, want %d", got, goroutines*perGoroutine)
		}
	})
}

func TestLoadNeverGoesBackwards(t *testing.T) {
	eachImpl(t, func(t *testing.T, newCounter func() Counter) {
		c := newCounter()
		done := make(chan struct{})

		var wg sync.WaitGroup
		for range 8 {
			wg.Go(func() {
				for {
					select {
					case <-done:
						return
					default:
						c.Inc()
					}
				}
			})
		}

		var last int64
		for range 10_000 {
			now := c.Load()
			if now < last {
				close(done)
				wg.Wait()
				t.Fatalf("Load went backwards: %d after %d", now, last)
			}
			last = now
		}
		close(done)
		wg.Wait()
	})
}

// Run with: make bench M=00
func BenchmarkInc(b *testing.B) {
	for _, impl := range implementations {
		b.Run(impl.name, func(b *testing.B) {
			if impl.new == nil {
				b.Skip("not implemented yet")
			}
			c := impl.new()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					c.Inc()
				}
			})
		})
	}
}

// Sharding makes Inc cheap and Load expensive; this shows the cost when reads are not rare.
func BenchmarkIncWithReads(b *testing.B) {
	for _, impl := range implementations {
		b.Run(impl.name, func(b *testing.B) {
			if impl.new == nil {
				b.Skip("not implemented yet")
			}
			c := impl.new()
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					i++
					if i%100 == 0 {
						_ = c.Load()
						continue
					}
					c.Inc()
				}
			})
		})
	}
}
