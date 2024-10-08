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
	flag "github.com/spf13/pflag"
	"k8s.io/klog"
)

type BaseConfig struct {
	Host               string `json:"host"`
	Index              string `json:"index"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify"`
}

func addBaseConfigFlags(flagSet *flag.FlagSet, cfg *BaseConfig) {
	flagSet.StringVar(&cfg.Host, "host", cfg.Host, "elasticsearch host: http://<user>:<password>@<host>:<port>")
	flagSet.StringVar(&cfg.Index, "index", cfg.Index, "elasticsearch index name")
	flagSet.BoolVarP(&cfg.InsecureSkipVerify, "insecure-skip-verify", "k", cfg.InsecureSkipVerify, "skip verify tls certificate")

	klogSet := gflag.NewFlagSet(os.Args[0], gflag.ContinueOnError)
	klog.InitFlags(klogSet)

	flagSet.AddGoFlagSet(klogSet)
}

func newBaseConfig() *BaseConfig {
	return &BaseConfig{
		Host:               "http://localhost:9200",
		Index:              "elasticdumptest",
		InsecureSkipVerify: false,
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
