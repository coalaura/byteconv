// Copyright 2026 coalaura. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package byteconv

const intSize = 32 << (^uint(0) >> 63)

// IntSize is the size in bits of an int or uint value.
const IntSize = intSize

const _F = 255

// decodeTable maps ASCII bytes to their integer values.
// 0xFF (_F) represents an invalid character.
// This completely eliminates branching (if/else) for character validation.
var decodeTable = [256]byte{
	_F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F,
	_F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F,
	_F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, _F, _F, _F, _F, _F, _F,
	_F, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
	25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, _F, _F, _F, _F, _F,
	_F, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
	25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, _F, _F, _F, _F, _F,
	_F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F,
	_F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F,
	_F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F,
	_F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F,
	_F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F,
	_F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F,
	_F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F,
	_F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F, _F,
}

// ParseUint is like ParseInt but for unsigned numbers.
// A sign prefix is not permitted.
func ParseUint(s []byte, base int, bitSize int) (uint64, error) {
	if len(s) == 0 {
		return 0, ErrSyntax
	}

	base0 := base == 0
	s0 := s

	switch {
	case 2 <= base && base <= 36:
		// valid base; nothing to do
	case base == 0:
		// Look for octal, hex prefix.
		base = 10
		if s[0] == '0' {
			switch {
			case len(s) >= 3 && (s[1]|0x20) == 'b':
				base = 2
				s = s[2:]
			case len(s) >= 3 && (s[1]|0x20) == 'o':
				base = 8
				s = s[2:]
			case len(s) >= 3 && (s[1]|0x20) == 'x':
				base = 16
				s = s[2:]
			default:
				base = 8
				s = s[1:]
			}
		}
	default:
		return 0, ErrBase
	}

	if bitSize == 0 {
		bitSize = IntSize
	} else if bitSize < 0 || bitSize > 64 {
		return 0, ErrBitSize
	}

	// Calculate max limits BEFORE fast paths so we can trap bitSize overflows
	maxVal := uint64(1)<<uint(bitSize) - 1
	if bitSize == 64 {
		maxVal = ^uint64(0)
	}

	// ==========================================
	// FAST PATHS (Bounds Check Eliminated)
	// ==========================================
	if len(s) > 0 {
		if base == 10 && len(s) <= 19 {
			_ = s[len(s)-1] // BCE Hint
			var n uint64
			for i := 0; i < len(s); i++ {
				v := decodeTable[s[i]]
				if v >= 10 {
					goto slowPath
				}
				n = n*10 + uint64(v)
			}
			if n > maxVal {
				return maxVal, ErrRange
			}
			return n, nil
		}
		if base == 16 && len(s) <= 16 {
			_ = s[len(s)-1]
			var n uint64
			for i := 0; i < len(s); i++ {
				v := decodeTable[s[i]]
				if v >= 16 {
					goto slowPath
				}
				n = n<<4 | uint64(v)
			}
			if n > maxVal {
				return maxVal, ErrRange
			}
			return n, nil
		}
		if base == 8 && len(s) <= 21 {
			_ = s[len(s)-1]
			var n uint64
			for i := 0; i < len(s); i++ {
				v := decodeTable[s[i]]
				if v >= 8 {
					goto slowPath
				}
				n = n<<3 | uint64(v)
			}
			if n > maxVal {
				return maxVal, ErrRange
			}
			return n, nil
		}
		if base == 2 && len(s) <= 64 {
			_ = s[len(s)-1]
			var n uint64
			for i := 0; i < len(s); i++ {
				v := decodeTable[s[i]]
				if v >= 2 {
					goto slowPath
				}
				n = n<<1 | uint64(v)
			}
			if n > maxVal {
				return maxVal, ErrRange
			}
			return n, nil
		}
	}

slowPath:
	var cutoff uint64
	switch base {
	case 10:
		cutoff = maxUint64/10 + 1
	case 16:
		cutoff = maxUint64/16 + 1
	default:
		cutoff = maxUint64/uint64(base) + 1
	}

	underscores := false
	var n uint64

	// Protect BCE hint against "" slice if stripped or bypassed to slowPath
	if len(s) > 0 {
		_ = s[len(s)-1]
		for i := 0; i < len(s); i++ {
			c := s[i]
			v := decodeTable[c]

			if v >= byte(base) {
				if c == '_' && base0 {
					underscores = true
					continue
				}
				return 0, ErrSyntax
			}

			if n >= cutoff {
				return maxVal, ErrRange
			}
			n *= uint64(base)

			n1 := n + uint64(v)
			if n1 < n || n1 > maxVal {
				return maxVal, ErrRange
			}
			n = n1
		}
	}

	if underscores && !underscoreOK(s0) {
		return 0, ErrSyntax
	}

	return n, nil
}

