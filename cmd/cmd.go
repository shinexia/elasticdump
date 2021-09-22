/*
Copyright 2021 Shine Xia <shine.xgh@gmail.com>.

Licensed under the MIT License.
*/

package cmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
)

// NewKubeadmCommand returns cobra.Command to run kubeadm command
func NewElasticDumpCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "elasticdump",
		Short: "elasticdump: easily dump and load data from elasticsearch",
		Long: dedent.Dedent(`

			    ┌──────────────────────────────────────────────────────────┐
			    │ ElasticDump                                              │
			    │ a simple dump/load data/mapping from elasticsearch       │
			    │                                                          │
			    │ Please give us feedback at:                              │
			    │ https://github.com/shinexia/elasticdump/issues           │
			    └──────────────────────────────────────────────────────────┘

			Example usage:

				elasticdump --host http://localhost:9200 --index elasticdumptest gen  testdata -v=10

				elasticdump --host http://localhost:9200 --index elasticdumptest dump mapping

				elasticdump --host http://localhost:9200 --index elasticdumptest dump data -v=4

				elasticdump --host http://localhost:9200 --index elasticdumptest load mapping --delete

				elasticdump --host http://localhost:9200 --index elasticdumptest load data
		`),
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmds.ResetFlags()

	cmds.AddCommand(newCmdDump(out))
	cmds.AddCommand(newCmdLoad(out))
	cmds.AddCommand(newCmdDelete(out))
	cmds.AddCommand(newCmdGen(out))

	return cmds
}

func newCmdDump(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump",
		Short: "dump mapping/data from elasticsearch",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(newCmdDumpMapping(out))
	cmd.AddCommand(newCmdDumpData(out))
	return cmd
}

func newCmdLoad(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load",
		Short: "load mapping/data from elasticsearch",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(newCmdLoadMapping(out))
	cmd.AddCommand(newCmdLoadData(out))
	return cmd
}

func newCmdDelete(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete index from elasticsearch",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(newCmdDeleteIndex(out))
	return cmd
}

func newCmdGen(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "gen testdata to elasticsearch",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(newCmdGenTestData(out))
	return cmd
}
