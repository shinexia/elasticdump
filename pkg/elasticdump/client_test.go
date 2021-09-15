package elasticdump

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

const (
	testHost      = "http://localhost:9200"
	testIndexName = "elasticdumptest"
)

func newTestClient() (*ESClient, error) {
	host := testHost
	client, err := NewElasticSearchClient(host)
	if err != nil {
		return nil, err
	}
	return NewESClient(host, client), nil
}

func TestGetMapping(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Error(err)
		return
	}
	res, err := client.DumpMapping(testIndexName)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("res: \n%v\n", string(res))
}

func TestCleanUpMapping(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Error(err)
		return
	}
	res, err := client.DumpMapping(testIndexName)
	if err != nil {
		t.Error(err)
		return
	}
	clean, err := client.CleanUpMapping(res)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("res: \n%v\n", string(clean))
}

func TestLoadMapping(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/test-mapping.json")
	if err != nil {
		t.Error(err)
		return
	}
	client, err := newTestClient()
	if err != nil {
		t.Error(err)
		return
	}
	reqData, err := client.CleanUpMapping(string(data))
	if err != nil {
		t.Error(err)
		return
	}
	res, err := client.LoadMapping(testIndexName+"2", reqData)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("res: \n%v\n", string(res))
}

func TestDeleteIndex(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Error(err)
		return
	}
	res, err := client.DeleteIndex(testIndexName + "2")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("res: \n%v\n", string(res))
}

func TestScrollStart(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Error(err)
		return
	}
	res, err := client.ScrollStart(testIndexName, 10, time.Minute)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("res: \n%v\n", string(res))
}

func TestScrollNext(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Error(err)
		return
	}
	res, err := client.ScrollNext("FGluY2x1ZGVfY29udGV4dF91dWlkDXF1ZXJ5QW5kRmV0Y2gBFmJ4TlZkQU9fU3pTVDRmMTN6QW90TGcAAAAAAAAgQhZKWkM2R1NtclRKYS1EMzVtZUFlT0l3", time.Minute)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("res: \n%v\n", string(res))
}

func TestScroll(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Error(err)
		return
	}
	err = client.DumpData(testIndexName, 100, -1, time.Minute, func(hits [][]byte) (bool, error) {
		for _, hit := range hits {
			fmt.Printf("hit: %s\n", string(hit))
		}
		return false, nil
	})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("succeed\n")
}
