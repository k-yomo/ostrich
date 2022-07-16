package indexer

import (
	"github.com/k-yomo/ostrich/analyzer"
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/internal/opstamp"
	"github.com/k-yomo/ostrich/postings"
	"github.com/k-yomo/ostrich/schema"
)

type SegmentWriter struct {
	maxDoc             index.DocID
	multifieldPostings *postings.PerFieldPostingsWriter
	segmentSerializer  *SegmentSerializer
	docOpstamps        []opstamp.OpStamp
	analyzer           analyzer.Analyzer
}

func newSegmentWriter(memoryBudget int, segment *index.Segment, schema *schema.Schema) (*SegmentWriter, error) {
	segmentSerializer, err := NewSegmentSerializer(segment)
	if err != nil {
		return nil, err
	}
	return &SegmentWriter{
		maxDoc:             0,
		segmentSerializer:  segmentSerializer,
		multifieldPostings: postings.NewMultiFieldPostingsWriter(schema),
		docOpstamps:        nil,
		analyzer:           segment.Index.Analyzer,
	}, nil
}

func (s *SegmentWriter) addDocument(addOperation *AddOperation, sc *schema.Schema) error {
	docID := s.maxDoc
	doc := addOperation.document
	s.docOpstamps = append(s.docOpstamps, addOperation.opstamp)
	for _, fieldAndFieldValues := range doc.SortedFieldValues() {
		fieldEntry := sc.Fields[fieldAndFieldValues.Field]
		switch fieldEntry.FieldType {
		case schema.FieldTypeText:
			var tokens []string
			for _, fieldValue := range fieldAndFieldValues.FieldValues {
				switch v := fieldValue.Value.(type) {
				case string:
					tokens = append(tokens, s.analyzer.Analyze(v)...)
				}
			}

			if len(tokens) == 0 {
				continue
			}
			postingsWriter := s.multifieldPostings.PostingsWriterForFiled(fieldAndFieldValues.Field)
			postingsWriter.IndexText(docID, fieldAndFieldValues.Field, tokens)
		}
	}

	docWriter := s.segmentSerializer.StoreWriter
	if err := docWriter.Store(doc); err != nil {
		return err
	}
	s.maxDoc++

	return nil
}

func (s *SegmentWriter) finalize() error {
	err := s.multifieldPostings.Serialize(s.segmentSerializer.PostingsSerializer)
	if err != nil {
		return err
	}
	return nil
}
