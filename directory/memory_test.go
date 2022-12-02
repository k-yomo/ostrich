package directory

import (
	"reflect"
	"testing"
)

func TestMemoryDirectory_AtomicRead(t *testing.T) {
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

			got, err := tt.memoryDirectory.AtomicRead(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("AtomicRead() got = %v, want %v", got, tt.want)
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
				_, err := got.Write(tt.writeBytes)
				if err != nil {
					t.Errorf("got.Write() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				gotBytes := tt.memoryDirectory.pathBytesMap[tt.args.path]
				if !reflect.DeepEqual(gotBytes, tt.writeBytes) {
					t.Errorf("OpenWrite(), Write() got = %v, want %v", gotBytes, tt.writeBytes)
				}
			}
		})
	}
}

func Test_memoryDirectory_Delete(t *testing.T) {
	t.Parallel()

	type args struct {
		path string
	}
	tests := []struct {
		name            string
		memoryDirectory *memoryDirectory
		args            args
		wantErr         bool
	}{
		{
			name: "delete path",
			memoryDirectory: func() *memoryDirectory {
				m := NewMemoryDirectory()
				m.AtomicWrite("test", []byte{})
				return m
			}(),
			args: args{path: "test"},
		},
		{
			name:            "delete not existing path won't cause error",
			memoryDirectory: NewMemoryDirectory(),
			args:            args{path: "not_exist"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.memoryDirectory.Delete(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			exists, err := tt.memoryDirectory.Exists(tt.args.path)
			if err != nil {
				t.Errorf("Delete(), Exists() error = %v, wantErr %v", err, false)
			}
			if exists {
				t.Errorf("Delete(), Exists() got = %v, want %v", exists, false)
			}
		})
	}
}
