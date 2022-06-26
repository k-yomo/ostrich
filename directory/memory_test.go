package directory

import (
	"io"
	"reflect"
	"testing"
)

func TestMemoryDirectory_Read(t *testing.T) {
	t.Parallel()

	type args struct {
		path string
	}
	tests := []struct {
		name            string
		memoryDirectory *memoryDirectory
		args            args
		want            []byte
		wantErr         bool
	}{
		{
			name: "when path exist, it returns io.ReadCloser",
			memoryDirectory: func() *memoryDirectory {
				m := NewMemoryDirectory()
				writer, _ := m.OpenWrite("test")
				_, _ = writer.Write([]byte("abc"))
				return m
			}(),
			args: args{
				path: "test",
			},
			want: []byte("abc"),
		},
		{
			name:            "when path doesn't exist, it returns error",
			memoryDirectory: NewMemoryDirectory(),
			args: args{
				path: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.memoryDirectory.Read(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				b := make([]byte, len(tt.want))
				if _, err := got.Read(b); err != nil && err != io.EOF {
					t.Error(err)
					return
				}
				if !reflect.DeepEqual(b, tt.want) {
					t.Errorf("Read() got = %v, want %v", b, tt.want)
				}
			}
		})
	}
}

func Test_memoryDirectory_OpenWrite(t *testing.T) {
	t.Parallel()

	type args struct {
		path string
	}
	tests := []struct {
		name            string
		memoryDirectory *memoryDirectory
		args            args
		writeBytes      []byte
		wantErr         bool
	}{
		{
			name:            "writes to the path",
			memoryDirectory: NewMemoryDirectory(),
			writeBytes:      []byte("abc"),
			args: args{
				path: "test",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.memoryDirectory.OpenWrite(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenWrite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				got.Write(tt.writeBytes)
				gotBytes := tt.memoryDirectory.pathBytesMap[tt.args.path]
				if !reflect.DeepEqual(gotBytes, tt.writeBytes) {
					t.Errorf("OpenWrite(), Write() got = %v, want %v", gotBytes, tt.writeBytes)
				}
			}
		})
	}
}
