/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"github.com/shinexia/elasticdump/pkg/helpers"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func newCmdDumpMapping(out io.Writer) *cobra.Command {
	type extraOption struct {
		OutputFile string
	}
	cfg := newBaseConfig()
	extra := &extraOption{}
	cmd := &cobra.Command{
		Use:   "mapping",
		Short: "dump mapping from elasticsearch",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			klog.V(5).Infof("cmd: %v\n", helpers.ToJSON(cfg))
			err = preprocessBaseConfig(cfg)
			if err != nil {
				return err
			}
			if extra.OutputFile == "" {
				extra.OutputFile = cfg.Index + "-mapping.json"
			}
			klog.V(5).Infof("cfg: %v\n", helpers.ToJSON(cfg))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			startTime := time.Now()
			client, err := helpers.NewElasticSearchClient(helpers.PathJoin(cfg.Host, cfg.Index))
			if err != nil {
				return err
			}
			res, err := client.Info(func(r *esapi.InfoRequest) {
				r.Pretty = true
			}, client.Info.WithHuman())
			if err != nil {
				return errors.Cause(err)
			}
			body, err := ioutil.ReadAll(res.Body)
			if res.IsError() || err != nil {
				return errors.Errorf("status: %d, body: %s", res.StatusCode, string(body))
			}
			outputFile := extra.OutputFile
			klog.V(5).Infof("writing mapping to: %s\n", outputFile)
			writer, err := os.Create(outputFile)
			if err != nil {
				return errors.Wrapf(err, "create file: %s failed", outputFile)
			}
			defer writer.Close()
			_, err = writer.Write(body)
			if err != nil {
				return errors.Wrapf(err, "dest: %s", outputFile)
			}
			cost := time.Since(startTime).Seconds()
			klog.Infof("dump mapping succeed, cost: %.3fs, index: %s, file: %s", cost, cfg.Index, outputFile)
			return nil
		},
		Args: cobra.NoArgs,
	}
	addBaseConfigFlags(cmd.Flags(), cfg)
	flagSet := cmd.Flags()

	flagSet.StringVarP(&extra.OutputFile, "file", "f", extra.OutputFile, "output file")

	return cmd
}
