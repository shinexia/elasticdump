package app

import (
	"io"
	"log"

	"github.com/shinexia/elasticdump/pkg/elasticdump"
	"github.com/spf13/cobra"
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
			log.Printf("origin: %v\n", elasticdump.ToJSON(cfg))
			err = preprocessBaseConfig(&cfg.BaseConfig)
			if err != nil {
				return err
			}
			log.Printf("cfg: %v\n", elasticdump.ToJSON(cfg))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dumper, err := newDumper(&cfg.BaseConfig)
			if err != nil {
				return err
			}
			return dumper.DeleteIndex(cfg.Index)
		},
		Args: cobra.NoArgs,
	}
	addBaseConfigFlags(cmd.Flags(), &cfg.BaseConfig)
	return cmd
}
