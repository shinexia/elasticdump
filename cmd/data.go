/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package cmd

import (
	"io"
	"strings"
	"time"

	"github.com/shinexia/elasticdump/pkg/elasticdump"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func newCmdDumpData(out io.Writer) *cobra.Command {
	type DumpDataConfig struct {
		BaseConfig `json:",inline"`
		Batch      int `json:"batch"`
		Limit      int `json:"limit"`
		TimeoutSec int `json:"timeoutSec"`
	}
	cfg := &DumpDataConfig{
		BaseConfig: *newBaseConfig(),
		Batch:      1000,
		Limit:      -1,
		TimeoutSec: 60,
	}
	cmd := &cobra.Command{
		Use:   "data",
		Short: "dump data from elasticsearch",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			klog.V(5).Infof("cmd: %v\n", elasticdump.ToJSON(cfg))
			err = preprocessBaseConfig(&cfg.BaseConfig)
			if err != nil {
				return err
			}
			if cfg.File == "" {
				cfg.File = cfg.Index + "-data.json"
			}
			klog.V(5).Infof("cfg: %v\n", elasticdump.ToJSON(cfg))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dumper, err := newDumper(&cfg.BaseConfig)
			if err != nil {
				return err
			}
			return dumper.DumpData(cfg.Index, cfg.File, cfg.Batch, cfg.Limit, time.Duration(cfg.TimeoutSec)*time.Second)
		},
		Args: cobra.NoArgs,
	}

	addBaseConfigFlags(cmd.Flags(), &cfg.BaseConfig)
	flagSet := cmd.Flags()
	flagSet.IntVarP(&cfg.Batch, "batch", "b", cfg.Batch, "batch size when scroll")
	flagSet.IntVarP(&cfg.Limit, "limit", "l", cfg.Limit, "limit size when scroll")
	flagSet.IntVar(&cfg.TimeoutSec, "timeout", cfg.TimeoutSec, "timeout (second) when scroll")
	return cmd
}

func newCmdLoadData(out io.Writer) *cobra.Command {
	type LoadDataConfig struct {
		BaseConfig `json:",inline"`
		Batch      int  `json:"batch"`
		Limit      int  `json:"limit"`
		BufSize    int  `json:"bufSize"`
		Delete     bool `json:"delete"`
	}
	cfg := &LoadDataConfig{
		BaseConfig: *newBaseConfig(),
		Batch:      1000,
		Limit:      -1,
		BufSize:    1024 * 1024 * 1024,
		Delete:     false,
	}
	cmd := &cobra.Command{
		Use:   "data",
		Short: "load data to elasticsearch",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			klog.V(5).Infof("cmd: %v\n", elasticdump.ToJSON(cfg))
			err = preprocessBaseConfig(&cfg.BaseConfig)
			if err != nil {
				return err
			}
			if cfg.File == "" {
				cfg.File = cfg.Index + "-data.json"
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
			return dumper.LoadData(cfg.Index, cfg.File, cfg.Batch, cfg.Limit, cfg.BufSize)
		},
		Args: cobra.NoArgs,
	}
	addBaseConfigFlags(cmd.Flags(), &cfg.BaseConfig)
	flagSet := cmd.Flags()
	flagSet.IntVarP(&cfg.Batch, "batch", "b", cfg.Batch, "batch size when scroll")
	flagSet.IntVarP(&cfg.Limit, "limit", "l", cfg.Limit, "limit size when scroll")
	flagSet.IntVar(&cfg.BufSize, "buf", cfg.BufSize, "buffer size (byte) when split data file to lines, must bigger than the largest line")
	flagSet.BoolVar(&cfg.Delete, "delete", cfg.Delete, "whether delete the index before load")
	return cmd
}
