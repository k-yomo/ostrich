package store

import (
	"encoding/json"
	"github.com/k-yomo/ostrich/directory"
	"github.com/k-yomo/ostrich/schema"
	"github.com/k0kubun/pp"
)

const BlockSize = 16_384

type StoreWriter struct {
	doc    int
	writer directory.WriteCloseFlasher
}

func NewStoreWriter(writer directory.WriteCloseFlasher) *StoreWriter {
	return &StoreWriter{
		writer: writer,
	}
}

func (s *StoreWriter) Store(document *schema.Document) error {
	docJSON, err := json.Marshal(document)
	if err != nil {
		return err
	}
	pp.Println(string(docJSON))
	// TODO: makes it possible to read
	if _, err := s.writer.Write(docJSON); err != nil {
		return err
	}
	s.writer.Close()
	s.doc++

	return nil
}
