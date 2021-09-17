package elasticdump

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"k8s.io/klog"
)

type ESClient struct {
	host   string // format: http://<username>:<password>@host:port
	client *elasticsearch.Client
}

// NewESClient create ESClient
//  dump mapping use `url` as address, `host` as address
func NewESClient(host string, client *elasticsearch.Client) *ESClient {
	return &ESClient{
		host:   host,
		client: client,
	}
}

func (ec *ESClient) makeResponse(res *esapi.Response, err error) (string, error) {
	if err != nil {
		return "", errors.WithStack(err)
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrapf(err, "read response body failed")
	}
	res.Body.Close()
	if res.IsError() {
		return "", errors.Errorf("status: %d, body: %s", res.StatusCode, string(resBody))
	}
	return string(resBody), nil
}

func (ec *ESClient) DumpMapping(index string) (string, error) {
	client, err := NewElasticSearchClient(PathJoin(ec.host, index))
	if err != nil {
		return "", err
	}
	res, err := client.Info(func(r *esapi.InfoRequest) {
		r.Pretty = true
		r.Human = true
	})
	return ec.makeResponse(res, err)
}

func (ec *ESClient) DeleteIndex(index string) (string, error) {
	client := ec.client
	res, err := client.Indices.Delete([]string{index})
	return ec.makeResponse(res, err)
}

func (ec *ESClient) LoadMapping(index string, data string) (string, error) {
	client := ec.client
	res, err := client.Indices.Create(index, client.Indices.Create.WithBody(bytes.NewReader([]byte(data))))
	return ec.makeResponse(res, err)
}

// cleanupMapping clean up some unnecessary field
func (ec *ESClient) CleanUpMapping(data string) (string, error) {
	var rootMap map[string]json.RawMessage
	err := json.Unmarshal([]byte(data), &rootMap)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if len(rootMap) != 1 {
		return "", errors.Errorf("multiple indexes: %v", len(rootMap))
	}
	var indexData json.RawMessage
	for _, v := range rootMap {
		indexData = v
	}
	var dataMap map[string]json.RawMessage
	err = json.Unmarshal(indexData, &dataMap)
	if err != nil {
		return "", errors.WithStack(err)
	}
	settingsData := dataMap["settings"]
	var settingsMap map[string]json.RawMessage
	err = json.Unmarshal(settingsData, &settingsMap)
	if err != nil {
		return "", errors.WithStack(err)
	}
	indexSettingData := settingsMap["index"]
	var indexSettingMap map[string]json.RawMessage
	err = json.Unmarshal(indexSettingData, &indexSettingMap)
	if err != nil {
		return "", errors.WithStack(err)
	}
	// delete .settings.index unused fields
	for _, key := range []string{"creation_date", "uuid", "version", "provided_name", "routing", "creation_date_string"} {
		delete(indexSettingMap, key)
	}
	newIndexSettingData, err := json.Marshal(indexSettingMap)
	if err != nil {
		return "", errors.WithStack(err)
	}
	settingsMap["index"] = json.RawMessage(newIndexSettingData)
	dataMap["settings"], err = json.Marshal(settingsMap)
	if err != nil {
		return "", errors.WithStack(err)
	}
	newData, err := json.MarshalIndent(dataMap, "", "  ")
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(newData), nil
}

func (ec *ESClient) ScrollStart(index string, batch int, timeout time.Duration) (string, error) {
	client := ec.client
	res, err := client.Search(client.Search.WithContext(context.Background()), client.Search.WithIndex(index), client.Search.WithScroll(timeout), client.Search.WithSize(batch), client.Search.WithQuery("*:*"))
	return ec.makeResponse(res, err)
}

func (ec *ESClient) ScrollNext(scrollId string, timeout time.Duration) (string, error) {
	client := ec.client
	req := []byte(fmt.Sprintf(`{"scroll": "%.0fs","scroll_id": "%s"}`, timeout.Seconds(), scrollId))
	res, err := client.Scroll(client.Scroll.WithContext(context.Background()), client.Scroll.WithBody(bytes.NewReader(req)))
	return ec.makeResponse(res, err)
}

func (ec *ESClient) DumpData(index string, batch int, limit int, timeout time.Duration, callback func(hit [][]byte) (bool, error)) error {
	if callback == nil {
		return errors.Errorf("callback is nil")
	}
	resData, err := ec.ScrollStart(index, batch, timeout)
	if err != nil {
		return err
	}
	count := 0
	for {
		klog.V(9).Infof("resData:\n%s\n", resData)
		var res = &ScrollResponse{}
		err = json.Unmarshal([]byte(resData), res)
		if err != nil {
			return errors.WithStack(err)
		}
		hits := res.Hits.Hits
		count += len(hits)
		if limit > 0 && count > limit {
			hits = hits[:count-limit]
		}
		hitsData := make([][]byte, len(hits))
		for i := range hits {
			hitsData[i] = hits[i]
		}
		stop, err := callback(hitsData)
		if err != nil {
			return err
		}
		if stop {
			break
		}
		if limit > 0 && count > limit {
			break
		}
		if len(res.Hits.Hits) < batch {
			break
		}
		resData, err = ec.ScrollNext(res.ScrollID, timeout)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ec *ESClient) LoadData(index string, items []*Hit) (string, error) {
	client := ec.client
	var buf bytes.Buffer
	for _, r := range items {
		meta := []byte(fmt.Sprintf(`{"create": {"_index": "%s", "_type": "_doc", "_id": "%s"}}%s`, index, r.ID, "\n"))
		buf.Write(meta)
		buf.Write(r.Source)
		buf.Write([]byte("\n"))
	}
	data := buf.Bytes()
	res, err := client.Bulk(bytes.NewReader(data))
	return ec.makeResponse(res, err)
}
