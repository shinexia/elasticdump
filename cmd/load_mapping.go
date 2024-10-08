/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package cmd

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/shinexia/elasticdump/pkg/helpers"
	"github.com/shinexia/elasticdump/pkg/mapping"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func newCmdLoadMapping(_ io.Writer) *cobra.Command {
	type extraOption struct {
		InputFile string
		Delete    bool `json:"delete"`
	}
	cfg := newBaseConfig()
	extra := &extraOption{
		InputFile: "",
		Delete:    false,
	}
	cmd := &cobra.Command{
		Use:   "mapping",
		Short: "load mapping to elasticsearch",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			klog.V(5).Infof("cfg: %v, extra: %v", helpers.ToJSON(cfg), helpers.ToJSON(extra))
			err = preprocessBaseConfig(cfg)
			if err != nil {
				return err
			}
			if extra.InputFile == "" {
				extra.InputFile = cfg.Index + "-mapping.json"
			}
			klog.V(5).Infof("cfg: %v\n", helpers.ToJSON(cfg))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := helpers.NewElasticSearchClient(cfg.Host, cfg.InsecureSkipVerify)
			if err != nil {
				return err
			}
			if extra.Delete {
				err = deleteIndexIfExists(client, cfg.Index)
				if err != nil {
					return err
				}
			}
			inputFile := extra.InputFile
			klog.V(5).Infof("reading file: %s\n", inputFile)
			startTime := time.Now()
			mappingData, err := os.ReadFile(inputFile)
			if err != nil {
				return errors.Wrapf(err, "read file: %s failed", inputFile)
			}
			reqData, err := mapping.CleanUpMapping(string(mappingData))
			if err != nil {
				return err
			}
			res, err := client.Indices.Create(cfg.Index, client.Indices.Create.WithBody(bytes.NewReader([]byte(reqData))))
			if err != nil {
				return err
			}
			if res.IsError() {
				return errors.New(res.String())
			}
			cost := time.Since(startTime).Seconds()
			klog.Infof("load mapping succeed, cost: %.3fs, index: %s, file: %s, message: %s\n", cost, cfg.Index, inputFile, res)
			return nil
		},
		Args: cobra.NoArgs,
	}
	addBaseConfigFlags(cmd.Flags(), cfg)
	flagSet := cmd.Flags()

	flagSet.StringVarP(&extra.InputFile, "file", "f", extra.InputFile, "input file")
	flagSet.BoolVar(&extra.Delete, "delete", extra.Delete, "whether delete the index before load")

	return cmd
}
