//go:build bughunt

package worker

import "context"

// Run processes jobs until ctx is cancelled.
func Run(ctx context.Context, jobs <-chan int, handle func(int)) {
	for job := range jobs {
		handle(job)
	}
}
