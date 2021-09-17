/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package elasticdump

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Hit struct {
	ID     string          `json:"_id"`
	Source json.RawMessage `json:"_source"`
}

type BulkResponse struct {
	Took   int  `json:"took"`
	Errors bool `json:"errors"`
	Items  []struct {
		Create struct {
			Index  string          `json:"_index"`
			Type   string          `json:"_type"`
			ID     string          `json:"_id"`
			Status int             `json:"status"`
			Result string          `json:"result"`
			Error  json.RawMessage `json:"error"`
		} `json:"create"`
	} `json:"items"`
}

type ScrollResponse struct {
	ScrollID string          `json:"_scroll_id"`
	Took     int             `json:"took"`
	TimeOut  bool            `json:"time_out"`
	Shards   json.RawMessage `json:"_shards"`
	Hits     struct {
		Total    json.RawMessage   `json:"total"`
		MaxScore float32           `json:"max_score"`
		Hits     []json.RawMessage `json:"hits"`
	} `json:"hits"`
}

// TestData test data structure
type TestData struct {
	Content   string `json:"content"`
	Title     string `json:"title"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func GenerateTestData(epoch, batch int, callback func(hits []*Hit) (bool, error)) error {
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
		stop, err := callback(items)
		if err != nil {
			return err
		}
		if stop {
			return nil
		}
	}
	return nil
}
