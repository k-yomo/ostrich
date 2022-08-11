package store

import (
	"encoding/json"

	"github.com/k-yomo/ostrich/directory"
	"github.com/k-yomo/ostrich/schema"
)

const BlockSize = 16_384

type StoreWriter struct {
	doc    int
	writer directory.WriteCloseSyncer
}

func NewStoreWriter(writer directory.WriteCloseSyncer) *StoreWriter {
	return &StoreWriter{
		writer: writer,
	}
}

func (s *StoreWriter) Store(document *schema.Document) error {
	docJSON, err := json.Marshal(document)
	if err != nil {
		return err
	}
	// TODO: makes it possible to read
	if _, err := s.writer.Write(docJSON); err != nil {
		return err
	}
	s.doc++

	return nil
}

func (s *StoreWriter) Close() error {
	return s.writer.Close()
}
