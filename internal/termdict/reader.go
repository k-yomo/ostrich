package termdict

import (
	"encoding/gob"
	"fmt"

	"github.com/k-yomo/ostrich/directory"
	"github.com/k-yomo/ostrich/schema"
)

func ReadPerFieldTermDict(termdictFile *directory.FileSlice) (map[schema.FieldID]TermDict, error) {
	termDict := map[schema.FieldID]TermDict{}
	if err := gob.NewDecoder(termdictFile.Reader()).Decode(&termDict); err != nil {
		return nil, fmt.Errorf("decode termdict file: %w", err)
	}
	return termDict, nil
}
