package qrcode

import "errors"

// alphanumericTable maps characters to their Code 128 alphanumeric values.
var alphanumericTable [128]int

func init() {
	for i := range alphanumericTable {
		alphanumericTable[i] = -1
	}
	for i, c := range "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:" {
		alphanumericTable[c] = i
	}
}

// detectMode returns the best encoding mode for the given data.
func detectMode(data []byte) mode {
	allNum := true
	allAlpha := true
	for _, b := range data {
		if b < '0' || b > '9' {
			allNum = false
		}
		if b >= 128 || alphanumericTable[b] < 0 {
			allAlpha = false
		}
	}
	if allNum && len(data) > 0 {
		return modeNumeric
	}
	if allAlpha {
		return modeAlphanumeric
	}
	return modeByte
}

// selectVersion picks the smallest version that can hold the data
// at the given EC level.
func selectVersion(data []byte, m mode, ecLevel ErrorCorrectionLevel) (int, error) {
	dataLen := len(data)
	for v := 1; v <= 40; v++ {
		ec := ecTable[v-1][ecLevel]
		capacity := ec.dataCapacity() * 8 // in bits

		// Calculate the bit length of the encoded data.
		bits := 4 + charCountBits(v, m)
		switch m {
		case modeNumeric:
			groups := dataLen / 3
			remainder := dataLen % 3
			bits += groups*10 + numericRemainderBits(remainder)
		case modeAlphanumeric:
			groups := dataLen / 2
			bits += groups*11 + (dataLen%2)*6
		case modeByte:
			bits += dataLen * 8
		}

		if bits <= capacity {
			return v, nil
		}
	}
	return 0, errors.New("qrcode: data too long for any version")
}

func numericRemainderBits(r int) int {
	switch r {
	case 1:
		return 4
	case 2:
		return 7
	}
	return 0
}

// encodeData produces the final data codewords (data + EC interleaved) for the QR code.
func encodeData(data []byte, ecLevel ErrorCorrectionLevel, version int) ([]byte, error) {
	m := detectMode(data)
	ec := ecTable[version-1][ecLevel]
	totalData := ec.dataCapacity()

	// Build the bit stream.
	buf := &bitBuffer{}

	// Mode indicator (4 bits).
	buf.put(modeIndicator(m), 4)

	// Character count indicator.
	buf.put(len(data), charCountBits(version, m))

	// Encode data segments.
	switch m {
	case modeNumeric:
		encodeNumeric(buf, data)
	case modeAlphanumeric:
		encodeAlphanumeric(buf, data)
	case modeByte:
		encodeByte(buf, data)
	}

	// Terminator (up to 4 zero bits, no more than capacity).
	capacityBits := totalData * 8
	termLen := 4
	if buf.len+termLen > capacityBits {
		termLen = capacityBits - buf.len
	}
	if termLen > 0 {
		buf.put(0, termLen)
	}

	// Pad to byte boundary.
	if buf.len%8 != 0 {
		buf.put(0, 8-buf.len%8)
	}

	// Pad codewords: alternate 0xEC and 0x11.
	padBytes := []byte{0xEC, 0x11}
	for i := 0; buf.len/8 < totalData; i++ {
		buf.putByte(padBytes[i%2])
	}

	dataBytes := buf.bytes()

	// Split into blocks and compute EC for each.
	return interleave(dataBytes, ec)
}

// interleave splits data into blocks, generates EC codewords, and
// interleaves the result per QR spec.
func interleave(data []byte, ec ecBlockInfo) ([]byte, error) {
	type block struct {
		data []byte
		ec   []byte
	}

	totalBlocks := ec.group1Blocks + ec.group2Blocks
	blocks := make([]block, totalBlocks)

	offset := 0
	for i := 0; i < ec.group1Blocks; i++ {
		blocks[i].data = data[offset : offset+ec.group1Data]
		offset += ec.group1Data
	}
	for i := 0; i < ec.group2Blocks; i++ {
		blocks[ec.group1Blocks+i].data = data[offset : offset+ec.group2Data]
		offset += ec.group2Data
	}

	// Compute EC codewords for each block.
	for i := range blocks {
		blocks[i].ec = rsEncode(blocks[i].data, ec.ecPerBlock)
	}

	// Interleave data codewords.
	var result []byte
	maxData := ec.group1Data
	if ec.group2Data > maxData {
		maxData = ec.group2Data
	}
	for col := 0; col < maxData; col++ {
		for _, b := range blocks {
			if col < len(b.data) {
				result = append(result, b.data[col])
			}
		}
	}

	// Interleave EC codewords.
	for col := 0; col < ec.ecPerBlock; col++ {
		for _, b := range blocks {
			if col < len(b.ec) {
				result = append(result, b.ec[col])
			}
		}
	}

	return result, nil
}

// encodeNumeric encodes digits in groups of 3.
func encodeNumeric(buf *bitBuffer, data []byte) {
	for i := 0; i < len(data); i += 3 {
		end := i + 3
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]
		val := 0
		for _, c := range chunk {
			val = val*10 + int(c-'0')
		}
		bits := numericRemainderBits(len(chunk))
		if len(chunk) == 3 {
			bits = 10
		}
		buf.put(val, bits)
	}
}

// encodeAlphanumeric encodes characters in pairs.
func encodeAlphanumeric(buf *bitBuffer, data []byte) {
	for i := 0; i < len(data); i += 2 {
		if i+1 < len(data) {
			val := alphanumericTable[data[i]]*45 + alphanumericTable[data[i+1]]
			buf.put(val, 11)
		} else {
			buf.put(alphanumericTable[data[i]], 6)
		}
	}
}

// encodeByte encodes each byte directly.
func encodeByte(buf *bitBuffer, data []byte) {
	for _, b := range data {
		buf.putByte(b)
	}
}
