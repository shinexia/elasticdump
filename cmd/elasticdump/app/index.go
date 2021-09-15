package app

import (
	"io"
	"log"

	"github.com/shinexia/elasticdump/pkg/elasticdump"
	"github.com/spf13/cobra"
)

func newCmdDeleteIndex(out io.Writer) *cobra.Command {
	baseCfg := newBaseConfig()
	cmd := &cobra.Command{
		Use:   "index",
		Short: "delete index from elasticsearch",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			log.Printf("origin: %v\n", elasticdump.ToJSON(baseCfg))
			baseCfg, err = preprocessBaseConfig(baseCfg)
			if err != nil {
				return err
			}
			log.Printf("cfg: %v\n", elasticdump.ToJSON(baseCfg))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dumper, err := newDumper(baseCfg)
			if err != nil {
				return err
			}
			return dumper.DeleteIndex(baseCfg.Index)
		},
		Args: cobra.NoArgs,
	}
	addBaseConfigFlags(cmd.Flags(), baseCfg)
	return cmd
}
