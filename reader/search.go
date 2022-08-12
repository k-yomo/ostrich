package reader

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

type options struct {
	goroutineNum int
}

type SearchOption func(*options)

func SearchOptionConcurrent(goroutineNum int) SearchOption {
	return func(o *options) {
		o.goroutineNum = goroutineNum
	}
}

func Search[T any](searcher *Searcher, q Query, c Collector[T], opt ...SearchOption) (T, error) {
	opts := options{}
	for _, o := range opt {
		o(&opts)
	}

	var zeroT T
	weight, err := q.Weight(searcher, false)
	if err != nil {
		return zeroT, fmt.Errorf("initialize weight: %w", err)
	}

	results := make([]T, 0, len(searcher.segmentReaders))
	if opts.goroutineNum > 1 {
		resultChan := make(chan T, len(searcher.segmentReaders))
		eg := errgroup.Group{}
		eg.SetLimit(opts.goroutineNum)
		for i, segmentReader := range searcher.segmentReaders {
			i, segmentReader := i, segmentReader
			eg.Go(func() error {
				result, err := c.CollectSegment(weight, i, segmentReader)
				if err != nil {
					return fmt.Errorf("collect segment %s: %w", segmentReader.SegmentID, err)
				}
				resultChan <- result
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return zeroT, fmt.Errorf("concurrent collect: %w", err)
		}
		close(resultChan)
		for result := range resultChan {
			results = append(results, result)
		}
	} else {
		for i, segmentReader := range searcher.segmentReaders {
			result, err := c.CollectSegment(weight, i, segmentReader)
			if err != nil {
				return nil, fmt.Errorf("collect segment %s: %w", segmentReader.SegmentID, err)
			}
			results = append(results, result)
		}
	}
	return c.MergeResults(results), nil
}
