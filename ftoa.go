// Copyright 2026 coalaura. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package byteconv

type floatInfo struct {
	mantbits uint
	expbits  uint
	bias     int
}

const (
	float32MantBits = 23
	float32ExpBits  = 8
	float32Bias     = -127
	float64MantBits = 52
	float64ExpBits  = 11
	float64Bias     = -1023
)

var (
	float32info = floatInfo{float32MantBits, float32ExpBits, float32Bias}
	float64info = floatInfo{float64MantBits, float64ExpBits, float64Bias}
)
