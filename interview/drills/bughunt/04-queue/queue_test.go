//go:build bughunt

package queue

import (
	"sort"
	"testing"
	"testing/synctest"
)

func TestTwoConsumersTwoItems(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		q := New()
		got := make(chan int, 2)
		for range 2 {
			go func() {
				got <- q.Take()
			}()
		}

		synctest.Wait()
		q.Put(1)
		synctest.Wait()
		q.Put(2)

		values := []int{<-got, <-got}
		sort.Ints(values)
		if values[0] != 1 || values[1] != 2 {
			t.Fatalf("got %v, want [1 2]", values)
		}
	})
}
