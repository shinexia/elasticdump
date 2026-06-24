/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package helpers

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/pkg/errors"
	"k8s.io/klog"
)

// NewElasticSearchClient create elasticsearch.Client
func NewElasticSearchClient(host string, insecureSkipVerify bool) (*elasticsearch.Client, error) {
	options := []elasticsearch.Option{
		elasticsearch.WithAddresses(host),
	}
	if insecureSkipVerify {
		options = append(options, elasticsearch.WithTransportOptions(
			elastictransport.WithTransport(&http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}),
		))
	}
	if klog.V(4) {
		options = append(options, elasticsearch.WithLogger(
			&elastictransport.TextLogger{
				Output:             os.Stdout,
				EnableRequestBody:  bool(klog.V(5)),
				EnableResponseBody: bool(klog.V(6)),
			},
		))
	}
	client, err := elasticsearch.New(options...)
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
