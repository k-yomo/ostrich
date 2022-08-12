package main

import (
	"fmt"
	"runtime"

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
	tupleResult, err := reader.Search(
		searcher,
		q,
		tupleCollector,
		reader.SearchOptionConcurrent(runtime.GOMAXPROCS(0)),
	)
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
