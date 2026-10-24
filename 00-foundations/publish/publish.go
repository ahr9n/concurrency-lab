// Package publish compares ways to hand a fully built value from one goroutine to others.
package publish

type Config struct {
	Addr    string
	Timeout int
	Tags    []string
}

// Mailbox delivers one *Config from a single writer to any number of readers.
// Publish is called exactly once; Await returns the published value, waiting if necessary.
type Mailbox interface {
	Publish(cfg *Config)
	Await() *Config
}
