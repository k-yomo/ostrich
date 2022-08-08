package indexer

import (
	"fmt"
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
)

type IndexMerger struct {
	schema         *schema.Schema
	segmentReaders []*reader.SegmentReader
	maxDoc         schema.DocID
}

func NewIndexMerger(schema *schema.Schema, segments []*index.Segment) (*IndexMerger, error) {
	segmentReaders := make([]*reader.SegmentReader, 0, len(segments))
	for _, segment := range segments {
		if segment.Meta().DocNum() > 0 {
			segmentReader, err := reader.NewSegmentReader(segment)
			if err != nil {
				return nil, fmt.Errorf("open segment reader: %w", err)
			}
			segmentReaders = append(segmentReaders, segmentReader)
		}
	}
	return &IndexMerger{
		schema:         schema,
		segmentReaders: segmentReaders,
	}, nil
}

func (i *IndexMerger) Write(serializer *SegmentSerializer) (schema.DocID, error) {
	return i.maxDoc, nil
}
