/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package cmd

import (
	"io"

	"github.com/shinexia/elasticdump/pkg/helpers"
	"github.com/shinexia/elasticdump/pkg/loaddata"

	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func newCmdTestGenData(out io.Writer) *cobra.Command {
	type extraOption struct {
		Epoch  int
		Batch  int
		Delete bool
	}
	cfg := newBaseConfig()
	extra := &extraOption{
		Epoch:  2,
		Batch:  1000,
		Delete: false,
	}
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "gen testdata to elasticsearch for test",
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
			klog.V(5).Infof("gen test data to index: %s, epoch: %v, batch: %v\n", cfg.Index, extra.Epoch, extra.Batch)
			queue := loaddata.NewDataQueue[*loaddata.Hit]()
			var rerr error
			//  async read records from file
			go func() {
				defer queue.Stop()
				err = loaddata.GenTestHits(queue, extra.Epoch, extra.Batch)
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
	flagSet.IntVarP(&extra.Epoch, "epoch", "e", extra.Epoch, "number of batch")
	flagSet.IntVarP(&extra.Batch, "batch", "b", extra.Batch, "batch size when send to elastic")
	flagSet.BoolVar(&extra.Delete, "delete", extra.Delete, "whether delete the index before load")
	return cmd
}
