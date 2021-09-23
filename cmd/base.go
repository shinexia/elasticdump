/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package cmd

import (
	gflag "flag"
	"net/url"
	"os"

	"github.com/pkg/errors"
	"github.com/shinexia/elasticdump/pkg/elasticdump"
	flag "github.com/spf13/pflag"
	"k8s.io/klog"
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

	klogSet := gflag.NewFlagSet(os.Args[0], gflag.ContinueOnError)
	klog.InitFlags(klogSet)

	flagSet.AddGoFlagSet(klogSet)
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
	if up.Scheme == "" || up.Host == "" {
		return "", errors.Errorf("invalid host url: %s", inHost)
	}
	if up.User != nil {
		password, _ := up.User.Password()
		host = up.Scheme + "://" + url.QueryEscape(up.User.Username()) + ":" + url.QueryEscape(password) + "@" + up.Host
	} else {
		host = up.Scheme + "://" + up.Host
	}
	return host, nil
}

func preprocessBaseConfig(cfg *BaseConfig) error {
	host, err := parseHost(cfg.Host)
	if err != nil {
		return err
	}
	cfg.Host = host
	return nil
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
