package reader

import (
	"encoding/gob"
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

	termdictFile directory.ReaderCloser
	storeFile    directory.ReaderCloser
	postingsFile directory.ReaderCloser

	schema *schema.Schema
}

func NewSegmentReader(segment *index.Segment) (*SegmentReader, error) {
	termdictFile, err := segment.OpenRead(index.SegmentComponentTerms)
	if err != nil {
		return nil, fmt.Errorf("open termdict file: %w", err)
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
		termdictFile: termdictFile,
		storeFile:    storeFile,
		postingsFile: postingsFile,
		schema:       segment.Schema(),
	}, nil
}

func (s *SegmentReader) InvertedIndex(fieldID schema.FieldID) (*postings.PostingsReader, error) {
	// fieldEntry := s.schema.FieldEntry(fieldID)
	termDict := map[schema.FieldID]map[string]*termdict.TermInfo{}
	if err := gob.NewDecoder(s.termdictFile).Decode(&termDict); err != nil {
		return nil, err
	}

	return postings.NewPostingsReader(termDict, s.postingsFile), nil
}
