package qrcode

// bitBuffer accumulates a sequence of bits for QR data encoding.
type bitBuffer struct {
	data []byte
	len  int // total number of bits
}

// put appends the low `length` bits of value to the buffer.
func (b *bitBuffer) put(value int, length int) {
	for i := length - 1; i >= 0; i-- {
		byteIdx := b.len >> 3
		if byteIdx >= len(b.data) {
			b.data = append(b.data, 0)
		}
		if (value>>i)&1 == 1 {
			b.data[byteIdx] |= 0x80 >> (uint(b.len) & 7)
		}
		b.len++
	}
}

// putByte appends a full byte.
func (b *bitBuffer) putByte(v byte) {
	b.put(int(v), 8)
}

// bytes returns the accumulated bytes, padding to the next byte boundary.
func (b *bitBuffer) bytes() []byte {
	n := (b.len + 7) / 8
	if n > len(b.data) {
		return b.data
	}
	return b.data[:n]
}
