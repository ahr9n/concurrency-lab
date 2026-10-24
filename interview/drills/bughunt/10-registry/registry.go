//go:build bughunt

package registry

import (
	"sync"
	"sync/atomic"
	"time"
)

type Client struct {
	Name string
}

type Registry struct {
	clients sync.Map
	dials   atomic.Int64
}

// Get returns the shared client for name. Dialing is expensive and must happen once per name.
func (r *Registry) Get(name string) *Client {
	if c, ok := r.clients.Load(name); ok {
		return c.(*Client)
	}

	c := r.dial(name)
	r.clients.Store(name, c)
	return c
}

func (r *Registry) Dials() int64 {
	return r.dials.Load()
}

func (r *Registry) dial(name string) *Client {
	r.dials.Add(1)
	time.Sleep(5 * time.Millisecond) // stands in for a network handshake
	return &Client{Name: name}
}
