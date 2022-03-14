/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package cmd

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/shinexia/elasticdump/pkg/dumpdata"
	"github.com/shinexia/elasticdump/pkg/helpers"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func newCmdDumpData(out io.Writer) *cobra.Command {
	type extraOption struct {
		OutputFile  string
		Batch       int
		SearchQuery string
		SearchBody  string
	}
	cfg := newBaseConfig()
	dopt := dumpdata.NewDumpDataOption()
	extra := &extraOption{
		Batch:       1000,
		SearchQuery: dumpdata.QUERY_ALL,
		SearchBody:  "",
	}
	cmd := &cobra.Command{
		Use:   "data",
		Short: "dump data from elasticsearch",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			klog.V(5).Infof("cfg: %v, opt: %s, extra: %s", helpers.ToJSON(cfg), helpers.ToJSON(dopt), helpers.ToJSON(extra))
			err = preprocessBaseConfig(cfg)
			if err != nil {
				return err
			}
			if extra.OutputFile == "" {
				extra.OutputFile = cfg.Index + "-data.json"
			}
			klog.V(5).Infof("cfg: %v\n", helpers.ToJSON(cfg))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			startTime := time.Now()
			client, err := helpers.NewElasticSearchClient(cfg.Host)
			if err != nil {
				return err
			}
			ops := []func(*esapi.SearchRequest){
				client.Search.WithContext(context.Background()),
				client.Search.WithIndex(cfg.Index),
				client.Search.WithScroll(time.Duration(dopt.TimeoutSec) * time.Second),
				client.Search.WithSize(extra.Batch),
			}
			if extra.SearchBody != "" {
				ops = append(ops, client.Search.WithBody(strings.NewReader(extra.SearchBody)))
			} else {
				ops = append(ops, client.Search.WithQuery(extra.SearchQuery))
			}
			writer := dumpdata.NewLazyDataWriter(extra.OutputFile)
			defer writer.Close()
			ncount, err := dumpdata.DumpData(client, dopt, writer.Write, ops...)
			if err != nil {
				return err
			}
			err = writer.Close()
			if err != nil {
				return err
			}
			cost := time.Since(startTime).Seconds()
			klog.Infof("dump data succeed, total: %d, index: %s, file: %s, cost: %.3fs\n", ncount, cfg.Index, extra.OutputFile, cost)
			return nil
		},
		Args: cobra.NoArgs,
	}

	addBaseConfigFlags(cmd.Flags(), cfg)
	flagSet := cmd.Flags()

	flagSet.IntVarP(&dopt.Limit, "limit", "l", dopt.Limit, "limit size when scroll")
	flagSet.IntVar(&dopt.TimeoutSec, "timeout", dopt.TimeoutSec, "timeout (second) when scroll")

	flagSet.StringVarP(&extra.OutputFile, "file", "f", extra.OutputFile, "output file")

	flagSet.IntVarP(&extra.Batch, "batch", "b", extra.Batch, "batch size when scroll")
	flagSet.StringVarP(&extra.SearchQuery, "search_query", "q", extra.SearchQuery, "search query")
	flagSet.StringVarP(&extra.SearchBody, "search_body", "d", extra.SearchBody, "search body")

	return cmd
}
