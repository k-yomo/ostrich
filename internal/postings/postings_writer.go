package postings

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/k-yomo/ostrich/internal/termdict"

	"github.com/k-yomo/ostrich/schema"
)

type UnorderedTermId = uint64

type PostingsWriter interface {
	Serialize(serializer *InvertedIndexSerializer, bytesOffset int) (writtenBytes int, _ error)
	IndexText(docId schema.DocID, field schema.FieldID, tokens []string)
	AddTermFreq(term string, docID schema.DocID, termFreq uint64)
}

type PerFieldPostingsWriter struct {
	perFieldPostingsWriters []PostingsWriter
}

func NewMultiFieldPostingsWriter(s *schema.Schema) *PerFieldPostingsWriter {
	return &PerFieldPostingsWriter{
		perFieldPostingsWriters: newPostingWriters(s.Fields),
	}
}

func (m *PerFieldPostingsWriter) PostingsWriterForFiled(fieldID schema.FieldID) PostingsWriter {
	return m.perFieldPostingsWriters[fieldID]
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
			return fmt.Errorf("serialize per field postings: %w", err)
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
		fieldID:         fieldEntry.ID,
		InvertedIndex:   map[string][]schema.DocID{},
		TermFrequencies: map[string][]uint64{},
	}
}

type SpecializedPostingsWriter struct {
	fieldID         schema.FieldID
	InvertedIndex   map[string][]schema.DocID
	TermFrequencies map[string][]uint64
}

func (s *SpecializedPostingsWriter) Serialize(serializer *InvertedIndexSerializer, bytesOffset int) (writtenBytes int, _ error) {
	var buf []byte
	// TODO: store total number of tokens for calculating BM25
	// buf := make([]byte, 8)
	// totalTokenNum := uint64(len(s.InvertedIndex))
	// binary.LittleEndian.PutUint64(buf, totalTokenNum)
	for term, postingList := range s.InvertedIndex {
		footer := Footer{}
		b := bytes.NewBuffer([]byte{})
		if err := gob.NewEncoder(b).Encode(postingList); err != nil {
			return 0, err
		}
		footer.postingsByteNum = uint64(b.Len())
		if err := gob.NewEncoder(b).Encode(s.TermFrequencies[term]); err != nil {
			return 0, err
		}
		footer.termFrequencyByteNum = uint64(b.Len()) - footer.postingsByteNum
		footer.Write(b)

		from := bytesOffset + len(buf)
		to := from + b.Len()
		serializer.termsWrite.AddTermInfo(s.fieldID, &termdict.TermInfo{
			Term:    term,
			DocFreq: len(postingList),
			PostingsRange: termdict.Range{
				From: from,
				To:   to,
			},
		})
		buf = append(buf, b.Bytes()...)
	}
	writtenBytes, err := serializer.postingsWrite.Write(buf)
	if err != nil {
		return 0, err
	}
	return writtenBytes, nil
}

func (s *SpecializedPostingsWriter) IndexText(docID schema.DocID, field schema.FieldID, terms []string) {
	termFreqMap := map[string]uint64{}
	for _, term := range terms {
		termFreqMap[term] += 1
	}
	for term, freq := range termFreqMap {
		s.AddTermFreq(term, docID, freq)
	}
}

func (s *SpecializedPostingsWriter) AddTermFreq(term string, docID schema.DocID, termFreq uint64) {
	s.InvertedIndex[term] = append(s.InvertedIndex[term], docID)
	s.TermFrequencies[term] = append(s.TermFrequencies[term], termFreq)
}
