//go:build bughunt

package inventory

import "sync"

type Inventory struct {
	mu    sync.RWMutex
	items map[string]int
}

func New() *Inventory {
	return &Inventory{items: map[string]int{}}
}

func (inv *Inventory) Set(name string, qty int) {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	inv.items[name] = qty
}

func (inv *Inventory) Count(name string) int {
	inv.mu.RLock()
	defer inv.mu.RUnlock()
	return inv.items[name]
}

func (inv *Inventory) Total(names ...string) int {
	inv.mu.RLock()
	defer inv.mu.RUnlock()

	sum := 0
	for _, name := range names {
		sum += inv.Count(name)
	}
	return sum
}
