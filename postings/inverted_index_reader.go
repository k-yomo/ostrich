package postings

import (
	"fmt"
	"github.com/k-yomo/ostrich/directory"
	"github.com/k-yomo/ostrich/termdict"
)

type InvertedIndexReader struct {
	termDict     map[string]*termdict.TermInfo
	postingsFile directory.ReaderCloser
}

func NewInvertedIndexReader(
	termDict map[string]*termdict.TermInfo,
	postingsFile directory.ReaderCloser,
) *InvertedIndexReader {
	return &InvertedIndexReader{
		termDict:     termDict,
		postingsFile: postingsFile,
	}
}

func (p *InvertedIndexReader) ReadPostings(term string) (*PostingsReader, error) {
	termInfo, ok := p.termDict[term]
	if !ok {
		return nil, nil
	}
	postingsBytes := make([]byte, termInfo.PostingsRange.Len())
	if _, err := p.postingsFile.ReadAt(postingsBytes, int64(termInfo.PostingsRange.From)); err != nil {
		return nil, fmt.Errorf("read posting list: %w", err)
	}
	return NewPostingsReader(postingsBytes)
}

func (p *InvertedIndexReader) DocFreq(term string) int {
	termInfo, ok := p.termDict[term]
	if !ok {
		return 0
	}
	return termInfo.DocFreq
}
