package app

import (
	"io"
	"log"

	"github.com/shinexia/elasticdump/pkg/elasticdump"
	"github.com/spf13/cobra"
)

func newCmdGenTestData(out io.Writer) *cobra.Command {
	baseCfg := newBaseConfig()
	type GenTestDataConfig struct {
		epoch int
		batch int
	}
	cfg := &GenTestDataConfig{
		epoch: 10,
		batch: 100,
	}
	cmd := &cobra.Command{
		Use:   "testdata",
		Short: "gen testdata to elasticsearch for test",
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
			return dumper.GenTestData(baseCfg.Index, cfg.epoch, cfg.batch)
		},
		Args: cobra.NoArgs,
	}
	addBaseConfigFlags(cmd.Flags(), baseCfg)
	return cmd
}
