package postings

import (
	"bytes"
	"encoding/gob"
	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/schema"
)

type UnorderedTermId = uint64

type PostingsWriter interface {
	Serialize(serializer *InvertedIndexSerializer, bytesOffset int) (writtenBytes int, termOffsetMap map[string]uint64, _ error)
	IndexText(docId index.DocID, field schema.FieldID, tokens []string)
}

type PerFieldPostingsWriter struct {
	perFieldPostingsWriters []PostingsWriter
}

func NewMultiFieldPostingsWriter(s *schema.Schema) *PerFieldPostingsWriter {
	return &PerFieldPostingsWriter{
		perFieldPostingsWriters: newPostingWriters(s.Fields),
	}
}

func (m *PerFieldPostingsWriter) PostingsWriterForFiled(field schema.FieldID) PostingsWriter {
	return m.perFieldPostingsWriters[field]
}

func (m *PerFieldPostingsWriter) IndexText(docID index.DocID, field schema.FieldID, tokens []string) {
	postingsWriter := m.perFieldPostingsWriters[field]
	postingsWriter.IndexText(
		docID,
		field,
		tokens,
	)
}

func (m *PerFieldPostingsWriter) Serialize(serializer *InvertedIndexSerializer) error {
	fieldTermPostingsOffsetMap := map[schema.FieldID]map[string]uint64{}
	var offset int
	for _, field := range serializer.schema.Fields {
		postingsWriter := m.PostingsWriterForFiled(field.ID)
		writtenBytes, termOffsetMap, err := postingsWriter.Serialize(serializer, offset)
		if err != nil {
			return err
		}
		fieldTermPostingsOffsetMap[field.ID] = termOffsetMap
		offset += writtenBytes
	}

	b := bytes.NewBuffer([]byte{})
	if err := gob.NewEncoder(b).Encode(fieldTermPostingsOffsetMap); err != nil {
		return err
	}
	if _, err := serializer.termsWrite.Write(b.Bytes()); err != nil {
		return err
	}

	return nil
}

func newPostingWriters(fields []*schema.FieldEntry) []PostingsWriter {
	var postingsWriter []PostingsWriter
	for _, field := range fields {
		postingsWriter = append(postingsWriter, newPostingWriterFromFieldEntry(field))
	}
	return postingsWriter
}

func newPostingWriterFromFieldEntry(fieldEntry *schema.FieldEntry) PostingsWriter {
	// // TODO: support other types
	// switch fieldEntry.FieldType {
	// case schema.FieldTypeText:
	// }
	return &SpecializedPostingsWriter{InvertedIndex: map[string][]index.DocID{}}
}

type SpecializedPostingsWriter struct {
	InvertedIndex map[string][]index.DocID
}

func (s *SpecializedPostingsWriter) Serialize(serializer *InvertedIndexSerializer, bytesOffset int) (writtenBytes int, termOffsetMap map[string]uint64, _ error) {
	termOffsetMap = map[string]uint64{}
	var buf []byte
	for term, postings := range s.InvertedIndex {
		b := bytes.NewBuffer([]byte{})
		if err := gob.NewEncoder(b).Encode(postings); err != nil {
			return 0, nil, err
		}
		termOffsetMap[term] = uint64(bytesOffset + len(buf) + 1)
		buf = append(buf, b.Bytes()...)
	}
	writtenBytes, err := serializer.postingsWrite.Write(buf)
	if err != nil {
		return 0, nil, nil
	}
	return writtenBytes, termOffsetMap, nil
}

func (s *SpecializedPostingsWriter) IndexText(docId index.DocID, field schema.FieldID, terms []string) {
	for _, term := range terms {
		s.InvertedIndex[term] = append(s.InvertedIndex[term], docId)
	}
}
