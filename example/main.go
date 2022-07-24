package main

import (
	"fmt"
	"github.com/k-yomo/ostrich/collector"
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/indexer"
	"github.com/k-yomo/ostrich/query"
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
	"time"
)

func main() {
	schemaBuilder := schema.NewBuilder()
	phraseField := schemaBuilder.AddTextField("phrase")
	descriptionField := schemaBuilder.AddTextField("description")
	//
	idx, err := index.NewBuilder(schemaBuilder.Build()).CreateInDir("tmp")
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
				Value:   "Down To The Wire",
			},
			{
				FieldID: descriptionField,
				Value:   "A tense situation where the outcome is decided only in the last few seconds.",
			},
		}},
		{FieldValues: []*schema.FieldValue{
			{
				FieldID: phraseField,
				Value:   "Eat My Hat",
			},
			{
				FieldID: descriptionField,
				Value:   "Having confidence in a specific outcome; being almost sure about something",
			},
		}},
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
	time.Sleep(5 * time.Second)
	// idx.DeleteDoc(3)

	indexReader, err := reader.NewIndexReader(idx)
	if err != nil {
		panic(err)
	}
	defer indexReader.Close()

	searcher := indexReader.Searcher()

	queryParser := query.NewParser(idx, []schema.FieldID{phraseField, descriptionField})
	q := queryParser.Parse("black cat")
	hits := index.Search(searcher, q, collector.NewTopDocsCollector(10, 0))
	for _, hit := range hits {
		fmt.Printf("docAddress: %+v, score: %v\n", hit.DocAddress, hit.Score)
	}
}
