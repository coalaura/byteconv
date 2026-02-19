// Copyright 2026 coalaura. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package byteconv

import "math"

var optimize = true

// special returns the floating-point value for inf, infinity, and NaN.
func special(s []byte) (f float64, n int, ok bool) {
	if len(s) == 0 {
		return 0, 0, false
	}
	sign := 1
	nsign := 0
	switch s[0] {
	case '+', '-':
		if s[0] == '-' {
			sign = -1
		}
		nsign = 1
		s = s[1:]
		if len(s) == 0 {
			return 0, 0, false
		}
		fallthrough // A signed value MUST be infinity, never NaN.
	case 'i', 'I':
		if len(s) >= 3 && (s[0]|0x20) == 'i' && (s[1]|0x20) == 'n' && (s[2]|0x20) == 'f' {
			if len(s) >= 8 && (s[3]|0x20) == 'i' && (s[4]|0x20) == 'n' && (s[5]|0x20) == 'i' && (s[6]|0x20) == 't' && (s[7]|0x20) == 'y' {
				return math.Inf(sign), nsign + 8, true
			}
			return math.Inf(sign), nsign + 3, true
		}
	case 'n', 'N':
		if len(s) >= 3 && (s[0]|0x20) == 'n' && (s[1]|0x20) == 'a' && (s[2]|0x20) == 'n' {
			return math.NaN(), 3, true
		}
	}
	return 0, 0, false
}

func (b *decimal) set(s []byte) (ok bool) {
	lenS := len(s)
	if lenS == 0 {
		return false
	}
	_ = s[lenS-1] // BCE Hint

	i := 0
	b.neg = false
	b.trunc = false

	if s[i] == '+' {
		i++
	} else if s[i] == '-' {
		i++
		b.neg = true
	}

	sawdot := false
	sawdigits := false

	// Register-only math. No memory lookups for large floats.
	for ; i < lenS; i++ {
		c := s[i]
		if c-'0' <= 9 {
			sawdigits = true
			if c == '0' && b.nd == 0 {
				b.dp--
				continue
			}
			if b.nd < len(b.d) {
				b.d[b.nd] = c
				b.nd++
			} else if c != '0' {
				b.trunc = true
			}
		} else if c == '.' {
			if sawdot {
				return false
			}
			sawdot = true
			b.dp = b.nd
		} else if c == '_' {
			continue
		} else {
			break
		}
	}

	if !sawdigits {
		return false
	}
	if !sawdot {
		b.dp = b.nd
	}

	if i < lenS && (s[i]|0x20) == 'e' {
		i++
		if i >= lenS {
			return false
		}
		esign := 1
		if s[i] == '+' {
			i++
		} else if s[i] == '-' {
			i++
			esign = -1
		}
		if i >= lenS || s[i]-'0' > 9 {
			return false
		}
		e := 0
		for ; i < lenS; i++ {
			c := s[i]
			if c == '_' {
				continue
			}
			if c-'0' <= 9 {
				if e < 10000 {
					e = e*10 + int(c-'0')
				}
			} else {
				break
			}
		}
		b.dp += e * esign
	}

	return i == lenS
}

