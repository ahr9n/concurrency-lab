//go:build bughunt

package collect

import (
	"slices"
	"testing"
)

func TestCollectAll(t *testing.T) {
	for range 100 {
		got := Collect([][]int{{1}, {2, 3, 4}, {5, 6}})
		slices.Sort(got)

		want := []int{1, 2, 3, 4, 5, 6}
		if !slices.Equal(got, want) {
			t.Fatalf("Collect = %v, want %v", got, want)
		}
	}
}
