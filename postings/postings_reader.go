package postings

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/k-yomo/ostrich/directory"
	"github.com/k-yomo/ostrich/schema"
	"sort"
)

type PostingsReader struct {
	postingList     []schema.DocID
	termFrequencies []uint64

	curIdx int
}

func NewPostingsReader(postingsFile *directory.FileSlice) (*PostingsReader, error) {
	postingsBytes, err := postingsFile.Read(0, postingsFile.Len())
	if err != nil {
		return nil, fmt.Errorf("read posting file: %w", err)
	}

	p := &PostingsReader{}
	footer := readFooter(postingsBytes)
	if err := gob.NewDecoder(bytes.NewReader(postingsBytes[:footer.postingsByteNum])).Decode(&p.postingList); err != nil {
		return nil, fmt.Errorf("decode posting list: %w", err)
	}
	if err := gob.NewDecoder(bytes.NewReader(postingsBytes[footer.postingsByteNum:])).Decode(&p.termFrequencies); err != nil {
		return nil, fmt.Errorf("decode term frequencies: %w", err)
	}

	return p, nil
}

func (p *PostingsReader) Advance() schema.DocID {
	if p.curIdx < len(p.postingList) {
		p.curIdx += 1
	}
	return p.Doc()
}

func (p *PostingsReader) Doc() schema.DocID {
	if p.curIdx >= len(p.postingList) {
		return schema.DocIDTerminated
	}
	return p.postingList[p.curIdx]
}

func (p *PostingsReader) Seek(target schema.DocID) schema.DocID {
	if p.Doc() >= target {
		return p.Doc()
	}
	nextIndex := sort.Search(len(p.postingList[p.curIdx:]), func(i int) bool {
		return p.postingList[i] >= target
	})
	p.curIdx = p.curIdx + nextIndex
	return p.Doc()
}

func (p *PostingsReader) SizeHint() uint32 {
	return uint32(len(p.postingList))
}

func (p *PostingsReader) TermFreq() uint64 {
	return p.termFrequencies[p.curIdx]
}
