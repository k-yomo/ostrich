package index

import (
	"fmt"
	"github.com/k-yomo/ostrich/directory"
	"github.com/k-yomo/ostrich/schema"
)

type SegmentReader struct {
	maxDoc       DocID
	termdictFile directory.ReaderCloser
	storeFile    directory.ReaderCloser
	postingsFile directory.ReaderCloser

	schema *schema.Schema
}

func NewSegmentReader(segment *Segment) (*SegmentReader, error) {
	termdictFile, err := segment.OpenRead(SegmentComponentTerms)
	if err != nil {
		return nil, fmt.Errorf("open termdict file: %w", err)
	}
	storeFile, err := segment.OpenRead(SegmentComponentStore)
	if err != nil {
		return nil, fmt.Errorf("open store file: %w", err)
	}
	postingsFile, err := segment.OpenRead(SegmentComponentPostings)
	if err != nil {
		return nil, fmt.Errorf("open positings file: %w", err)
	}

	return &SegmentReader{
		maxDoc:       segment.Meta().MaxDoc,
		termdictFile: termdictFile,
		storeFile:    storeFile,
		postingsFile: postingsFile,
		schema:       segment.Schema(),
	}, nil
}

func (s *SegmentReader) MaxDoc() DocID {
	return s.maxDoc
}
