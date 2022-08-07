package postings

import (
	"bytes"
	"encoding/binary"
)

const FooterByteLength = 16

type Footer struct {
	postingsByteNum      uint64
	termFrequencyByteNum uint64
}

func (f *Footer) Write(buf *bytes.Buffer) {
	b := make([]byte, FooterByteLength)
	binary.LittleEndian.PutUint64(b[:8], f.postingsByteNum)
	binary.LittleEndian.PutUint64(b[8:], f.termFrequencyByteNum)
	buf.Write(b)
}

func ReadFooter(b []byte) *Footer {
	termFrequencyByteNum := binary.LittleEndian.Uint64(b[len(b)-8:])
	postingsByteNum := binary.LittleEndian.Uint64(b[len(b)-16 : len(b)-8])
	return &Footer{
		postingsByteNum:      postingsByteNum,
		termFrequencyByteNum: termFrequencyByteNum,
	}
}
