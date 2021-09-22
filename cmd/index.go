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

func newCmdDeleteIndex(out io.Writer) *cobra.Command {
	type DeleteConfig struct {
		BaseConfig `json:",inline"`
	}
	cfg := &DeleteConfig{
		BaseConfig: *newBaseConfig(),
	}
	cmd := &cobra.Command{
		Use:   "index",
		Short: "delete index from elasticsearch",
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
			err = dumper.DeleteIndex(cfg.Index)
			if err != nil {
				if !strings.Contains(err.Error(), "index_not_found_exception") {
					return err
				}
				klog.Infof("%s\n", err.Error())
			}
			return nil
		},
		Args: cobra.NoArgs,
	}
	addBaseConfigFlags(cmd.Flags(), &cfg.BaseConfig)
	return cmd
}
