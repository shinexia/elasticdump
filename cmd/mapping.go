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

func newCmdDumpMapping(out io.Writer) *cobra.Command {
	type DumpMappingConfig struct {
		BaseConfig `json:",inline"`
	}
	cfg := &DumpMappingConfig{
		BaseConfig: *newBaseConfig(),
	}
	cmd := &cobra.Command{
		Use:   "mapping",
		Short: "dump mapping from elasticsearch",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			klog.V(5).Infof("cmd: %v\n", elasticdump.ToJSON(cfg))
			err = preprocessBaseConfig(&cfg.BaseConfig)
			if err != nil {
				return err
			}
			if cfg.File == "" {
				cfg.File = cfg.Index + "-mapping.json"
			}
			klog.V(5).Infof("cfg: %v\n", elasticdump.ToJSON(cfg))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dumper, err := newDumper(&cfg.BaseConfig)
			if err != nil {
				return err
			}
			return dumper.DumpMapping(cfg.Index, cfg.File)
		},
		Args: cobra.NoArgs,
	}
	addBaseConfigFlags(cmd.Flags(), &cfg.BaseConfig)
	return cmd
}

func newCmdLoadMapping(out io.Writer) *cobra.Command {
	type LoadMappingConfig struct {
		BaseConfig `json:",inline"`
		Delete     bool `json:"delete"`
	}
	cfg := &LoadMappingConfig{
		BaseConfig: *newBaseConfig(),
		Delete:     false,
	}
	cmd := &cobra.Command{
		Use:   "mapping",
		Short: "load mapping to elasticsearch",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			klog.V(5).Infof("cmd: %v\n", elasticdump.ToJSON(cfg))
			err = preprocessBaseConfig(&cfg.BaseConfig)
			if err != nil {
				return err
			}
			if cfg.File == "" {
				cfg.File = cfg.Index + "-mapping.json"
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
			return dumper.LoadMapping(cfg.Index, cfg.File)
		},
		Args: cobra.NoArgs,
	}
	addBaseConfigFlags(cmd.Flags(), &cfg.BaseConfig)
	flagSet := cmd.Flags()
	flagSet.BoolVar(&cfg.Delete, "delete", cfg.Delete, "whether delete the index before load")
	return cmd
}
