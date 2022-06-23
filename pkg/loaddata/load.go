/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package loaddata

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	"k8s.io/klog"
)

func LoadHits(queue *DataQueue[*Hit], in io.Reader, maxLineLength int) error {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, maxLineLength), maxLineLength)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		hit := &Hit{}
		err := json.Unmarshal(line, hit)
		if err != nil {
			return errors.WithStack(err)
		}
		ok := queue.Push(hit)
		if !ok {
			return nil
		}
	}
	if scanner.Err() != nil {
		return errors.WithStack(scanner.Err())
	}
	return nil
}

func GenTestHits(queue *DataQueue[*Hit], epoch, batch int) error {
	type TestData struct {
		Content   string `json:"content"`
		Title     string `json:"title"`
		CreatedAt int64  `json:"created_at"`
		UpdatedAt int64  `json:"updated_at"`
	}
	id := 1
	for i := 0; i < epoch; i++ {
		items := make([]*Hit, batch)
		for j := 0; j < batch; j++ {
			now := time.Now().UnixMilli()
			data := &TestData{
				Content:   "content-" + strconv.Itoa(id),
				Title:     "title-" + strconv.Itoa(id),
				CreatedAt: now,
				UpdatedAt: now,
			}
			dj, err := json.Marshal(data)
			if err != nil {
				return errors.WithStack(err)
			}
			items[j] = &Hit{
				ID:     "id-" + strconv.Itoa(id),
				Source: dj,
			}
			id++
		}
		ok := queue.Push(items...)
		if !ok {
			return nil
		}
	}
	return nil
}

func LoadData(client *elasticsearch.Client, queue *DataQueue[*Hit], batch int, index string) error {
	startTime := time.Now()
	// pop and send records to elasticsearch
	totalSecceed := 0
	totalError := 0
	for {
		hits := queue.Pop(batch)
		if len(hits) == 0 {
			break
		}
		klog.V(5).Infof("received lines: %v\n", len(hits))
		startTime2 := time.Now()
		var buf bytes.Buffer
		for _, r := range hits {
			meta := []byte(fmt.Sprintf(`{"create": {"_index": "%s", "_type": "_doc", "_id": "%s"}}%s`, index, r.ID, "\n"))
			// have routing field
			if len(r.Routing) > 0 {
				meta = []byte(fmt.Sprintf(`{"create": {"_index": "%s", "_type": "_doc", "_id": "%s", "routing": "%s"}}%s`, index, r.ID, r.Routing, "\n"))
			}
			buf.Write(meta)
			buf.Write(r.Source)
			buf.Write([]byte("\n"))
		}
		res, err := client.Bulk(bytes.NewReader(buf.Bytes()))
		if err != nil {
			return errors.WithStack(err)
		}
		body, err := ioutil.ReadAll(res.Body)
		if res.IsError() || err != nil {
			return errors.Errorf("status: %d, body: %s", res.StatusCode, string(body))
		}
		var result = &BulkResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return errors.WithMessagef(err, "status: %d, body: %s", res.StatusCode, string(body))
		}
		succeedCount := 0
		errorCount := 0
		for _, d := range result.Items {
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
		klog.Infof("indexed succeed: %v/%v, failed: %v/%v, cost: %.3fs\n", succeedCount, totalSecceed, errorCount, totalError, cost2)
	}
	cost := time.Since(startTime).Seconds()
	klog.Infof("load data succeed, indexed: %d, failed: %v, index: %s, cost: %.3fs\n", totalSecceed, totalError, index, cost)
	return nil
}
