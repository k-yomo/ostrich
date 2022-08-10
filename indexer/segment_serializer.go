package indexer

import (
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/postings"
	"github.com/k-yomo/ostrich/store"
)

type SegmentSerializer struct {
	StoreWriter        *store.StoreWriter
	PostingsSerializer *postings.InvertedIndexSerializer
}

func NewSegmentSerializer(segment *index.Segment) (*SegmentSerializer, error) {
	storeWrite, err := segment.OpenWrite(index.SegmentComponentStore)
	if err != nil {
		return nil, err
	}

	postingsSerializer, err := postings.NewInvertedIndexSerializer(segment)
	if err != nil {
		return nil, err
	}

	return &SegmentSerializer{
		StoreWriter:        store.NewStoreWriter(storeWrite),
		PostingsSerializer: postingsSerializer,
	}, nil
}

func (s *SegmentSerializer) Close() error {
	if err := s.PostingsSerializer.Close(); err != nil {
		return err
	}
	if err := s.StoreWriter.Close(); err != nil {
		return err
	}
	return nil
}
