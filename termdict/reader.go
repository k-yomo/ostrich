package termdict

import (
	"encoding/gob"
	"fmt"
	"github.com/k-yomo/ostrich/directory"
)

func ReadTermDict(termdictFile *directory.FileSlice) (TermDict, error) {
	termDict := TermDict{}
	if err := gob.NewDecoder(termdictFile.Reader()).Decode(&termDict); err != nil {
		return nil, fmt.Errorf("decode termdict file: %w", err)
	}
	return termDict, nil
}
