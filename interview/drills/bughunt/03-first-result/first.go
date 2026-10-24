//go:build bughunt

package first

import "context"

// First queries every replica and returns the fastest answer.
func First(ctx context.Context, replicas []func(context.Context) string) string {
	results := make(chan string)
	for _, query := range replicas {
		go func() {
			results <- query(ctx)
		}()
	}
	return <-results
}
