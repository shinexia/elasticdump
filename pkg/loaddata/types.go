/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package loaddata

import (
	"encoding/json"
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
