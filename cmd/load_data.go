/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package cmd

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/shinexia/elasticdump/pkg/helpers"
	"github.com/shinexia/elasticdump/pkg/loaddata"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func newCmdLoadData(out io.Writer) *cobra.Command {
	type extraOption struct {
		InputFile string
		Batch     int
		Limit     int
		BufSize   int
		Delete    bool
	}
	cfg := newBaseConfig()
	extra := &extraOption{
		InputFile: "",
		Batch:     1000,
		Limit:     -1,
		BufSize:   1024 * 1024 * 1024,
		Delete:    false,
	}
	cmd := &cobra.Command{
		Use:   "data",
		Short: "load data to elasticsearch",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			klog.V(5).Infof("cmd: %v\n", helpers.ToJSON(cfg))
			err = preprocessBaseConfig(cfg)
			if err != nil {
				return err
			}
			if extra.InputFile == "" {
				extra.InputFile = cfg.Index + "-data.json"
			}
			klog.V(5).Infof("cfg: %v\n", helpers.ToJSON(cfg))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := helpers.NewElasticSearchClient(cfg.Host)
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
			klog.V(5).Infof("load data to index: %s, from: %s, batch: %v, limit: %v, bufSize: %v\n", cfg.Index, inputFile, extra.Batch, extra.Limit, extra.BufSize)
			queue := loaddata.NewDataQueue[*loaddata.Hit]()
			var rerr error
			//  async read records from file
			go func() {
				defer queue.Stop()
				file, err := os.Open(inputFile)
				if err != nil {
					rerr = errors.WithMessagef(err, "file: %s", inputFile)
					return
				}
				defer file.Close()
				err = loaddata.LoadHits(queue, file, extra.BufSize)
				if err != nil {
					rerr = err
				}
			}()
			err = loaddata.LoadData(client, queue, extra.Batch, cfg.Index)
			if err != nil {
				return err
			}
			return rerr
		},
		Args: cobra.NoArgs,
	}
	addBaseConfigFlags(cmd.Flags(), cfg)
	flagSet := cmd.Flags()

	flagSet.StringVarP(&extra.InputFile, "file", "f", extra.InputFile, "input file")

	flagSet.IntVarP(&extra.Batch, "batch", "b", extra.Batch, "batch size when scroll")
	flagSet.IntVarP(&extra.Limit, "limit", "l", extra.Limit, "limit size when scroll")
	flagSet.IntVar(&extra.BufSize, "buf", extra.BufSize, "buffer size (byte) when split data file to lines, must bigger than the largest line")
	flagSet.BoolVar(&extra.Delete, "delete", extra.Delete, "whether delete the index before load")
	return cmd
}
