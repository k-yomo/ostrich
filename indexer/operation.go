package indexer

import (
	"github.com/k-yomo/ostrich/internal/opstamp"
	"github.com/k-yomo/ostrich/schema"
)

type AddOperation struct {
	opstamp  opstamp.OpStamp
	document *schema.Document
	result   func(error)
}
