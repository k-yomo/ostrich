package cmd

import (
	"fmt"
	"io"
	"time"

	"github.com/k-yomo/ostrich/collector"
	"github.com/k-yomo/ostrich/directory"
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/query"
	"github.com/k-yomo/ostrich/reader"
	"github.com/spf13/cobra"
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
				return fmt.Errorf("get limit flag: %w", err)
			}
			dir, err := directory.NewMMapDirectory(indexPath)
			if err != nil {
				return fmt.Errorf("initialize mmap directory: %w", err)
			}
			idx, err := index.OpenIndex(dir)
			if err != nil {
				return fmt.Errorf("open index: %w", err)
			}
			indexReader, err := reader.NewIndexReader(idx)
			if err != nil {
				return fmt.Errorf("open index reader: %w", err)
			}
			searcher := indexReader.Searcher()

			start := time.Now()
			queryParser := query.NewParser(idx.Schema(), idx.Schema().FieldIDs())
			q, err := queryParser.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse query: %w", err)
			}
			hits, err := reader.Search(searcher, q, collector.NewTopDocsCollector(limit, 0))
			if err != nil {
				return fmt.Errorf("search: %w", err)
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
