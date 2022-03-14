/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package dumpdata

import (
	"encoding/json"
)

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
