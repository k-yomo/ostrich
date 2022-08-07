package termdict

import "github.com/k-yomo/ostrich/schema"

type TermDict map[schema.FieldID]map[string]*TermInfo
