package util

import "encoding/binary"

type decoder struct {
	packet []byte
	offset int
}

func NewDecoder(packet []byte) *decoder {
	d := new(decoder)
	d.offset = 0
	d.packet = packet
	return d
}

func (d *decoder) SetPacket(packet []byte) {
	d.offset = 0
	d.packet = packet
}

func (d *decoder) LenLeft() int {
	return len(d.packet) - d.offset
}

func (d *decoder) GetByte() byte {
	data := d.packet[d.offset]
	d.offset++
	return data
}

func (d *decoder) GetBytes(n int) []byte {
	data := d.packet[d.offset : d.offset+n]
	d.offset += n
	return data
}

func (d *decoder) GetUint16() uint16 {
	return binary.BigEndian.Uint16(d.GetBytes(2))
}

func (d *decoder) GetUint32() uint32 {
	return binary.BigEndian.Uint32(d.GetBytes(4))
}

func (d *decoder) GetUint64() uint64 {
	return binary.BigEndian.Uint64(d.GetBytes(8))
}
