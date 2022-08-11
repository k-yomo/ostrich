package indexer

import (
	"math"
	"sort"

	"github.com/k-yomo/ostrich/index"
)

const DefaultLevelLogSize = 0.75
const DefaultMinLayerSize = 10_000
const DefaultMinNumSegmentsToMerge = 8
const DefaultMaxDocsBeforeMerge = 10_000_000
const DefaultDelDocsRatioBeforeMerge = 1.0

type LogMergePolicy struct {
	minNumSegments          int
	maxDocsBeforeMerge      uint32
	minLayerSize            uint32
	levelLogSize            float64
	delDocsRatioBeforeMerge float32
}

func NewLogMergePolicy() *LogMergePolicy {
	return &LogMergePolicy{
		minNumSegments:          DefaultMinNumSegmentsToMerge,
		maxDocsBeforeMerge:      DefaultMaxDocsBeforeMerge,
		minLayerSize:            DefaultMinLayerSize,
		levelLogSize:            DefaultLevelLogSize,
		delDocsRatioBeforeMerge: DefaultDelDocsRatioBeforeMerge,
	}
}

func (l *LogMergePolicy) ComputeMergeCandidates(segmentMetas []*index.SegmentMeta) [][]index.SegmentID {
	candidateSegments := make([]*index.SegmentMeta, 0, len(segmentMetas))
	for _, segmentMeta := range segmentMetas {
		if segmentMeta.DocNum() < l.maxDocsBeforeMerge {
			candidateSegments = append(candidateSegments, segmentMeta)
		}
	}
	if len(segmentMetas) == 0 {
		return nil
	}

	sort.Slice(candidateSegments, func(i, j int) bool {
		return candidateSegments[i].DocNum() > candidateSegments[j].DocNum()
	})

	curMaxLogSize := math.MaxFloat64
	levels := map[float64][]*index.SegmentMeta{}
	for _, candidateSegment := range candidateSegments {
		segmentLogSize := float64(l.clipMinSize(candidateSegment.DocNum()))
		if segmentLogSize < curMaxLogSize-l.levelLogSize {
			curMaxLogSize = segmentLogSize
		}
		levels[curMaxLogSize] = append(levels[curMaxLogSize], candidateSegment)
	}

	candidates := make([][]index.SegmentID, 0, len(levels))
	for _, level := range levels {
		if len(level) < l.minNumSegments {
			continue
		}
		candidate := make([]index.SegmentID, 0, len(level))
		for _, segmentMeta := range level {
			candidate = append(candidate, segmentMeta.SegmentID)
		}
		candidates = append(candidates, candidate)
	}
	return candidates
}

func (l *LogMergePolicy) clipMinSize(size uint32) uint32 {
	if l.minLayerSize > size {
		return l.minLayerSize
	}
	return size
}
