package query

import (
	"reflect"
	"testing"

	"github.com/k-yomo/ostrich/analyzer"
	"github.com/k-yomo/ostrich/reader"
	"github.com/k-yomo/ostrich/schema"
)

func TestParser_Parse(t *testing.T) {
	indexSchema := schema.NewSchema()
	f1 := indexSchema.AddTextField("f1", analyzer.DefaultAnalyzerName)
	f2 := indexSchema.AddTextField("f2", analyzer.DefaultAnalyzerName)
	parser := NewParser(indexSchema, []schema.FieldID{f1, f2})
	tests := []struct {
		name    string
		query   string
		want    reader.Query
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.Parse(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
