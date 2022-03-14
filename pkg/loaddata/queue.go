/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package loaddata

import (
	"sync"
)

// DataQueue a buffered FIFO data queue
type DataQueue[T any] struct {
	lock    sync.Mutex
	cond    *sync.Cond
	data    []T
	stopped bool
}

func NewDataQueue[T any]() *DataQueue[T] {
	q := &DataQueue[T]{}
	q.cond = sync.NewCond(&q.lock)
	return q
}

func (q *DataQueue[T]) Push(data ...T) bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.stopped {
		return false
	}
	q.data = append(q.data, data...)
	if len(q.data) > 0 {
		q.cond.Broadcast()
	}
	return true
}

func (q *DataQueue[T]) Stop() {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.stopped = true
	q.cond.Broadcast()
}

func (q *DataQueue[T]) Pop(limit int) []T {
	q.lock.Lock()
	defer q.lock.Unlock()
	for {
		if len(q.data) > 0 {
			if limit > 0 && len(q.data) > limit {
				ret := q.data[:limit]
				q.data = q.data[limit:]
				return ret
			}
			ret := q.data
			q.data = nil
			return ret
		}
		if q.stopped {
			return nil
		}
		q.cond.Wait()
	}
}
