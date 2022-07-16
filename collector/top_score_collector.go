package collector

import "github.com/k-yomo/ostrich/index"

type TopDocsCollector struct {
	limit  int
	offset int
}

type TopDocsResult struct {
	DocID uint32
	Score float64
}

func NewTopDocsCollector(limit int, offset int) Collector[*TopDocsResult] {
	return &TopDocsCollector{
		limit:  limit,
		offset: offset,
	}
}

func (t *TopDocsCollector) CollectSegment(segmentReader *index.SegmentReader) []*TopDocsResult {
	// TODO implement me
	panic("implement me")
}

func (t *TopDocsCollector) MergeResults(results [][]*TopDocsResult) []*TopDocsResult {
	// TODO implement me
	panic("implement me")
}
