package services

func Utf8Decode(p []byte) (rune, int) {
	if len(p) == 0 {
		return 0, 0 // No bytes to decode
	}

	// The first byte tells us how many bytes the character uses.
	b0 := p[0]

	// 1. One-byte character (0xxxxxxx)
	if b0 < 0x80 {
		// This is an ASCII character, the rune is simply the byte itself.
		return rune(b0), 1
	}

	// 2. Two-byte character (110xxxxx 10xxxxxx)
	if b0 < 0xE0 {
		// Check for enough bytes and that the second byte is a continuation byte.
		if len(p) < 2 || (p[1]>>6) != 0b10 {
			return 0, 0 // Invalid or incomplete
		}
		// The `&` bitmask gets the lower 5 bits of b0 and lower 6 of b1.
		// The `|` bitwise OR combines them.
		r := rune(b0&0x1F)<<6 | rune(p[1]&0x3F)
		return r, 2
	}

	// 3. Three-byte character (1110xxxx 10xxxxxx 10xxxxxx)
	if b0 < 0xF0 {
		if len(p) < 3 || (p[1]>>6) != 0b10 || (p[2]>>6) != 0b10 {
			return 0, 0
		}
		// Reconstruct the rune from the three bytes.
		r := rune(b0&0x0F)<<12 | rune(p[1]&0x3F)<<6 | rune(p[2]&0x3F)
		return r, 3
	}

	// 4. Four-byte character (11110xxx 10xxxxxx 10xxxxxx 10xxxxxx)
	if b0 < 0xF8 {
		if len(p) < 4 || (p[1]>>6) != 0b10 || (p[2]>>6) != 0b10 || (p[3]>>6) != 0b10 {
			return 0, 0
		}
		// Reconstruct the rune from the four bytes.
		r := rune(b0&0x07)<<18 | rune(p[1]&0x3F)<<12 | rune(p[2]&0x3F)<<6 | rune(p[3]&0x3F)
		return r, 4
	}

	// If the first byte doesn't match any of the patterns, it's invalid.
	return 0, 0
}
