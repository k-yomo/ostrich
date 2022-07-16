package index

import (
	"github.com/k-yomo/ostrich/directory"
	"github.com/k-yomo/ostrich/schema"
)

type Segment struct {
	Index *Index
	meta  *SegmentMeta
}

func newSegment(idx *Index, meta *SegmentMeta) *Segment {
	return &Segment{
		Index: idx,
		meta:  meta,
	}
}

func (s *Segment) Schema() *schema.Schema {
	return s.Index.schema
}

func (s *Segment) Meta() *SegmentMeta {
	return s.meta
}

func (s *Segment) WithMaxDoc(maxDoc uint32) *Segment {
	return newSegment(s.Index, s.meta.WithMaxDoc(maxDoc))
}

func (s *Segment) OpenRead(component SegmentComponent) (directory.ReaderCloser, error) {
	return s.Index.Directory().OpenRead(s.meta.RelativePath(component))
}

func (s *Segment) OpenWrite(component SegmentComponent) (directory.WriteCloseSyncer, error) {
	return s.Index.Directory().OpenWrite(s.meta.RelativePath(component))
}
