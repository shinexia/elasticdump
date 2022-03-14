/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package dumpdata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
)

const (
	QUERY_ALL = "*:*"
)

type DumpDataOption struct {
	Limit      int
	TimeoutSec int
}

func NewDumpDataOption() *DumpDataOption {
	return &DumpDataOption{
		Limit:      0,
		TimeoutSec: 60,
	}
}

func DumpData(client *elasticsearch.Client, dumpOption *DumpDataOption, writeFunc WriteDataFunc, o ...func(*esapi.SearchRequest)) (int, error) {
	res, err := client.Search(o...)
	count := 0
	for {
		if err != nil {
			return count, errors.Cause(err)
		}
		body, err := ioutil.ReadAll(res.Body)
		if res.IsError() {
			return count, errors.New(res.String())
		}
		response := &ScrollResponse{}
		err = json.Unmarshal(body, response)
		if err != nil {
			return count, errors.WithStack(err)
		}
		hits := response.Hits.Hits
		nhits := len(hits)
		if nhits == 0 {
			return count, nil
		}
		if dumpOption.Limit > 0 && count+nhits > dumpOption.Limit {
			hits = hits[:count+nhits-dumpOption.Limit]
		}
		n, err := writeFunc(hits)
		count += n
		if err != nil {
			return count, err
		}
		if dumpOption.Limit > 0 && count >= dumpOption.Limit {
			break
		}
		scrollReq := []byte(fmt.Sprintf(`{"scroll": "%ds","scroll_id": "%s"}`, dumpOption.TimeoutSec, response.ScrollID))
		res, err = client.Scroll(client.Scroll.WithContext(context.Background()), client.Scroll.WithBody(bytes.NewReader(scrollReq)))
	}
	return count, nil
}
