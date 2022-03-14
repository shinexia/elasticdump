/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package cmd

import (
	"io"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	"github.com/shinexia/elasticdump/pkg/helpers"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func newCmdDeleteIndex(out io.Writer) *cobra.Command {
	cfg := newBaseConfig()
	cmd := &cobra.Command{
		Use:   "index",
		Short: "delete index from elasticsearch",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			klog.V(5).Infof("cmd: %v\n", helpers.ToJSON(cfg))
			err = preprocessBaseConfig(cfg)
			if err != nil {
				return err
			}
			klog.V(5).Infof("cfg: %v\n", helpers.ToJSON(cfg))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			index := cfg.Index
			klog.V(5).Infof("deleting index: %s\n", index)
			startTime := time.Now()
			client, err := helpers.NewElasticSearchClient(cfg.Host)
			if err != nil {
				return err
			}
			res, err := client.Indices.Delete([]string{index})
			if err != nil {
				return errors.Cause(err)
			}
			if res.IsError() {
				return errors.New(res.String())
			}
			cost := time.Since(startTime).Seconds()
			klog.Infof("deleting index succeed, cost: %.3fs, index: %s, message: %s\n", cost, index, res.String())
			return nil
		},
		Args: cobra.NoArgs,
	}
	addBaseConfigFlags(cmd.Flags(), cfg)
	return cmd
}

func deleteIndexIfExists(client *elasticsearch.Client, index string) error {
	klog.V(5).Infof("deleting index: %s\n", index)
	startTime := time.Now()
	res, err := client.Indices.Delete([]string{index})
	if err != nil {
		if !strings.Contains(err.Error(), "index_not_found_exception") {
			return err
		}
		klog.Infof("%s\n", err.Error())
		return nil
	}
	if res.IsError() {
		return errors.New(res.String())
	}
	cost := time.Since(startTime).Seconds()
	klog.Infof("deleting index succeed, cost: %.3fs, index: %s, message: %s\n", cost, index, res.String())
	return nil
}
