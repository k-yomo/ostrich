package index

import (
	"encoding/json"
	"fmt"
	"github.com/k-yomo/ostrich/directory"
	"github.com/k-yomo/ostrich/internal/opstamp"
	"github.com/k-yomo/ostrich/schema"
)

const metaFileName = "meta.json"

type IndexMeta struct {
	Segments []*SegmentMeta `json:"segments"`
	Schema   *schema.Schema `json:"schema"`
	// last commit operation's id
	Opstamp opstamp.OpStamp `json:"opstamp"`
}

type SegmentMeta struct {
	SegmentID SegmentID   `json:"segmentId"`
	MaxDoc    DocID       `json:"maxDoc"`
	Deletes   *DeleteMeta `json:"deletes"`
}

type SegmentMetaInventory struct {
	inventory []*SegmentMeta
}

type DeleteMeta struct {
	NumDeletedDocs uint32          `json:"numDeletedDocs"`
	Opstamp        opstamp.OpStamp `json:"operationId"`
}

func NewIndexMeta(schema *schema.Schema) *IndexMeta {
	return &IndexMeta{
		Segments: nil,
		Schema:   schema,
		Opstamp:  0,
	}
}

func (s *SegmentMeta) WithMaxDoc(maxDoc DocID) *SegmentMeta {
	return &SegmentMeta{
		SegmentID: s.SegmentID,
		MaxDoc:    maxDoc,
		Deletes:   nil,
	}
}

func (s *SegmentMeta) DocNum() uint32 {
	// Currently we don't support delete
	// return s.MaxDoc - s.Deletes.NumDeletedDocs
	return uint32(s.MaxDoc)
}

func (s *SegmentMeta) RelativePath(component SegmentComponent) string {
	path := s.SegmentID
	switch component {
	case SegmentComponentPostings:
		return fmt.Sprintf("%s.idx", path)
	case SegmentComponentTerms:
		return fmt.Sprintf("%s.term", path)
	case SegmentComponentStore:
		return fmt.Sprintf("%s.store", path)
	case SegmentComponentDelete:
		return fmt.Sprintf("%s.%d.del", path, s.Deletes.Opstamp)
	default:
		panic(fmt.Sprintf("invalid component: %v", component))
	}
}

func (i *SegmentMetaInventory) NewSegmentMeta(segmentID SegmentID, maxDoc DocID) *SegmentMeta {
	segmentMeta := &SegmentMeta{
		SegmentID: segmentID,
		MaxDoc:    maxDoc,
		Deletes:   nil,
	}
	// TODO: Make it thread safe
	i.inventory = append(i.inventory, segmentMeta)
	return segmentMeta
}

func SaveMetas(meta *IndexMeta, dir directory.Directory) error {
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	return dir.AtomicWrite(metaFileName, metaJSON)
}
