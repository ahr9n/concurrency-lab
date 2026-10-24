//go:build bughunt

package queue

import "sync"

type Queue struct {
	mu    sync.Mutex
	cond  *sync.Cond
	items []int
}

func New() *Queue {
	q := &Queue{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func (q *Queue) Put(v int) {
	q.mu.Lock()
	q.items = append(q.items, v)
	q.mu.Unlock()
	q.cond.Broadcast()
}

func (q *Queue) Take() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		q.cond.Wait()
	}
	v := q.items[0]
	q.items = q.items[1:]
	return v
}
