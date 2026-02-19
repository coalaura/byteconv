// Copyright 2026 coalaura. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package byteconv

// ParseComplex converts the []byte s to a complex number
// with the precision specified by bitSize: 64 for complex64, or 128 for complex128.
// When bitSize=64, the result still has type complex128, but it will be
// convertible to complex64 without changing its value.
func ParseComplex(s []byte, bitSize int) (complex128, error) {
	size := 64
	if bitSize == 64 {
		size = 32 // complex64 uses float32 parts
	}

	lenS := len(s)
	if lenS == 0 {
		return 0, ErrSyntax
	}

	// Remove parentheses without re-evaluating length continuously.
	_ = s[lenS-1] // BCE Hint
	if lenS >= 2 && s[0] == '(' && s[lenS-1] == ')' {
		s = s[1 : lenS-1]
		lenS -= 2
	}

	var pending error

	// Read real part (possibly imaginary part if followed by 'i').
	re, n, err := parseFloatPrefix(s, size)
	if err != nil {
		if err != ErrRange {
			return 0, err
		}
		pending = err
	}
	s = s[n:]
	lenS -= n

	// If we have nothing left, we're done.
	if lenS == 0 {
		return complex(re, 0), pending
	}

	_ = s[lenS-1] // BCE Hint for remaining slice
	switch s[0] {
	case '+':
		if lenS > 1 && s[1] != '+' {
			s = s[1:]
			lenS--
		}
	case '-':
		// ok
	case 'i':
		if lenS == 1 {
			return complex(0, re), pending
		}
		fallthrough
	default:
		return 0, ErrSyntax
	}

	// Read imaginary part.
	im, n, err := parseFloatPrefix(s, size)
	if err != nil {
		if err != ErrRange {
			return 0, err
		}
		pending = err
	}
	s = s[n:]

	if len(s) != 1 || s[0] != 'i' {
		return 0, ErrSyntax
	}

	return complex(re, im), pending
}