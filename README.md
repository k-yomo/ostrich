# ostrich

![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)
[![Test](https://github.com/k-yomo/ostrich/actions/workflows/test.yml/badge.svg)](https://github.com/k-yomo/ostrich/actions/workflows/test.yml)
[![Codecov](https://codecov.io/gh/k-yomo/ostrich/branch/main/graph/badge.svg?token=P3pNbMGbeN)](https://codecov.io/gh/k-yomo/ostrich)
[![Go Report Card](https://goreportcard.com/badge/k-yomo/ostrich)](https://goreportcard.com/report/k-yomo/ostrich)

Full text search engine library written in Go with 1.18+ Generics, heavily inspired by Tantivy

※ This library is not production ready, don't use it in production.

## Features
- Full-text search
- Configurable analyzer
- Concurrent indexing in batch
- Segment merge with LogMergePolicy
- Mmap directory
- Natural query language (e.g. "(go OR golang) AND (search or fts)")
- Concurrent search
- TF-IDF scoring (will be replaced with BM25)

#### Supported field types:
  - Text
#### Supported query types:
  - Term, Conjunction, Disjunction, Boolean

※ We'll support more and more types

## Example

```go
package main

import (
	"fmt"

	"github.com/k-yomo/ostrich/analyzer"
	"github.com/k-yomo/ostrich/collector"
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/indexer"
	"github.com/k-yomo/ostrich/query"
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
)

func main() {
	indexSchema := schema.NewSchema()
	analyzer.Register("en_stem", analyzer.NewEnglishAnalyzer())
	phraseField := indexSchema.AddTextField("phrase", "en_stem")
	descriptionField := indexSchema.AddTextField("description", "en_stem")

	idx, err := index.NewBuilder(indexSchema).OpenOrCreate("tmp")
	if err != nil {
		panic(err)
	}

	indexWriter, err := indexer.NewIndexWriter(idx, 100_000_000)
	if err != nil {
		panic(err)
	}
	defer indexWriter.Close()

	doc := &schema.Document{
		FieldValues: []*schema.FieldValue{
			{
				FieldID: phraseField,
				Value:   "When the Rubber Hits the Road",
			},
			{
				FieldID: descriptionField,
				Value:   "When something is about to begin, get serious, or put to the test.",
			},
		},
	}
	indexWriter.AddDocument(doc)
	if _, err := indexWriter.Commit(); err != nil {
		panic(err)
	}

	indexReader, err := reader.NewIndexReader(idx)
	if err != nil {
		panic(err)
	}
	defer indexReader.Close()

	queryParser := query.NewParser(idx.Schema(), idx.Schema().FieldIDs())
	q, err := queryParser.Parse("phrase:hat OR description:serious")
	if err != nil {
		panic(err)
	}
	tupleCollector := collector.NewTupleCollector(
		collector.NewTopScoreCollector(10, 0),
		collector.NewCountCollector(),
	)

	searcher := indexReader.Searcher()
	tupleResult, err := reader.Search(searcher, q, tupleCollector)
	if err != nil {
		panic(err)
	}

	hits := tupleResult.Left
	count := tupleResult.Right
	fmt.Println("total hit:", count)
	for _, hit := range hits {
		fmt.Printf("docAddress: %+v, score: %v\n", hit.DocAddress, hit.Score)
	}
}
```