func readFloat(s []byte) (mantissa uint64, exp int, neg, trunc, hex bool, i int, ok bool) {
	lenS := len(s)
	if lenS == 0 {
		return
	}
	_ = s[lenS-1] // BCE Hint
	underscores := false

	if s[i] == '+' {
		i++
	} else if s[i] == '-' {
		i++
		neg = true
	}

	if i >= lenS {
		return
	}

	if i+2 < lenS && s[i] == '0' && (s[i+1]|0x20) == 'x' {
		hex = true
		i += 2
	}

	sawdot := false
	sawdigits := false
	nd := 0
	ndMant := 0
	dp := 0

	// Split the hot paths. This allows the compiler to perfectly inline
	// and unroll the decimal path without branch pollution from hex checks.
	if !hex {
		// DECIMAL HOT PATH (Register Arithmetic Only)
		for ; i < lenS; i++ {
			c := s[i]
			if c-'0' <= 9 {
				sawdigits = true
				if c == '0' && nd == 0 {
					dp--
					continue
				}
				nd++
				if ndMant < 19 {
					mantissa = mantissa*10 + uint64(c-'0')
					ndMant++
				} else if c != '0' {
					trunc = true
				}
			} else if c == '.' {
				if sawdot {
					break
				}
				sawdot = true
				dp = nd
			} else if c == '_' {
				underscores = true
			} else {
				break
			}
		}
	} else {
		// HEX HOT PATH
		for ; i < lenS; i++ {
			c := s[i]
			v := c
			isDigit := false

			if c-'0' <= 9 {
				v = c - '0'
				isDigit = true
			} else if (c|0x20)-'a' <= 5 {
				v = (c | 0x20) - 'a' + 10
				isDigit = true
			}

			if isDigit {
				sawdigits = true
				if v == 0 && nd == 0 {
					dp--
					continue
				}
				nd++
				if ndMant < 16 {
					mantissa = mantissa<<4 | uint64(v)
					ndMant++
				} else if v != 0 {
					trunc = true
				}
			} else if c == '.' {
				if sawdot {
					break
				}
				sawdot = true
				dp = nd
			} else if c == '_' {
				underscores = true
			} else {
				break
			}
		}
	}

	if !sawdigits {
		return
	}
	if !sawdot {
		dp = nd
	}

	if hex {
		dp *= 4
		ndMant *= 4
	}

	expChar := byte('e')
	if hex {
		expChar = 'p'
	}

	if i < lenS && (s[i]|0x20) == expChar {
		i++
		if i >= lenS {
			return
		}
		esign := 1
		if s[i] == '+' {
			i++
		} else if s[i] == '-' {
			i++
			esign = -1
		}
		if i >= lenS || s[i]-'0' > 9 {
			return
		}
		e := 0
		for ; i < lenS; i++ {
			c := s[i]
			if c-'0' <= 9 {
				if e < 10000 {
					e = e*10 + int(c-'0')
				}
			} else if c == '_' {
				underscores = true
			} else {
				break
			}
		}
		dp += e * esign
	} else if hex {
		return
	}

	if mantissa != 0 {
		exp = dp - ndMant
	}

	if underscores && !underscoreOK(s[:i]) {
		return
	}

	ok = true
	return
}

var powtab = []int{1, 3, 6, 9, 13, 16, 19, 23, 26}

func (d *decimal) floatBits(flt *floatInfo) (b uint64, overflow bool) {
	var exp int
	var mant uint64

	if d.nd == 0 {
		mant = 0
		exp = flt.bias
		goto out
	}

	if d.dp > 310 {
		goto overflow
	}
	if d.dp < -330 {
		mant = 0
		exp = flt.bias
		goto out
	}

	exp = 0
	for d.dp > 0 {
		var n int
		if d.dp >= len(powtab) {
			n = 27
		} else {
			n = powtab[d.dp]
		}
		d.Shift(-n)
		exp += n
	}
	for d.dp < 0 || d.dp == 0 && d.d[0] < '5' {
		var n int
		if -d.dp >= len(powtab) {
			n = 27
		} else {
			n = powtab[-d.dp]
		}
		d.Shift(n)
		exp -= n
	}

	exp--

	if exp < flt.bias+1 {
		n := flt.bias + 1 - exp
		d.Shift(-n)
		exp += n
	}

	if exp-flt.bias >= 1<<flt.expbits-1 {
		goto overflow
	}

	d.Shift(int(1 + flt.mantbits))
	mant = d.RoundedInteger()

	if mant == 2<<flt.mantbits {
		mant >>= 1
		exp++
		if exp-flt.bias >= 1<<flt.expbits-1 {
			goto overflow
		}
	}

	if mant&(1<<flt.mantbits) == 0 {
		exp = flt.bias
	}
	goto out

overflow:
	mant = 0
	exp = 1<<flt.expbits - 1 + flt.bias
	overflow = true

out:
	bits := mant & (uint64(1)<<flt.mantbits - 1)
	bits |= uint64((exp-flt.bias)&(1<<flt.expbits-1)) << flt.mantbits
	if d.neg {
		bits |= 1 << flt.mantbits << flt.expbits
	}
	return bits, overflow
}

var float64pow10 = []float64{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
	1e20, 1e21, 1e22,
}
var float32pow10 = []float32{1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10}

func atof64exact(mantissa uint64, exp int, neg bool) (f float64, ok bool) {
	if mantissa>>float64info.mantbits != 0 {
		return
	}
	f = float64(mantissa)
	if neg {
		f = -f
	}
	switch {
	case exp == 0:
		return f, true
	case exp > 0 && exp <= 15+22:
		if exp > 22 {
			f *= float64pow10[exp-22]
			exp = 22
		}
		if f > 1e15 || f < -1e15 {
			return
		}
		return f * float64pow10[exp], true
	case exp < 0 && exp >= -22:
		return f / float64pow10[-exp], true
	}
	return
}

