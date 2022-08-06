package termdict

import (
	"bytes"
	"encoding/gob"
	"github.com/k-yomo/ostrich/directory"
	"github.com/k-yomo/ostrich/schema"
)

type TermWriter struct {
	termInfos    map[schema.FieldID]map[string]*TermInfo
	termDictFile directory.WriteCloseSyncer
}

func NewTermWriter(termDictFile directory.WriteCloseSyncer) *TermWriter {
	return &TermWriter{
		termInfos:    map[schema.FieldID]map[string]*TermInfo{},
		termDictFile: termDictFile,
	}
}

func (t *TermWriter) AddTermInfo(fieldID schema.FieldID, termInfo *TermInfo) {
	if _, ok := t.termInfos[fieldID]; !ok {
		t.termInfos[fieldID] = map[string]*TermInfo{}
	}
	t.termInfos[fieldID][termInfo.Term] = termInfo
}

func (t *TermWriter) Serialize() error {
	b := bytes.NewBuffer([]byte{})
	if err := gob.NewEncoder(b).Encode(t.termInfos); err != nil {
		return err
	}
	if _, err := t.termDictFile.Write(b.Bytes()); err != nil {
		return err
	}
	return nil
}

func (t *TermWriter) Close() error {
	return t.termDictFile.Close()
}
