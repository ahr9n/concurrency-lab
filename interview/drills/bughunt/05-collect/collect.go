//go:build bughunt

package collect

// Collect reads every source concurrently and returns all values.
func Collect(sources [][]int) []int {
	out := make(chan int)
	for _, src := range sources {
		go func() {
			for _, v := range src {
				out <- v
			}
			close(out)
		}()
	}

	var all []int
	for v := range out {
		all = append(all, v)
	}
	return all
}
