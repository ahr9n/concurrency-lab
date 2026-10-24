//go:build bughunt

package stats

import "sync"

type Stats struct {
	mu       sync.Mutex
	requests int
}

func (s Stats) Record() {
	s.mu.Lock()
	s.requests++
	s.mu.Unlock()
}

func (s *Stats) Requests() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.requests
}
