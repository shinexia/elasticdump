package app

import (
	"net/url"

	"github.com/pkg/errors"
	"github.com/shinexia/elasticdump/pkg/elasticdump"
	flag "github.com/spf13/pflag"
)

type BaseConfig struct {
	Host  string `json:"host"`
	Index string `json:"index"`
	File  string `json:"file"`
}

func addBaseConfigFlags(flagSet *flag.FlagSet, cfg *BaseConfig) {
	flagSet.StringVar(&cfg.Host, "host", cfg.Host, "elasticsearch host: http://<user>:<password>@<host>:<port>")
	flagSet.StringVar(&cfg.Index, "index", cfg.Index, "elasticsearch index name")
	flagSet.StringVarP(&cfg.File, "file", "f", cfg.File, "filename")
}

func newBaseConfig() *BaseConfig {
	return &BaseConfig{
		Host:  "http://localhost:9200",
		Index: "elasticdumptest",
		File:  "",
	}
}

func parseHost(inHost string) (host string, err error) {
	if inHost == "" {
		return "", errors.Errorf("empty elasticsearch host")
	}
	up, err := url.Parse(inHost)
	if err != nil {
		return "", errors.Wrapf(err, "parse host url: %s failed", inHost)
	}
	if up.User != nil {
		password, _ := up.User.Password()
		host = up.Scheme + "://" + url.QueryEscape(up.User.Username()) + ":" + url.QueryEscape(password) + "@" + up.Host
	} else {
		host = up.Scheme + "://" + up.Host
	}
	return host, nil
}

func preprocessBaseConfig(old *BaseConfig) (*BaseConfig, error) {
	var cfg = *old
	host, err := parseHost(old.Host)
	if err != nil {
		return nil, err
	}
	cfg.Host = host
	return &cfg, nil
}

func newDumper(cfg *BaseConfig) (*elasticdump.Dumper, error) {
	es, err := elasticdump.NewElasticSearchClient(cfg.Host)
	if err != nil {
		return nil, err
	}
	client := elasticdump.NewESClient(cfg.Host, es)
	dumper := elasticdump.NewDumper(client)
	return dumper, nil
}
