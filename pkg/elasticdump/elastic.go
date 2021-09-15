package elasticdump

import (
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/estransport"
	"github.com/pkg/errors"
)

// NewElasticSearchClient create elasticsearch.Client
func NewElasticSearchClient(host string) (*elasticsearch.Client, error) {
	conf := elasticsearch.Config{
		Addresses:         []string{host},
		EnableMetrics:     true,
		EnableDebugLogger: true,
		Logger: &estransport.TextLogger{
			Output:             os.Stdout,
			EnableRequestBody:  false,
			EnableResponseBody: false,
		},
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
