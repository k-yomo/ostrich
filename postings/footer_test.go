package postings

import (
	"bytes"
	"reflect"
	"testing"
)

func TestFooter(t *testing.T) {
	f := &Footer{
		postingsByteNum:      350,
		termFrequencyByteNum: 200,
	}
	buf := bytes.NewBuffer([]byte{})
	f.Write(buf)

	got := readFooter(buf.Bytes())
	if !reflect.DeepEqual(f, got) {
		t.Errorf("readFooter() got = %v, want %v", got, f)
	}
}
