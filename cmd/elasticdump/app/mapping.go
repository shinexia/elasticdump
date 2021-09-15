package app

import (
	"io"
	"log"
	"strings"

	"github.com/shinexia/elasticdump/pkg/elasticdump"
	"github.com/spf13/cobra"
)

func newCmdDumpMapping(out io.Writer) *cobra.Command {
	baseCfg := newBaseConfig()
	cmd := &cobra.Command{
		Use:   "mapping",
		Short: "dump mapping from elasticsearch",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			log.Printf("origin: %v\n", elasticdump.ToJSON(baseCfg))
			baseCfg, err = preprocessBaseConfig(baseCfg)
			if err != nil {
				return err
			}
			if baseCfg.File == "" {
				baseCfg.File = baseCfg.Index + "-mapping.json"
			}
			log.Printf("cfg: %v\n", elasticdump.ToJSON(baseCfg))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dumper, err := newDumper(baseCfg)
			if err != nil {
				return err
			}
			return dumper.DumpMapping(baseCfg.Index, baseCfg.File)
		},
		Args: cobra.NoArgs,
	}
	addBaseConfigFlags(cmd.Flags(), baseCfg)
	return cmd
}

func newCmdLoadMapping(out io.Writer) *cobra.Command {
	baseCfg := newBaseConfig()
	type LoadMappingConfig struct {
		delete bool
	}
	cfg := &LoadMappingConfig{
		delete: false,
	}
	cmd := &cobra.Command{
		Use:   "mapping",
		Short: "load mapping to elasticsearch",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			log.Printf("origin: %v\n", elasticdump.ToJSON(baseCfg))
			baseCfg, err = preprocessBaseConfig(baseCfg)
			if err != nil {
				return err
			}
			if baseCfg.File == "" {
				baseCfg.File = baseCfg.Index + "-mapping.json"
			}
			log.Printf("cfg: %v\n", elasticdump.ToJSON(baseCfg))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dumper, err := newDumper(baseCfg)
			if err != nil {
				return err
			}
			if cfg.delete {
				err = dumper.DeleteIndex(baseCfg.Index)
				if err != nil {
					if !strings.Contains(err.Error(), "index_not_found_exception") {
						return err
					}
					log.Printf("index: %s not found\n", baseCfg.Index)
				}
			}
			return dumper.LoadMapping(baseCfg.Index, baseCfg.File)
		},
		Args: cobra.NoArgs,
	}
	addBaseConfigFlags(cmd.Flags(), baseCfg)
	flagSet := cmd.Flags()
	flagSet.BoolVar(&cfg.delete, "delete", cfg.delete, "whether delete the index before load")
	return cmd
}
