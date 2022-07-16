package main

import (
	"fmt"
	"github.com/k-yomo/ostrich/collector"
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/indexer"
	"github.com/k-yomo/ostrich/query"
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
	"math/rand"
	"time"
)

func main() {
	schemaBuilder := schema.NewBuilder()
	titleField := schemaBuilder.AddTextField("title")
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

	docs := []*schema.Document{
		{FieldValues: []*schema.FieldValue{
			{
				FieldID: titleField,
				Value:   "there is a white cat",
			},
			{
				FieldID: descriptionField,
				Value:   "this is a description",
			},
		}},
		// {ID: 2, Text: "black hair cat"},
		// {ID: 3, Text: "black cat"},
		// {ID: 4, Text: "white dog"},
		// {ID: 5, Text: "blue cat"},
		// {ID: 6, Text: "black tiger"},
		// {ID: 7, Text: "white hair dog"},
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(docs), func(i, j int) { docs[i], docs[j] = docs[j], docs[i] })
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
	searcher := indexReader.Searcher()

	queryParser := query.NewParser(idx, []schema.FieldID{titleField, descriptionField})
	q := queryParser.Parse("black cat")
	hits := index.Search(searcher, q, collector.NewTopDocsCollector(10, 0))
	for _, hit := range hits {
		fmt.Printf("docID: %v, score: %v\n", hit.DocID, hit.Score)
	}
}
