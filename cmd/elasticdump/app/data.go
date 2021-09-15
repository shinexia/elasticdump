package app

import (
	"io"
	"log"
	"strings"
	"time"

	"github.com/shinexia/elasticdump/pkg/elasticdump"
	"github.com/spf13/cobra"
)

func newCmdDumpData(out io.Writer) *cobra.Command {
	baseCfg := newBaseConfig()
	type DumpDataConfig struct {
		batch      int
		limit      int
		timeoutSec int
	}
	cfg := &DumpDataConfig{
		batch:      1000,
		limit:      -1,
		timeoutSec: 60,
	}
	cmd := &cobra.Command{
		Use:   "data",
		Short: "dump data from elasticsearch",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			log.Printf("origin: %v\n", elasticdump.ToJSON(baseCfg))
			baseCfg, err = preprocessBaseConfig(baseCfg)
			if err != nil {
				return err
			}
			if baseCfg.File == "" {
				baseCfg.File = baseCfg.Index + "-data.json"
			}
			log.Printf("cfg: %v\n", elasticdump.ToJSON(baseCfg))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dumper, err := newDumper(baseCfg)
			if err != nil {
				return err
			}
			return dumper.DumpData(baseCfg.Index, baseCfg.File, cfg.batch, cfg.limit, time.Duration(cfg.timeoutSec)*time.Second)
		},
		Args: cobra.NoArgs,
	}

	addBaseConfigFlags(cmd.Flags(), baseCfg)
	flagSet := cmd.Flags()
	flagSet.IntVarP(&cfg.batch, "batch", "b", cfg.batch, "batch size when scroll")
	flagSet.IntVarP(&cfg.batch, "limit", "l", cfg.limit, "limit size when scroll")
	flagSet.IntVar(&cfg.batch, "timeout", cfg.timeoutSec, "timeout (second) when scroll")
	return cmd
}

func newCmdLoadData(out io.Writer) *cobra.Command {
	baseCfg := newBaseConfig()
	type LoadDataConfig struct {
		batch   int
		limit   int
		bufSize int
		delete  bool
	}
	cfg := &LoadDataConfig{
		batch:   1000,
		limit:   -1,
		bufSize: 1024 * 1024 * 1024,
		delete:  false,
	}
	cmd := &cobra.Command{
		Use:   "data",
		Short: "load data to elasticsearch",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			log.Printf("origin: %v\n", elasticdump.ToJSON(baseCfg))
			baseCfg, err = preprocessBaseConfig(baseCfg)
			if err != nil {
				return err
			}
			if baseCfg.File == "" {
				baseCfg.File = baseCfg.Index + "-data.json"
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
			return dumper.LoadData(baseCfg.Index, baseCfg.File, cfg.batch, cfg.limit, cfg.bufSize)
		},
		Args: cobra.NoArgs,
	}
	addBaseConfigFlags(cmd.Flags(), baseCfg)
	flagSet := cmd.Flags()
	flagSet.IntVarP(&cfg.batch, "batch", "b", cfg.batch, "batch size when scroll")
	flagSet.IntVarP(&cfg.batch, "limit", "l", cfg.limit, "limit size when scroll")
	flagSet.IntVar(&cfg.bufSize, "buf", cfg.bufSize, "buffer size (byte) when split data file to lines, must bigger than the largest line")
	flagSet.BoolVar(&cfg.delete, "delete", cfg.delete, "whether delete the index before load")
	return cmd
}
