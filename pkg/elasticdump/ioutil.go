/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package elasticdump

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

func ToJSON(obj interface{}) string {
	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Sprintf("%v", obj)
	}
	return string(data)
}

func ReadFileByLines(filename string, bufSize int, callback func(line []byte) (bool, error)) error {
	file, err := os.Open(filename)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, bufSize), bufSize)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Bytes()
		stop, err := callback(line)
		if err != nil {
			return err
		}
		if stop {
			return nil
		}
	}
	if scanner.Err() != nil {
		return errors.WithStack(scanner.Err())
	}
	return nil
}

// DataQueue a buffered FIFO data queue
type DataQueue struct {
	lock    sync.Mutex
	cond    *sync.Cond
	data    []interface{}
	stopped bool
	err     error
}

func NewDataQueue() *DataQueue {
	q := &DataQueue{}
	q.cond = sync.NewCond(&q.lock)
	return q
}

func (q *DataQueue) Push(data []interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if data == nil {
		q.stopped = true
	} else {
		q.data = append(q.data, data...)
	}
	q.cond.Broadcast()
}

func (q *DataQueue) PushError(err error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.err = err
	q.cond.Broadcast()
}

func (q *DataQueue) Pop(limit int) ([]interface{}, error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	for {
		if len(q.data) > 0 {
			if limit > 0 && len(q.data) > limit {
				ret := q.data[:limit]
				q.data = q.data[limit:]
				return ret, nil
			}
			ret := q.data
			q.data = nil
			return ret, nil
		}
		if q.err != nil {
			return nil, q.err
		}
		if q.stopped {
			return nil, nil
		}
		q.cond.Wait()
	}
}

func Bytes2Queue(data [][]byte) []interface{} {
	rs := make([]interface{}, len(data))
	for i, d := range data {
		rs[i] = d
	}
	return rs
}

func Hits2Queue(data []*Hit) []interface{} {
	rs := make([]interface{}, len(data))
	for i, d := range data {
		rs[i] = d
	}
	return rs
}

func Queue2Hits(data []interface{}) []*Hit {
	rs := make([]*Hit, len(data))
	for i, d := range data {
		rs[i] = d.(*Hit)
	}
	return rs
}

type AtomicBool struct{ flag int32 }

func (b *AtomicBool) Set(value bool) {
	var i int32 = 0
	if value {
		i = 1
	}
	atomic.StoreInt32(&(b.flag), int32(i))
}

func (b *AtomicBool) Get() bool {
	return atomic.LoadInt32(&(b.flag)) != 0
}
