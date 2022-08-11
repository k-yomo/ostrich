package index

import (
	"encoding/json"
	"fmt"
	"sync"

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
	SegmentID SegmentID    `json:"segmentId"`
	MaxDoc    schema.DocID `json:"maxDoc"`
	Deletes   *DeleteMeta  `json:"deletes"`
}

type SegmentMetaInventory struct {
	inventory []*SegmentMeta
	mu        sync.Mutex
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

func (s *SegmentMeta) WithMaxDoc(maxDoc schema.DocID) *SegmentMeta {
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

func (i *SegmentMetaInventory) NewSegmentMeta(segmentID SegmentID, maxDoc schema.DocID) *SegmentMeta {
	segmentMeta := &SegmentMeta{
		SegmentID: segmentID,
		MaxDoc:    maxDoc,
		Deletes:   nil,
	}
	i.mu.Lock()
	i.inventory = append(i.inventory, segmentMeta)
	i.mu.Unlock()
	return segmentMeta
}

func SaveMeta(meta *IndexMeta, dir directory.Directory) error {
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("marshal index meta: %w", err)
	}
	return dir.AtomicWrite(metaFileName, metaJSON)
}

func LoadMeta(directory directory.Directory, inventory *SegmentMetaInventory) (*IndexMeta, error) {
	metaData, err := directory.AtomicRead(metaFileName)
	if err != nil {
		return nil, err
	}
	var indexMeta IndexMeta
	if err := json.Unmarshal(metaData, &indexMeta); err != nil {
		return nil, err
	}

	return &indexMeta, nil
}
