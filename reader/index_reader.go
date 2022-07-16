package reader

import (
	"fmt"
	"github.com/k-yomo/ostrich/index"
)

type IndexReader struct {
	index    *index.Index
	searcher *index.Searcher
}

func NewIndexReader(idx *index.Index) (*IndexReader, error) {
	indexReader := &IndexReader{
		index: idx,
	}
	if err := indexReader.Reload(); err != nil {
		return nil, fmt.Errorf("reload searcher: %w", err)
	}
	return indexReader, nil
}

func (i *IndexReader) Reload() error {
	segmentReaders, err := i.openSegmentReaders()
	if err != nil {
		return fmt.Errorf("open segment readers: %w", err)
	}
	// TODO: update to have searcher pool
	searcher := index.NewSearcher(i.index.Schema(), i.index, segmentReaders)
	i.searcher = searcher
	return nil
}

func (i *IndexReader) Searcher() *index.Searcher {
	return i.searcher
}

func (i *IndexReader) openSegmentReaders() ([]*index.SegmentReader, error) {
	searchableSegments, err := i.index.SearchableSegments()
	if err != nil {
		return nil, fmt.Errorf("get searchable segments: %w", err)
	}

	segmentReaders := make([]*index.SegmentReader, 0, len(searchableSegments))
	for _, segment := range searchableSegments {
		segmentReader, err := index.NewSegmentReader(segment)
		if err != nil {
			return nil, fmt.Errorf("initialize segment reader: %w", err)
		}
		segmentReaders = append(segmentReaders, segmentReader)
	}

	return segmentReaders, nil
}
