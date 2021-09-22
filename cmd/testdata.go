/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package cmd

import (
	"io"
	"strings"

	"github.com/shinexia/elasticdump/pkg/elasticdump"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func newCmdGenTestData(out io.Writer) *cobra.Command {
	type GenTestDataConfig struct {
		BaseConfig `json:",inline"`
		Epoch      int
		Batch      int
		Delete     bool `json:"delete"`
	}
	cfg := &GenTestDataConfig{
		BaseConfig: *newBaseConfig(),
		Epoch:      2,
		Batch:      1000,
		Delete:     false,
	}
	cmd := &cobra.Command{
		Use:   "testdata",
		Short: "gen testdata to elasticsearch for test",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			klog.V(5).Infof("cmd: %v\n", elasticdump.ToJSON(cfg))
			err = preprocessBaseConfig(&cfg.BaseConfig)
			if err != nil {
				return err
			}
			klog.V(5).Infof("cfg: %v\n", elasticdump.ToJSON(cfg))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dumper, err := newDumper(&cfg.BaseConfig)
			if err != nil {
				return err
			}
			if cfg.Delete {
				err = dumper.DeleteIndex(cfg.Index)
				if err != nil {
					if !strings.Contains(err.Error(), "index_not_found_exception") {
						return err
					}
					klog.Infof("index: %s not found\n", cfg.Index)
				}
			}
			return dumper.GenTestData(cfg.Index, cfg.Epoch, cfg.Batch)
		},
		Args: cobra.NoArgs,
	}
	addBaseConfigFlags(cmd.Flags(), &cfg.BaseConfig)
	flagSet := cmd.Flags()
	flagSet.IntVarP(&cfg.Epoch, "epoch", "e", cfg.Epoch, "number of batch")
	flagSet.IntVarP(&cfg.Batch, "batch", "b", cfg.Batch, "batch size when send to elastic")
	flagSet.BoolVar(&cfg.Delete, "delete", cfg.Delete, "whether delete the index before load")
	return cmd
}
