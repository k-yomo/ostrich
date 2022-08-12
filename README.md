# ostrich

![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)
[![Test](https://github.com/k-yomo/ostrich/actions/workflows/test.yml/badge.svg)](https://github.com/k-yomo/ostrich/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/k-yomo/ostrich)](https://goreportcard.com/report/k-yomo/ostrich)

Full text search engine library written in Go heavily inspired by Tantivy

※ This library is not production ready, don't use it in production.

## example
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
	schemaBuilder := schema.NewBuilder()
	analyzer.Register("en_stem", analyzer.NewEnglishAnalyzer())
	phraseField := schemaBuilder.AddTextField("phrase", "en_stem")
	descriptionField := schemaBuilder.AddTextField("description", "en_stem")

	idx, err := index.NewBuilder(schemaBuilder.Build()).OpenOrCreate("tmp")
	if err != nil {
		panic(err)
	}

	indexWriter, err := indexer.NewIndexWriter(idx, 100_000_000)
	if err != nil {
		panic(err)
	}
	defer indexWriter.Close()

	docs := []*schema.Document{
		{FieldValues: []*schema.FieldValue{
			{
				FieldID: phraseField,
				Value:   "When the Rubber Hits the Road",
			},
			{
				FieldID: descriptionField,
				Value:   "When something is about to begin, get serious, or put to the test.",
			},
		}},
	}

	for _, doc := range docs {
		indexWriter.AddDocument(doc)
	}
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
		collector.NewTopDocsCollector(10, 0),
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
