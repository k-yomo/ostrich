package cmd

import (
	"fmt"
	"github.com/k-yomo/ostrich/collector"
	"github.com/k-yomo/ostrich/directory"
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/query"
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
	"github.com/spf13/cobra"
	"io"
	"time"
)

// newSearchCmd returns the command to search index
func newSearchCmd(out io.Writer) *cobra.Command {
	const flagNameLimit = "limit"

	command := &cobra.Command{
		Use:     "search QUERY",
		Short:   "search Ostrich index",
		Long:    "search Ostrich index with given query",
		Example: "search 'query' --path=/path/to/index",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			limit, err := cmd.Flags().GetInt(flagNameLimit)
			if err != nil {
				return err
			}
			dir, err := directory.NewMMapDirectory(indexPath)
			if err != nil {
				return err
			}
			idx, err := index.OpenIndex(dir)
			if err != nil {
				return err
			}
			indexReader, err := reader.NewIndexReader(idx)
			if err != nil {
				return err
			}
			searcher := indexReader.Searcher()

			start := time.Now()
			// TODO: fix to use query parser
			term := schema.NewTermFromText(idx.Schema().Fields[0].ID, args[0])
			termQuery := query.NewTermQuery(term)
			hits, err := reader.Search(searcher, termQuery, collector.NewTopDocsCollector(limit, 0))
			if err != nil {
				panic(err)
			}
			elapsedTime := time.Since(start).String()
			for _, hit := range hits {
				out.Write([]byte(fmt.Sprintf("docAddress: %+v, score: %.6f\n", hit.DocAddress, hit.Score)))
			}
			out.Write([]byte(fmt.Sprintf("elapsed time: %s\n", elapsedTime)))
			return nil
		},
	}
	command.PersistentFlags().IntP(flagNameLimit, "l", 10, "how many documents to collect")
	command.SetOut(out)
	return command
}
