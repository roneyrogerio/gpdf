package font

import (
	"fmt"
	"sort"
	"strings"
)

// cmapEntry is a glyph-ID-to-rune pair used during CMap generation.
type cmapEntry struct {
	gid uint16
	r   rune
}

// GenerateToUnicodeCMap generates a ToUnicode CMap stream for PDF text
// extraction. The mapping associates glyph IDs (as CID values) with their
// corresponding Unicode code points, allowing PDF readers to extract
// searchable and copyable text from the document.
func GenerateToUnicodeCMap(mapping map[rune]uint16) []byte {
	if len(mapping) == 0 {
		return nil
	}

	// Build sorted list of (gid, rune) pairs for deterministic output.
	entries := make([]cmapEntry, 0, len(mapping))
	for r, gid := range mapping {
		entries = append(entries, cmapEntry{gid: gid, r: r})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].gid < entries[j].gid
	})

	var b strings.Builder
	writeCMapHeader(&b)
	writeBFCharEntries(&b, entries)
	writeCMapFooter(&b)

	return []byte(b.String())
}

// writeCMapHeader writes the CMap header boilerplate.
func writeCMapHeader(b *strings.Builder) {
	b.WriteString("/CIDInit /ProcSet findresource begin\n")
	b.WriteString("12 dict begin\n")
	b.WriteString("begincmap\n")
	b.WriteString("/CIDSystemInfo << /Registry (Adobe) /Ordering (UCS) /Supplement 0 >> def\n")
	b.WriteString("/CMapName /Adobe-Identity-UCS def\n")
	b.WriteString("/CMapType 2 def\n")
	b.WriteString("1 begincodespacerange\n")
	b.WriteString("<0000> <FFFF>\n")
	b.WriteString("endcodespacerange\n")
}

// writeBFCharEntries writes beginbfchar/endbfchar blocks.
// PDF limits each block to 100 entries, so we chunk accordingly.
func writeBFCharEntries(b *strings.Builder, entries []cmapEntry) {
	const maxPerBlock = 100

	for len(entries) > 0 {
		chunk := entries
		if len(chunk) > maxPerBlock {
			chunk = entries[:maxPerBlock]
		}
		entries = entries[len(chunk):]

		fmt.Fprintf(b, "%d beginbfchar\n", len(chunk))
		for _, e := range chunk {
			if e.r <= 0xFFFF {
				fmt.Fprintf(b, "<%04X> <%04X>\n", e.gid, e.r)
			} else {
				// Supplementary plane: encode as UTF-16 surrogate pair.
				high, low := utf16SurrogatePair(e.r)
				fmt.Fprintf(b, "<%04X> <%04X%04X>\n", e.gid, high, low)
			}
		}
		b.WriteString("endbfchar\n")
	}
}

// writeCMapFooter writes the CMap footer boilerplate.
func writeCMapFooter(b *strings.Builder) {
	b.WriteString("endcmap\n")
	b.WriteString("CMapName currentdict /CMap defineresource pop\n")
	b.WriteString("end\n")
	b.WriteString("end\n")
}

// utf16SurrogatePair returns the UTF-16 surrogate pair for a supplementary
// plane code point (U+10000 and above).
func utf16SurrogatePair(r rune) (high, low uint16) {
	r -= 0x10000
	high = 0xD800 + uint16(r>>10)
	low = 0xDC00 + uint16(r&0x3FF)
	return
}
