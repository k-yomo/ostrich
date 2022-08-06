package postings

import (
	"bytes"
	"encoding/gob"
	"github.com/k-yomo/ostrich/schema"
	"github.com/k-yomo/ostrich/termdict"
)

type UnorderedTermId = uint64

type PostingsWriter interface {
	Serialize(serializer *InvertedIndexSerializer, bytesOffset int) (writtenBytes int, _ error)
	IndexText(docId schema.DocID, field schema.FieldID, tokens []string)
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

func (m *PerFieldPostingsWriter) IndexText(docID schema.DocID, field schema.FieldID, tokens []string) {
	postingsWriter := m.perFieldPostingsWriters[field]
	postingsWriter.IndexText(
		docID,
		field,
		tokens,
	)
}

func (m *PerFieldPostingsWriter) Serialize(serializer *InvertedIndexSerializer) error {
	var offset int
	for _, field := range serializer.schema.Fields {
		postingsWriter := m.PostingsWriterForFiled(field.ID)
		writtenBytes, err := postingsWriter.Serialize(serializer, offset)
		if err != nil {
			return err
		}
		offset += writtenBytes
	}

	if err := serializer.termsWrite.Serialize(); err != nil {
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
	return &SpecializedPostingsWriter{
		fieldID:       fieldEntry.ID,
		InvertedIndex: map[string][]schema.DocID{},
	}
}

type SpecializedPostingsWriter struct {
	fieldID       schema.FieldID
	InvertedIndex map[string][]schema.DocID
}

func (s *SpecializedPostingsWriter) Serialize(serializer *InvertedIndexSerializer, bytesOffset int) (writtenBytes int, _ error) {
	var buf []byte
	for term, postingList := range s.InvertedIndex {
		b := bytes.NewBuffer([]byte{})
		if err := gob.NewEncoder(b).Encode(postingList); err != nil {
			return 0, err
		}
		from := uint64(bytesOffset + len(buf))
		to := from + uint64(b.Len())
		serializer.termsWrite.AddTermInfo(s.fieldID, &termdict.TermInfo{
			Term:    term,
			DocFreq: 1,
			PostingsRange: termdict.Range{
				From: from,
				To:   to,
			},
		})
		buf = append(buf, b.Bytes()...)
	}
	writtenBytes, err := serializer.postingsWrite.Write(buf)
	if err != nil {
		return 0, nil
	}
	return writtenBytes, nil
}

func (s *SpecializedPostingsWriter) IndexText(docId schema.DocID, field schema.FieldID, terms []string) {
	for _, term := range terms {
		s.InvertedIndex[term] = append(s.InvertedIndex[term], docId)
	}
}
