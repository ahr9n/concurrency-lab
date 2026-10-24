//go:build bughunt

package square

import "testing"

func TestSquareAllReturnsEverything(t *testing.T) {
	nums := make([]int, 200)
	for i := range nums {
		nums[i] = i
	}

	for range 50 {
		got := SquareAll(nums)
		if len(got) != len(nums) {
			t.Fatalf("len(SquareAll) = %d, want %d", len(got), len(nums))
		}
	}
}
