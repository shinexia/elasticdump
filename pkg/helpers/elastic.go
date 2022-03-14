/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package helpers

import (
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/estransport"
	"github.com/pkg/errors"
	"k8s.io/klog"
)

// NewElasticSearchClient create elasticsearch.Client
func NewElasticSearchClient(host string) (*elasticsearch.Client, error) {
	conf := elasticsearch.Config{
		Addresses:         []string{host},
		EnableMetrics:     true,
		EnableDebugLogger: bool(klog.V(4)),
		Logger:            nil,
	}
	if klog.V(4) {
		conf.Logger = &estransport.TextLogger{
			Output:             os.Stdout,
			EnableRequestBody:  bool(klog.V(5)),
			EnableResponseBody: bool(klog.V(6)),
		}
	}
	client, err := elasticsearch.NewClient(conf)
	if err != nil {
		return nil, errors.Wrapf(err, "host=%s", host)
	}
	return client, nil
}

func PathJoin(a, b string) string {
	if strings.HasSuffix(a, "/") && strings.HasPrefix(b, "/") {
		return a + b[1:]
	} else if strings.HasSuffix(a, "/") || strings.HasPrefix(b, "/") {
		return a + b
	} else {
		return a + "/" + b
	}
}
