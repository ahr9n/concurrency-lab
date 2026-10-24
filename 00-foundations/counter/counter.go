// Package counter compares ways to count from many goroutines at once.
package counter

// Counter is implemented five ways in this package (see ../roadmap.md).
// Register each implementation in the table in counter_test.go.
type Counter interface {
	Inc()
	Load() int64
}
