package postings

import (
	"github.com/k-yomo/ostrich/directory"
	"github.com/k-yomo/ostrich/termdict"
)

type InvertedIndexReader struct {
	termDict     map[string]*termdict.TermInfo
	postingsFile *directory.FileSlice
}

func NewInvertedIndexReader(
	termDict map[string]*termdict.TermInfo,
	postingsFile *directory.FileSlice,
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
	return NewPostingsReader(p.postingsFile.Slice(termInfo.PostingsRange.From, termInfo.PostingsRange.To))
}

func (p *InvertedIndexReader) DocFreq(term string) int {
	termInfo, ok := p.termDict[term]
	if !ok {
		return 0
	}
	return termInfo.DocFreq
}
