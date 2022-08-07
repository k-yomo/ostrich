package reader

import (
	"fmt"
	"github.com/k-yomo/ostrich/directory"
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/postings"
	"github.com/k-yomo/ostrich/schema"
	"github.com/k-yomo/ostrich/termdict"
)

type SegmentReader struct {
	SegmentID index.SegmentID
	MaxDoc    schema.DocID

	termDict     termdict.TermDict
	storeFile    *directory.FileSlice
	postingsFile *directory.FileSlice

	schema *schema.Schema
}

func NewSegmentReader(segment *index.Segment) (*SegmentReader, error) {
	termdictFile, err := segment.OpenRead(index.SegmentComponentTerms)
	if err != nil {
		return nil, fmt.Errorf("open termdict file: %w", err)
	}
	termDict, err := termdict.ReadTermDict(termdictFile)
	if err != nil {
		return nil, fmt.Errorf("read termdict: %w", err)
	}
	storeFile, err := segment.OpenRead(index.SegmentComponentStore)
	if err != nil {
		return nil, fmt.Errorf("open store file: %w", err)
	}
	postingsFile, err := segment.OpenRead(index.SegmentComponentPostings)
	if err != nil {
		return nil, fmt.Errorf("open positings file: %w", err)
	}

	return &SegmentReader{
		SegmentID:    segment.Meta().SegmentID,
		MaxDoc:       segment.Meta().MaxDoc,
		termDict:     termDict,
		storeFile:    storeFile,
		postingsFile: postingsFile,
		schema:       segment.Schema(),
	}, nil
}

func (s *SegmentReader) InvertedIndex(fieldID schema.FieldID) (*postings.InvertedIndexReader, error) {
	return postings.NewInvertedIndexReader(s.termDict[fieldID], s.postingsFile), nil
}
