package main

import "encoding/json"

type DataRecord struct {
	Index  string          `json:"_index"`
	Type   string          `json:"_type"`
	ID     string          `json:"_id"`
	Score  json.RawMessage `json:"_score"`
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
