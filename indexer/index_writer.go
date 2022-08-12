package indexer

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"unsafe"

	"github.com/k-yomo/go-batch"

	"github.com/k-yomo/ostrich/index"
	"github.com/k-yomo/ostrich/internal/opstamp"
	"github.com/k-yomo/ostrich/schema"
)

const (
	MaxThreadNum              = 8
	MarginInBytes         int = 1e6 // 1MB
	HeapSizeMin               = MarginInBytes * 3
	HeapSizeMax               = math.MaxUint32 - MarginInBytes
	MaxOperationQueueSize     = 10000
)

type IndexWriter struct {
	index *index.Index

	operationBatcher *batch.Batch[*AddOperation]
	segmentUpdater   *SegmentUpdater

	stamper          *opstamp.Stamper
	committedOpstamp opstamp.OpStamp

	closed bool
}

func NewIndexWriter(idx *index.Index, overallHeapBytes int) (*IndexWriter, error) {
	indexMeta, err := idx.LoadMeta()
	if err != nil {
		return nil, fmt.Errorf("load meta: %w", err)
	}

	currentOpstamp := indexMeta.Opstamp

	stamper := opstamp.NewStamper(currentOpstamp)
	segmentUpdater := NewSegmentUpdater(idx, indexMeta, stamper)

	i := &IndexWriter{
		index: idx,

		segmentUpdater: segmentUpdater,

		stamper:          stamper,
		committedOpstamp: currentOpstamp,
	}

	threadNum := int(math.Min(float64(runtime.GOMAXPROCS(0)), 8))
	b := batch.New(func(operations []*AddOperation) {
		err := i.indexDocuments(operations)
		for _, op := range operations {
			op.result(err)
		}
	})
	b.HandlerLimit = threadNum
	b.BundleCountThreshold = MaxOperationQueueSize
	b.BundleByteThreshold = HeapSizeMin
	b.BundleByteLimit = HeapSizeMax
	b.BufferedByteLimit = overallHeapBytes

	i.operationBatcher = b

	return i, nil
}

type AddDocumentResult struct {
	OpStamp opstamp.OpStamp
	// Result blocks until the operation is finished
	Result func() error
}

func (i *IndexWriter) AddDocument(document *schema.Document) *AddDocumentResult {
	if i.closed {
		return &AddDocumentResult{
			OpStamp: 0,
			Result: func() error {
				return errors.New("writer is already closed")
			},
		}
	}

	resultChan := make(chan error, 1)
	resultFunc := func() error {
		return <-resultChan
	}

	opStamp := i.stamper.Stamp()
	addOperation := &AddOperation{
		opstamp:  opStamp,
		document: document,
		result: func(err error) {
			resultChan <- err
		},
	}
	if err := i.operationBatcher.Add(addOperation, int(unsafe.Sizeof(addOperation))); err != nil {
		return &AddDocumentResult{
			OpStamp: opStamp,
			Result: func() error {
				return err
			},
		}
	}

	return &AddDocumentResult{
		OpStamp: opStamp,
		Result:  resultFunc,
	}
}

func (i *IndexWriter) indexDocuments(operations []*AddOperation) error {
	segment := i.index.NewSegment()
	indexSchema := segment.Schema()
	segmentWriter, err := newSegmentWriter(segment, indexSchema)
	if err != nil {
		return fmt.Errorf("initialize segment writer: %w", err)
	}
	for _, op := range operations {
		if err := segmentWriter.addDocument(op, indexSchema); err != nil {
			return fmt.Errorf("add document: %w", err)
		}
	}

	if err := segmentWriter.finalize(); err != nil {
		return fmt.Errorf("finalize segment: %w", err)
	}

	segmentWithMaxDoc := segment.WithMaxDoc(segmentWriter.maxDoc)
	segmentEntry := NewSegmentEntry(segmentWithMaxDoc.Meta())
	i.segmentUpdater.AddSegment(segmentEntry)

	return nil
}

func (i *IndexWriter) Commit() (opstamp.OpStamp, error) {
	i.operationBatcher.Flush()
	commitOpstamp := i.stamper.Stamp()
	if err := i.segmentUpdater.Commit(commitOpstamp); err != nil {
		return 0, err
	}
	go i.segmentUpdater.considerMergeOptions()
	return commitOpstamp, nil
}

func (i *IndexWriter) Merge(segmentIDs []index.SegmentID) (*index.SegmentMeta, error) {
	operation := NewMergeOperation(i.segmentUpdater.activeMeta.Opstamp, segmentIDs)
	return i.segmentUpdater.startMerge(operation)
}

func (i *IndexWriter) Close() error {
	i.closed = true
	i.operationBatcher.Flush()
	return nil
}
