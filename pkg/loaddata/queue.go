/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package loaddata

import (
	"sync"
)

const minQueueCapacity = 16

// DataQueue a buffered FIFO data queue backed by a circular ring buffer
type DataQueue[T any] struct {
	mu      sync.Mutex
	cond    *sync.Cond
	buf     []T
	head    int
	count   int
	stopped bool
}

func NewDataQueue[T any]() *DataQueue[T] {
	q := &DataQueue[T]{
		buf: make([]T, minQueueCapacity),
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// grow doubles the buffer; must be called with lock held and count == len(buf).
func (q *DataQueue[T]) grow() {
	newBuf := make([]T, len(q.buf)*2)
	if q.head+q.count <= len(q.buf) {
		copy(newBuf, q.buf[q.head:q.head+q.count])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.count-n])
	}
	q.buf = newBuf
	q.head = 0
}

func (q *DataQueue[T]) Push(data ...T) bool {
	if len(data) == 0 {
		return true
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.stopped {
		return false
	}
	for _, item := range data {
		if q.count == len(q.buf) {
			q.grow()
		}
		q.buf[(q.head+q.count)%len(q.buf)] = item
		q.count++
	}
	q.cond.Broadcast()
	return true
}

func (q *DataQueue[T]) Stop() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.stopped = true
	q.cond.Broadcast()
}

func (q *DataQueue[T]) Pop(limit int) []T {
	q.mu.Lock()
	defer q.mu.Unlock()
	for {
		if q.count > 0 {
			n := q.count
			if limit > 0 && n > limit {
				n = limit
			}
			ret := make([]T, n)
			if q.head+n <= len(q.buf) {
				copy(ret, q.buf[q.head:q.head+n])
			} else {
				first := len(q.buf) - q.head
				copy(ret, q.buf[q.head:])
				copy(ret[first:], q.buf[:n-first])
			}
			// zero out slots to allow GC to collect referenced objects
			var zero T
			for i := 0; i < n; i++ {
				q.buf[(q.head+i)%len(q.buf)] = zero
			}
			q.head = (q.head + n) % len(q.buf)
			q.count -= n
			return ret
		}
		if q.stopped {
			return nil
		}
		q.cond.Wait()
	}
}
