//go:build bughunt

package square

import "sync"

func SquareAll(nums []int) []int {
	var (
		mu  sync.Mutex
		wg  sync.WaitGroup
		out = make([]int, 0, len(nums))
	)
	for _, n := range nums {
		go func() {
			wg.Add(1)
			defer wg.Done()

			mu.Lock()
			out = append(out, n*n)
			mu.Unlock()
		}()
	}
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	return out
}
