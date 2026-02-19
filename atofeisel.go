// Copyright 2026 coalaura. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package byteconv

import (
	"math"
	"math/bits"
)

func eiselLemire64(man uint64, exp10 int, neg bool) (f float64, ok bool) {
	if man == 0 {
		if neg {
			f = math.Float64frombits(0x8000000000000000) // Negative zero.
		}
		return f, true
	}
	pow, exp2, ok := pow10(exp10)
	if !ok {
		return 0, false
	}

	clz := bits.LeadingZeros64(man)
	man <<= uint(clz)
	retExp2 := uint64(exp2+63-float64Bias) - uint64(clz)

	xHi, xLo := bits.Mul64(man, pow.Hi)

	if xHi&0x1FF == 0x1FF && xLo+man < man {
		yHi, yLo := bits.Mul64(man, pow.Lo)
		mergedHi, mergedLo := xHi, xLo+yHi
		if mergedLo < xLo {
			mergedHi++
		}
		if mergedHi&0x1FF == 0x1FF && mergedLo+1 == 0 && yLo+man < man {
			return 0, false
		}
		xHi, xLo = mergedHi, mergedLo
	}

	msb := xHi >> 63
	retMantissa := xHi >> (msb + 9)
	retExp2 -= 1 ^ msb

	if xLo == 0 && xHi&0x1FF == 0 && retMantissa&3 == 1 {
		return 0, false
	}

	retMantissa += retMantissa & 1
	retMantissa >>= 1
	if retMantissa>>53 > 0 {
		retMantissa >>= 1
		retExp2 += 1
	}

	if retExp2-1 >= 0x7FF-1 {
		return 0, false
	}
	retBits := retExp2<<float64MantBits | retMantissa&(1<<float64MantBits-1)
	if neg {
		retBits |= 0x8000000000000000
	}
	return math.Float64frombits(retBits), true
}

func eiselLemire32(man uint64, exp10 int, neg bool) (f float32, ok bool) {
	if man == 0 {
		if neg {
			f = math.Float32frombits(0x80000000) // Negative zero.
		}
		return f, true
	}
	pow, exp2, ok := pow10(exp10)
	if !ok {
		return 0, false
	}

	clz := bits.LeadingZeros64(man)
	man <<= uint(clz)
	retExp2 := uint64(exp2+63-float32Bias) - uint64(clz)

	xHi, xLo := bits.Mul64(man, pow.Hi)

	if xHi&0x3FFFFFFFFF == 0x3FFFFFFFFF && xLo+man < man {
		yHi, yLo := bits.Mul64(man, pow.Lo)
		mergedHi, mergedLo := xHi, xLo+yHi
		if mergedLo < xLo {
			mergedHi++
		}
		if mergedHi&0x3FFFFFFFFF == 0x3FFFFFFFFF && mergedLo+1 == 0 && yLo+man < man {
			return 0, false
		}
		xHi, xLo = mergedHi, mergedLo
	}

	msb := xHi >> 63
	retMantissa := xHi >> (msb + 38)
	retExp2 -= 1 ^ msb

	if xLo == 0 && xHi&0x3FFFFFFFFF == 0 && retMantissa&3 == 1 {
		return 0, false
	}

	retMantissa += retMantissa & 1
	retMantissa >>= 1
	if retMantissa>>24 > 0 {
		retMantissa >>= 1
		retExp2 += 1
	}

	if retExp2-1 >= 0xFF-1 {
		return 0, false
	}
	retBits := retExp2<<float32MantBits | retMantissa&(1<<float32MantBits-1)
	if neg {
		retBits |= 0x80000000
	}
	return math.Float32frombits(uint32(retBits)), true
}