func atof32exact(mantissa uint64, exp int, neg bool) (f float32, ok bool) {
	if mantissa>>23 != 0 {
		return
	}
	f = float32(mantissa)
	if neg {
		f = -f
	}
	switch {
	case exp == 0:
		return f, true
	case exp > 0 && exp <= 7+10:
		if exp > 10 {
			f *= float32pow10[exp-10]
			exp = 10
		}
		if f > 1e7 || f < -1e7 {
			return
		}
		return f * float32pow10[exp], true
	case exp < 0 && exp >= -10:
		return f / float32pow10[-exp], true
	}
	return
}

func atofHex(s []byte, flt *floatInfo, mantissa uint64, exp int, neg, trunc bool) (float64, error) {
	maxExp := 1<<flt.expbits + flt.bias - 2
	minExp := flt.bias + 1
	exp += int(flt.mantbits)

	for mantissa != 0 && mantissa>>(flt.mantbits+2) == 0 {
		mantissa <<= 1
		exp--
	}
	if trunc {
		mantissa |= 1
	}
	for mantissa>>(1+flt.mantbits+2) != 0 {
		mantissa = mantissa>>1 | mantissa&1
		exp++
	}

	for mantissa > 1 && exp < minExp-2 {
		mantissa = mantissa>>1 | mantissa&1
		exp++
	}

	round := mantissa & 3
	mantissa >>= 2
	round |= mantissa & 1
	exp += 2
	if round == 3 {
		mantissa++
		if mantissa == 1<<(1+flt.mantbits) {
			mantissa >>= 1
			exp++
		}
	}

	if mantissa>>flt.mantbits == 0 {
		exp = flt.bias
	}
	var err error
	if exp > maxExp {
		mantissa = 1 << flt.mantbits
		exp = maxExp + 1
		err = ErrRange
	}

	bits := mantissa & (1<<flt.mantbits - 1)
	bits |= uint64((exp-flt.bias)&(1<<flt.expbits-1)) << flt.mantbits
	if neg {
		bits |= 1 << flt.mantbits << flt.expbits
	}
	if flt == &float32info {
		return float64(math.Float32frombits(uint32(bits))), err
	}
	return math.Float64frombits(bits), err
}

func atof32(s []byte) (f float32, n int, err error) {
	if val, n, ok := special(s); ok {
		return float32(val), n, nil
	}

	mantissa, exp, neg, trunc, hex, n, ok := readFloat(s)
	if !ok {
		return 0, n, ErrSyntax
	}

	if hex {
		flt64, err := atofHex(s[:n], &float32info, mantissa, exp, neg, trunc)
		return float32(flt64), n, err
	}

	if optimize {
		if !trunc {
			if f, ok := atof32exact(mantissa, exp, neg); ok {
				return f, n, nil
			}
		}
		f, ok := eiselLemire32(mantissa, exp, neg)
		if ok {
			if !trunc {
				return f, n, nil
			}
			fUp, ok := eiselLemire32(mantissa+1, exp, neg)
			if ok && f == fUp {
				return f, n, nil
			}
		}
	}

	var d decimal
	if !d.set(s[:n]) {
		return 0, n, ErrSyntax
	}
	b, ovf := d.floatBits(&float32info)
	f = math.Float32frombits(uint32(b))
	if ovf {
		err = ErrRange
	}
	return f, n, err
}

func atof64(s []byte) (f float64, n int, err error) {
	if val, n, ok := special(s); ok {
		return val, n, nil
	}

	mantissa, exp, neg, trunc, hex, n, ok := readFloat(s)
	if !ok {
		return 0, n, ErrSyntax
	}

	if hex {
		f, err := atofHex(s[:n], &float64info, mantissa, exp, neg, trunc)
		return f, n, err
	}

	if optimize {
		if !trunc {
			if f, ok := atof64exact(mantissa, exp, neg); ok {
				return f, n, nil
			}
		}
		f, ok := eiselLemire64(mantissa, exp, neg)
		if ok {
			if !trunc {
				return f, n, nil
			}
			fUp, ok := eiselLemire64(mantissa+1, exp, neg)
			if ok && f == fUp {
				return f, n, nil
			}
		}
	}

	var d decimal
	if !d.set(s[:n]) {
		return 0, n, ErrSyntax
	}
	b, ovf := d.floatBits(&float64info)
	f = math.Float64frombits(b)
	if ovf {
		err = ErrRange
	}
	return f, n, err
}

func ParseFloat(s []byte, bitSize int) (float64, error) {
	f, n, err := parseFloatPrefix(s, bitSize)
	if n != len(s) {
		return 0, ErrSyntax
	}
	return f, err
}

func parseFloatPrefix(s []byte, bitSize int) (float64, int, error) {
	if bitSize == 32 {
		f, n, err := atof32(s)
		return float64(f), n, err
	}
	return atof64(s)
}
