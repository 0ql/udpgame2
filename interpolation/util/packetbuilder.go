package util

import (
	"encoding/binary"
)

type Builder struct {
	Buf   []byte
	buf16 []byte
	buf32 []byte
	buf64 []byte
}

func NewBuilder() *Builder {
	b := new(Builder)
	b.buf16 = make([]byte, 2)
	b.buf32 = make([]byte, 4)
	b.buf64 = make([]byte, 8)
	return b
}

func (b *Builder) Reset() {
	b.Buf = make([]byte, 0)
}

func (b *Builder) AddBytes(slice []byte) {
	b.Buf = append(b.Buf, slice...)
}

func (b *Builder) AddByte(byt byte) {
	b.Buf = append(b.Buf, byt)
}

func (b *Builder) AddUint16(v uint16) {
	binary.BigEndian.PutUint16(b.buf16, v)
	b.Buf = append(b.Buf, b.buf16...)
}

func (b *Builder) AddUint32(v uint32) {
	binary.BigEndian.PutUint32(b.buf32, v)
	b.Buf = append(b.Buf, b.buf32...)
}

func (b *Builder) AddUint64(v uint64) {
	binary.BigEndian.PutUint64(b.buf64, v)
	b.Buf = append(b.Buf, b.buf64...)
}
