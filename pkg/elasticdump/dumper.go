package elasticdump

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
	"k8s.io/klog"
)

type Dumper struct {
	client *ESClient
}

func NewDumper(client *ESClient) *Dumper {
	return &Dumper{
		client: client,
	}
}

func (d *Dumper) DumpMapping(index string, dest string) error {
	startTime := time.Now()
	res, err := d.client.DumpMapping(index)
	if err != nil {
		return err
	}
	klog.V(5).Infof("writing mapping to: %s\n", dest)
	writer, err := os.Create(dest)
	if err != nil {
		return errors.Wrapf(err, "create file: %s failed", dest)
	}
	defer writer.Close()
	_, err = writer.WriteString(res)
	if err != nil {
		return errors.Wrapf(err, "dest: %s", dest)
	}
	cost := time.Since(startTime).Seconds()
	klog.Infof("dump mapping succeed, cost: %.3fs, index: %s\n", cost, index)
	return nil
}

func (d *Dumper) DeleteIndex(index string) error {
	klog.V(5).Infof("deleting index: %s\n", index)
	startTime := time.Now()
	res, err := d.client.DeleteIndex(index)
	if err != nil {
		return err
	}
	cost := time.Since(startTime).Seconds()
	klog.Infof("deleting index succeed, cost: %.3fs, index: %s, message: %s\n", cost, index, res)
	return nil
}

func (d *Dumper) LoadMapping(index string, filename string) error {
	klog.V(5).Infof("reading file: %s\n", filename)
	startTime := time.Now()
	mappingData, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrapf(err, "read file: %s failed", filename)
	}
	reqData, err := d.client.CleanUpMapping(string(mappingData))
	if err != nil {
		return err
	}
	res, err := d.client.LoadMapping(index, reqData)
	if err != nil {
		return err
	}
	cost := time.Since(startTime).Seconds()
	klog.Infof("load mapping succeed, cost: %.3fs, index: %s, file: %s, message: %s\n", cost, index, filename, res)
	return nil
}

func (d *Dumper) DumpData(index string, filename string, batch int, limit int, timeout time.Duration) error {
	klog.V(5).Infof("dump data from index: %s, to: %s, batch: %v, limit: %v\n", index, filename, batch, limit)
	startTime := time.Now()
	queue := NewDataQueue()
	stopped := &AtomicBool{}
	// async read records from elasticsearch
	go func() {
		err := d.client.DumpData(index, batch, limit, timeout, func(hits [][]byte) (bool, error) {
			if stopped.Get() {
				return true, nil
			}
			klog.V(5).Infof("recieved hits: %d\n", len(hits))
			queue.Push(Bytes2Queue(hits))
			return false, nil
		})
		if err != nil {
			queue.PushError(err)
		} else {
			queue.Push(nil)
		}
	}()
	var writer *os.File
	defer func() {
		if writer != nil {
			writer.Close()
		}
	}()
	// pop and write records to file
	newLine := []byte("\n")
	numWrited := 0
	for {
		hits, err := queue.Pop(-1)
		if err != nil {
			return err
		}
		if len(hits) == 0 {
			break
		}
		if writer == nil {
			var err error
			writer, err = os.Create(filename)
			if err != nil {
				stopped.Set(true)
				return errors.WithStack(err)
			}
			klog.V(5).Infof("created file: %s\n", filename)
		}
		startTime2 := time.Now()
		for _, hit := range hits {
			_, err = writer.Write(hit.([]byte))
			if err == nil {
				_, err = writer.Write(newLine)
			}
			if err != nil {
				stopped.Set(true)
				return errors.WithStack(err)
			}
			numWrited++
		}
		cost2 := time.Since(startTime2).Seconds()
		klog.Infof("writed: %d/%d, cost: %.3fs\n", numWrited, len(hits), cost2)
	}
	cost := time.Since(startTime).Seconds()
	klog.Infof("dump data succeed, total: %d, index: %s, file: %s, cost: %.3fs\n", numWrited, index, filename, cost)
	return nil
}

func (d *Dumper) LoadData(index string, filename string, batch int, limit int, bufSize int) error {
	klog.V(5).Infof("load data to index: %s, from: %s, batch: %v, limit: %v, bufSize: %v\n", index, filename, batch, limit, bufSize)
	queue := NewDataQueue()
	stopped := &AtomicBool{}
	//  async read records from file
	go func() {
		err := ReadFileByLines(filename, bufSize, func(line []byte) (bool, error) {
			if stopped.Get() {
				return true, nil
			}
			if len(line) == 0 {
				return false, nil
			}
			hit := &Hit{}
			err := json.Unmarshal(line, hit)
			if err != nil {
				stopped.Set(true)
				return false, errors.Wrapf(err, "line: %s", string(line))
			}
			queue.Push([]interface{}{hit})
			return false, nil
		})
		if err != nil {
			queue.PushError(err)
		} else {
			queue.Push(nil)
		}
	}()
	return d.doLoadData(queue, stopped, index, filename, batch)
}

func (d *Dumper) GenTestData(index string, epoch, batch int) error {
	klog.V(5).Infof("gen test data to index: %s, epoch: %d, batch: %d\n", index, epoch, batch)
	queue := NewDataQueue()
	stopped := &AtomicBool{}
	//  async read records from file
	go func() {
		err := GenerateTestData(epoch, batch, func(hits []*Hit) (bool, error) {
			if stopped.Get() {
				return true, nil
			}
			queue.Push(Hits2Queue(hits))
			return false, nil
		})
		if err != nil {
			queue.PushError(err)
		} else {
			queue.Push(nil)
		}
	}()
	return d.doLoadData(queue, stopped, index, "testgen", batch)
}

func (d *Dumper) doLoadData(queue *DataQueue, stopped *AtomicBool, index string, filename string, batch int) error {
	startTime := time.Now()
	// pop and send records to elasticsearch
	totalSecceed := 0
	totalError := 0
	for {
		hitsQ, err := queue.Pop(batch)
		if err != nil {
			return err
		}
		if len(hitsQ) == 0 {
			break
		}
		klog.V(5).Infof("received hits: %v\n", len(hitsQ))
		startTime2 := time.Now()
		hits := Queue2Hits(hitsQ)
		resData, err := d.client.LoadData(index, hits)
		if err != nil {
			stopped.Set(true)
			return err
		}
		var res = &BulkResponse{}
		err = json.Unmarshal([]byte(resData), res)
		if err != nil {
			stopped.Set(true)
			return err
		}
		succeedCount := 0
		errorCount := 0
		for _, d := range res.Items {
			// ... so for any HTTP status above 201 ...
			if d.Create.Status > 201 {
				errorCount++
				klog.Infof("Error [%d]: %s\n", d.Create.Status, d.Create.Error)
			} else {
				succeedCount++
			}
		}
		totalSecceed += succeedCount
		totalError += errorCount
		cost2 := time.Since(startTime2).Seconds()
		klog.Infof("indexed succeed: %v/%v, failed: %v/%v, cost: %.3fs\n", totalSecceed, succeedCount, totalError, errorCount, cost2)
	}
	cost := time.Since(startTime).Seconds()
	klog.Infof("load data succeed, indexed: %d, failed: %v, index: %s, file: %s, cost: %.3fs\n", totalSecceed, totalError, index, filename, cost)
	return nil
}