// ParseInt interprets a byte slice s in the given base (0, 2 to 36) and
// bit size (0 to 64) and returns the corresponding value i.
func ParseInt(s []byte, base int, bitSize int) (i int64, err error) {
	if len(s) == 0 {
		return 0, ErrSyntax
	}

	neg := false
	if s[0] == '+' {
		s = s[1:]
	} else if s[0] == '-' {
		neg = true
		s = s[1:]
	}

	var un uint64
	un, err = ParseUint(s, base, bitSize)
	if err != nil && err != ErrRange {
		return 0, err
	}

	if bitSize == 0 {
		bitSize = IntSize
	}

	cutoff := uint64(1 << uint(bitSize-1))
	if !neg && un >= cutoff {
		return int64(cutoff - 1), ErrRange
	}
	if neg && un > cutoff {
		return -int64(cutoff), ErrRange
	}
	n := int64(un)
	if neg {
		n = -n
	}
	return n, nil
}

// Atoi is equivalent to ParseInt(s, 10, 0), converted to type int.
func Atoi(s []byte) (int, error) {
	sLen := len(s)
	if sLen == 0 {
		return 0, ErrSyntax
	}

	s0 := s // Save original slice containing the sign for the slow path

	neg := false
	if s[0] == '-' {
		neg = true
		s = s[1:]
		sLen--
	} else if s[0] == '+' {
		s = s[1:]
		sLen--
	}

	if sLen == 0 {
		return 0, ErrSyntax
	}

	// Fast path max length: 18 digits safely fit in 64-bit, 9 in 32-bit.
	limit := 18
	if intSize == 32 {
		limit = 9
	}

	if sLen <= limit {
		_ = s[sLen-1] // BCE Hint
		n := 0
		for i := 0; i < sLen; i++ {
			v := int(decodeTable[s[i]])
			if v >= 10 {
				goto slowPath // Fallback if invalid or underscore
			}
			n = n*10 + v
		}
		if neg {
			return -n, nil
		}
		return n, nil
	}

slowPath:
	// Use s0 here, NOT s, because ParseInt needs to handle the sign logic itself.
	i64, err := ParseInt(s0, 10, 0)
	return int(i64), err
}

// underscoreOK reports whether the underscores in s are allowed.
func underscoreOK(s []byte) bool {
	saw := '^'
	i := 0

	if len(s) >= 1 && (s[0] == '-' || s[0] == '+') {
		s = s[1:]
	}

	hex := false
	if len(s) >= 2 && s[0] == '0' {
		switch s[1] | 0x20 {
		case 'b', 'o', 'x':
			i = 2
			saw = '0'
			hex = (s[1] | 0x20) == 'x'
		}
	}

	_ = s[len(s)-1] // BCE Hint
	for ; i < len(s); i++ {
		c := s[i]

		isDigit := false
		if hex {
			isDigit = decodeTable[c] < 16
		} else {
			isDigit = decodeTable[c] < 10
		}

		if isDigit {
			saw = '0'
			continue
		}
		if c == '_' {
			if saw != '0' {
				return false
			}
			saw = '_'
			continue
		}
		if saw == '_' {
			return false
		}
		saw = '!'
	}
	return saw != '_'
}
