//go:build bughunt

package counter

type Counter struct {
	hits map[string]int
}

func New() *Counter {
	return &Counter{hits: map[string]int{}}
}

func (c *Counter) Inc(key string) {
	c.hits[key]++
}

func (c *Counter) Get(key string) int {
	return c.hits[key]
}
