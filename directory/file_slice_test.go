package directory

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/k-yomo/ostrich/pkg/xrange"
)

func TestFileSlice_Read(t *testing.T) {
	t.Parallel()

	type args struct {
		from int
		size int
	}
	tests := []struct {
		name      string
		fileSlice *FileSlice
		args      args
		want      []byte
		wantErr   bool
	}{
		{
			name:      "read all",
			fileSlice: NewFileSlice(bytes.NewReader([]byte{1, 2, 3, 4}), func() error { return nil }),
			args: args{
				from: 0,
				size: 4,
			},
			want: []byte{1, 2, 3, 4},
		},
		{
			name:      "read from middle",
			fileSlice: NewFileSlice(bytes.NewReader([]byte{1, 2, 3, 4}), func() error { return nil }),
			args: args{
				from: 2,
				size: 2,
			},
			want: []byte{3, 4},
		},
		{
			name:      "reach EOF",
			fileSlice: NewFileSlice(bytes.NewReader([]byte{1, 2, 3, 4}), func() error { return nil }),
			args: args{
				from: 2,
				size: 3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fileSlice.Read(tt.args.from, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileSlice_Slice(t *testing.T) {
	t.Parallel()

	fileSlice := NewFileSlice(bytes.NewReader([]byte{1, 2, 3, 4}), func() error { return nil })
	defer fileSlice.Close()
	buf := make([]byte, 2)
	dataReader := fileSlice.Slice(xrange.Range{From: 1, To: 2}).Reader()
	_, err := dataReader.Read(buf)
	if err != nil {
		t.Errorf("Slice(), Reader(), Read() = %v, want %v", err, nil)
		return
	}
	want := []byte{2, 3}
	if !reflect.DeepEqual(buf, want) {
		t.Errorf("Slice() = %v, want %v", buf, want)
	}

	buf = make([]byte, 1)
	_, err = dataReader.Read(buf)
	if err == nil {
		t.Errorf("Slice(), Reader(), Read() = %v, want %v", nil, "io.EOF")
		return
	}
}
