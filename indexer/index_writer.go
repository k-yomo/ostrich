package indexer

import (
	"fmt"
	"math"
	"runtime"

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
	index             *index.Index
	heapSizePerThread int

	operationSender   chan<- []*AddOperation
	operationReceiver <-chan []*AddOperation
	segmentUpdater    *SegmentUpdater

	workerID  int
	ThreadNum int

	stamper          *opstamp.Stamper
	committedOpstamp opstamp.OpStamp
}

func NewIndexWriter(idx *index.Index, overallHeapBytes int) (*IndexWriter, error) {
	threadNum := int(math.Min(float64(runtime.GOMAXPROCS(0)), 8))
	heapSizePerThread := overallHeapBytes / threadNum
	if heapSizePerThread < HeapSizeMin {
		threadNum = int(math.Max(float64(overallHeapBytes/HeapSizeMin), 1))
	}
	if heapSizePerThread < HeapSizeMin {
		return nil, fmt.Errorf("heap size per thread needs to be at least %d", HeapSizeMin)
	}

	indexMeta, err := idx.LoadMetas()
	if err != nil {
		return nil, fmt.Errorf("load metas: %v", err)
	}

	currentOpstamp := indexMeta.Opstamp

	stamper := opstamp.NewStamper(currentOpstamp)
	segmentUpdater := NewSegmentUpdater(idx, indexMeta, stamper)
	operationChan := make(chan []*AddOperation, MaxOperationQueueSize)

	i := &IndexWriter{
		index: idx,

		operationSender:   operationChan,
		operationReceiver: operationChan,
		segmentUpdater:    segmentUpdater,

		workerID:  0,
		ThreadNum: threadNum,

		stamper:          stamper,
		committedOpstamp: currentOpstamp,
	}

	i.startWorkers()

	return i, nil
}

func (i *IndexWriter) AddDocument(document *schema.Document) opstamp.OpStamp {
	opStamp := i.stamper.Stamp()
	addOperation := &AddOperation{
		opstamp:  opStamp,
		document: document,
	}
	i.operationSender <- []*AddOperation{addOperation}

	return opStamp
}

func (i *IndexWriter) startWorkers() {
	for j := 0; j < i.ThreadNum; j++ {
		i.addIndexWorker()
	}
}

func (i *IndexWriter) addIndexWorker() {
	go func() {
		for {
			operations := <-i.operationReceiver
			if err := i.indexDocuments(operations); err != nil {
				// TODO: logging?
				fmt.Println(err)
			}
		}
	}()

	i.workerID += 1
}

func (i *IndexWriter) indexDocuments(operations []*AddOperation) error {
	segment := i.index.NewSegment()
	indexSchema := segment.Schema()
	segmentWriter, err := newSegmentWriter(i.heapSizePerThread, segment, indexSchema)
	if err != nil {
		return err
	}
	for _, op := range operations {
		if err := segmentWriter.addDocument(op, indexSchema); err != nil {
			return err
		}
	}

	if err := segmentWriter.finalize(); err != nil {
		return err
	}

	segmentWithMaxDoc := segment.WithMaxDoc(segmentWriter.maxDoc)
	segmentEntry := NewSegmentEntry(segmentWithMaxDoc.Meta())
	i.segmentUpdater.AddSegment(segmentEntry)

	return nil
}

func (i *IndexWriter) Commit() (opstamp.OpStamp, error) {
	commitOpstamp := i.stamper.Stamp()
	return commitOpstamp, i.segmentUpdater.Commit(commitOpstamp)
}
