package index

type Document struct {
	fields map[string]interface{}
}

type DocID uint32

type DocAddress struct {
	SegmentOrd int
	DocID      DocID
}

type DocSet interface {
	Advance() (DocID, error)
	Doc() (DocID, error)
}
