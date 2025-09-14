package services

func Utf8Encode(r rune) []byte {
	// A slice to hold the encoded bytes.
	var encodedBytes []byte

	// 1. One-byte encoding (0x00 - 0x7F):
	// Characters in this range are single-byte ASCII.
	if r <= 0x7F {
		encodedBytes = append(encodedBytes, byte(r))
		return encodedBytes
	}

	// 2. Two-byte encoding (0x80 - 0x7FF):
	// The first byte has the form 110xxxxx, and the second is 10xxxxxx.
	if r <= 0x7FF {
		// First byte: 110xxxxx
		// Shift the rune 6 bits to the right to get the top 5 bits (0x07F).
		// OR with the 110xxxx pattern (0xC0).
		b1 := byte(0xC0 | (r >> 6))

		// Second byte: 10xxxxxx
		// Use a bitmask to get the lower 6 bits (0x3F).
		// OR with the 10xxxxxx pattern (0x80).
		b2 := byte(0x80 | (r & 0x3F))

		encodedBytes = append(encodedBytes, b1, b2)
		return encodedBytes
	}

	// 3. Three-byte encoding (0x800 - 0xFFFF):
	// First byte: 1110xxxx. Second: 10xxxxxx. Third: 10xxxxxx.
	if r <= 0xFFFF {
		// First byte: 1110xxxx
		b1 := byte(0xE0 | (r >> 12))

		// Second byte: 10xxxxxx
		b2 := byte(0x80 | ((r >> 6) & 0x3F))

		// Third byte: 10xxxxxx
		b3 := byte(0x80 | (r & 0x3F))

		encodedBytes = append(encodedBytes, b1, b2, b3)
		return encodedBytes
	}

	// 4. Four-byte encoding (0x10000 - 0x10FFFF):
	// First byte: 11110xxx. Second, third, fourth: 10xxxxxx.
	if r <= 0x10FFFF {
		// First byte: 11110xxx
		b1 := byte(0xF0 | (r >> 18))

		// Second byte: 10xxxxxx
		b2 := byte(0x80 | ((r >> 12) & 0x3F))

		// Third byte: 10xxxxxx
		b3 := byte(0x80 | ((r >> 6) & 0x3F))

		// Fourth byte: 10xxxxxx
		b4 := byte(0x80 | (r & 0x3F))

		encodedBytes = append(encodedBytes, b1, b2, b3, b4)
		return encodedBytes
	}

	// Fallback for invalid or out-of-range code points.
	return []byte{0xEF, 0xBF, 0xBD}
}
