package postings

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/k-yomo/ostrich/directory"
	"github.com/k-yomo/ostrich/schema"
	"github.com/k-yomo/ostrich/termdict"
)

type PostingsReader struct {
	termDict     map[schema.FieldID]map[string]*termdict.TermInfo
	postingsFile directory.ReaderCloser
	curDoc       int
}

func NewPostingsReader(
	termDict map[schema.FieldID]map[string]*termdict.TermInfo,
	postingsFile directory.ReaderCloser,
) *PostingsReader {
	return &PostingsReader{
		termDict:     termDict,
		postingsFile: postingsFile,
	}
}

func (p *PostingsReader) ReadPostings(fieldID schema.FieldID, term string) ([]schema.DocID, error) {
	termInfo, ok := p.termDict[fieldID][term]
	if !ok {
		return nil, nil
	}
	postingsBytes := make([]byte, termInfo.PostingsRange.Len())
	if _, err := p.postingsFile.ReadAt(postingsBytes, int64(termInfo.PostingsRange.From)); err != nil {
		return nil, fmt.Errorf("read posting list: %w", err)
	}
	var postingList []schema.DocID
	if err := gob.NewDecoder(bytes.NewReader(postingsBytes)).Decode(&postingList); err != nil {
		return nil, fmt.Errorf("decode posting list: %w", err)
	}
	return postingList, nil
}
