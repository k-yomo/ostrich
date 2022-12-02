package indexer

import (
	"github.com/k-yomo/ostrich/analyzer"
	"github.com/k-yomo/ostrich/collector"
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/internal/opstamp"
	"github.com/k-yomo/ostrich/query"
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
	"reflect"
	"testing"
)

func TestIndexWriter_AddDocument(t *testing.T) {
	s := schema.NewSchema()
	titleField := s.AddTextField("title", analyzer.DefaultAnalyzerName)
	s.AddTextField("description", analyzer.DefaultAnalyzerName)

	idx, err := index.NewBuilder(s).CreateInMemory()
	if err != nil {
		t.Fatalf("failed to create index in memory")
	}
	indexWriter, err := NewIndexWriter(idx, 100000)
	if err != nil {
		t.Fatalf("failed to create index in memory")
	}

	docs := []schema.Document{
		{
			FieldValues: []*schema.FieldValue{{
				FieldID: titleField,
				Value:   "test title",
			}},
		}, {
			FieldValues: []*schema.FieldValue{{
				FieldID: titleField,
				Value:   "abc",
			}},
		},
	}
	for i, doc := range docs {
		addDocResult := indexWriter.AddDocument(&doc)
		if addDocResult.Result() != nil {
			t.Errorf("AddDocument() = %v, want %v", addDocResult.Result(), nil)
		}
		if !reflect.DeepEqual(addDocResult.OpStamp, opstamp.OpStamp(i+1)) {
			t.Errorf("AddDocumentResult.OpStamp = %v, want %v", addDocResult.OpStamp, opstamp.OpStamp(1))
		}
	}
	if _, err := indexWriter.Commit(); err != nil {
		t.Errorf("Commit() = %v, want %v", err, nil)
	}

	indexReader, err := reader.NewIndexReader(idx)
	if err != nil {
		t.Errorf("reader.NewIndexReader() = %v, want %v", err, nil)
	}

	termQuery := query.NewTermQuery(schema.NewTermFromText(titleField, "test"))
	searchResult, err := reader.Search(indexReader.Searcher(), termQuery, collector.NewTopScoreCollector(1, 0))
	if err != nil {
		t.Errorf("reader.Search() = %v, want %v", err, nil)
	}
	if !reflect.DeepEqual(len(searchResult), 1) {
		t.Errorf("hit count = %v, want %v", len(searchResult), 1)
	}
	wantDocAddress := index.DocAddress{SegmentOrd: 0, DocID: 0}
	if !reflect.DeepEqual(searchResult[0].DocAddress, wantDocAddress) {
		t.Errorf("hit's DocAddress = %v, want %v", searchResult[0].DocAddress, wantDocAddress)
	}
}
