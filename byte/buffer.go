package byte

import "encoding/binary"

type Buffer struct {
	Data  []byte
	Index int
}

func (buf *Buffer) Len() int {
	return len(buf.Data)
}

func (buf *Buffer) readByte() byte {
	buf.Index++
	return buf.Data[buf.Index-1]
}

func (buf *Buffer) readBytes(count int) []byte {
	buf.Index += count
	return buf.Data[buf.Index-count : buf.Index]
}

func (buf *Buffer) readU16() uint16 {
	buf.Index += 2
	return binary.BigEndian.Uint16(buf.Data[buf.Index-2 : buf.Index])
}

func (buf *Buffer) readU32() uint32 {
	buf.Index += 4
	return binary.BigEndian.Uint32(buf.Data[buf.Index-4 : buf.Index])
}

func (buf *Buffer) writeByte(b byte) {
	buf.Data = append(buf.Data, b)
	buf.Index++
}

func (buf *Buffer) writeBytes(b []byte) {
	buf.Data = append(buf.Data, b...)
	buf.Index += len(b)
}

func (buf *Buffer) writeU16(u uint16) {
	binary.BigEndian.PutUint16(buf.Data[buf.Index:buf.Index+2], u)
	buf.Index += 2
}

func (buf *Buffer) writeU32(u uint32) {
	binary.BigEndian.PutUint32(buf.Data[buf.Index:buf.Index+4], u)
	buf.Index += 4
}
