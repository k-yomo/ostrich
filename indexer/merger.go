package indexer

import (
	"fmt"
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/postings"
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
)

type IndexMerger struct {
	schema         *schema.Schema
	segmentReaders []*reader.SegmentReader
	maxDoc         schema.DocID
}

func NewIndexMerger(indexSchema *schema.Schema, segments []*index.Segment) (*IndexMerger, error) {
	segmentReaders := make([]*reader.SegmentReader, 0, len(segments))
	var maxDoc uint32
	for _, segment := range segments {
		if segment.Meta().DocNum() > 0 {
			segmentReader, err := reader.NewSegmentReader(segment)
			if err != nil {
				return nil, fmt.Errorf("open segment reader: %w", err)
			}
			maxDoc += segmentReader.DocNum()
			segmentReaders = append(segmentReaders, segmentReader)
		}
	}
	return &IndexMerger{
		schema:         indexSchema,
		segmentReaders: segmentReaders,
		maxDoc:         schema.DocID(maxDoc),
	}, nil
}

func (i *IndexMerger) Write(serializer *SegmentSerializer) (schema.DocID, error) {
	newToOldDocIDMappping := i.GenerateNewToDocIDMapping()

	oldToNewDocIDMap := make([]map[schema.DocID]schema.DocID, len(i.segmentReaders))
	for newDocID, oldDocAddress := range newToOldDocIDMappping {
		newDocID := schema.DocID(newDocID)
		if m := oldToNewDocIDMap[oldDocAddress.SegmentOrd]; m == nil {
			oldToNewDocIDMap[oldDocAddress.SegmentOrd] = map[schema.DocID]schema.DocID{}
		}
		oldToNewDocIDMap[oldDocAddress.SegmentOrd][oldDocAddress.DocID] = newDocID
	}

	perFieldPostingsWriter := postings.NewMultiFieldPostingsWriter(i.schema)
	for _, fieldEntry := range i.schema.Fields {
		fieldTermMap := map[string]struct{}{}
		for _, segmentReader := range i.segmentReaders {
			fieldReader := segmentReader.InvertedIndex(fieldEntry.ID)
			for term := range fieldReader.TermDict() {
				fieldTermMap[term] = struct{}{}
			}
		}
		postingsWriter := perFieldPostingsWriter.PostingsWriterForFiled(fieldEntry.ID)
		for term := range fieldTermMap {
			termDocCount := 0
			for segmentOrd, segmentReader := range i.segmentReaders {
				invertedIndexReader := segmentReader.InvertedIndex(fieldEntry.ID)
				postingsReader, err := invertedIndexReader.ReadPostings(term)
				if err != nil {
					return 0, err
				}
				docID := postingsReader.Doc()
				for docID != schema.DocIDTerminated {
					newDocID := oldToNewDocIDMap[segmentOrd][docID]
					postingsWriter.AddTermFreq(term, newDocID, postingsReader.TermFreq())
					termDocCount += 1
					docID = postingsReader.Advance()
				}
			}
		}
	}

	if err := perFieldPostingsWriter.Serialize(serializer.PostingsSerializer); err != nil {
		return 0, fmt.Errorf("serialize postings: %w", err)
	}
	if err := serializer.Close(); err != nil {
		return 0, fmt.Errorf("close serializer: %w", err)
	}

	// TODO: merge documents
	// serializer.StoreWriter

	return i.maxDoc, nil
}

func (i *IndexMerger) GenerateNewToDocIDMapping() []index.DocAddress {
	var docIDMapping []index.DocAddress
	for segmentOrd, segmentReader := range i.segmentReaders {
		segmentDocNum := int(segmentReader.DocNum())
		for i := 0; i < segmentDocNum; i++ {
			docIDMapping = append(docIDMapping, index.DocAddress{
				SegmentOrd: segmentOrd,
				DocID:      schema.DocID(i),
			})
		}
	}
	return docIDMapping
}
