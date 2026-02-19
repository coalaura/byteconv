// Copyright 2026 coalaura. All rights reserved.

package byteconv_bench

import (
	"strconv"
	"testing"

	"github.com/coalaura/byteconv"
)

var (
	// Bools
	boolTrueStr  = "true"
	boolTrue     = []byte(boolTrueStr)
	boolFalseStr = "false"
	boolFalse    = []byte(boolFalseStr)
	boolMixedStr = "TrUe" // Tests case-insensitivity overhead
	boolMixed    = []byte(boolMixedStr)

	// Ints (Base 10)
	intSmallStr      = "123"
	intSmall         = []byte(intSmallStr)
	intLargeStr      = "1234567890123456789"
	intLarge         = []byte(intLargeStr)
	intNegativeStr   = "-1234567890123456789"
	intNegative      = []byte(intNegativeStr)
	intUnderscoreStr = "1_234_567_890" // Tests underscore bypass
	intUnderscore    = []byte(intUnderscoreStr)

	// Ints (Other Bases)
	intHexLowerStr = "0x1234abcd" // Tests normal hex
	intHexLower    = []byte(intHexLowerStr)
	intHexUpperStr = "0x1234ABCD" // Tests casing jump
	intHexUpper    = []byte(intHexUpperStr)
	intOctStr      = "0o1234567"
	intOct         = []byte(intOctStr)
	intBinStr      = "0b101010101010"
	intBin         = []byte(intBinStr)
	intBase36Str   = "zzyyxxwwvv" // Tests arbitrary base lookup table
	intBase36      = []byte(intBase36Str)

	// Floats
	floatSmallStr = "3.14" // standard fast path
	floatSmall    = []byte(floatSmallStr)
	floatLargeStr = "1.7976931348623157e+308" // max float bounds
	floatLarge    = []byte(floatLargeStr)
	floatSciStr   = "1.2345e-67"
	floatSci      = []byte(floatSciStr)
	floatHexStr   = "0x1.b7p-1" // hex float fast path
	floatHex      = []byte(floatHexStr)
	// floatSlowPath forces Eisel-Lemire to fail, dropping into decimal.go
	floatSlowPathStr = "3.14159265358979323846264338327950288419716939937510582097494459"
	floatSlowPath    = []byte(floatSlowPathStr)
	floatNaNStr      = "NaN"
	floatNaN         = []byte(floatNaNStr)
	floatInfStr      = "+Inf"
	floatInf         = []byte(floatInfStr)

	// Complex
	complexStandardStr = "(1.234+5.678i)"
	complexStandard    = []byte(complexStandardStr)
	complexNoParensStr = "1.234+5.678i"
	complexNoParens    = []byte(complexNoParensStr)
	complexImagStr     = "5.678i"
	complexImag        = []byte(complexImagStr)
)

// ==========================================
// PARSE BOOL
// ==========================================

func BenchmarkParseBool_Standard_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseBool(boolTrueStr)
		strconv.ParseBool(boolFalseStr)
	}
}
func BenchmarkParseBool_Standard_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseBool(boolTrue)
		byteconv.ParseBool(boolFalse)
	}
}

func BenchmarkParseBool_MixedCase_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseBool(boolMixedStr)
	}
}
func BenchmarkParseBool_MixedCase_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseBool(boolMixed)
	}
}

// ==========================================
// ATOI & PARSE INT (Base 10)
// ==========================================

func BenchmarkAtoi_Small_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.Atoi(intSmallStr)
	}
}
func BenchmarkAtoi_Small_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.Atoi(intSmall)
	}
}

func BenchmarkAtoi_Large_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.Atoi(intLargeStr)
	}
}
func BenchmarkAtoi_Large_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.Atoi(intLarge)
	}
}

func BenchmarkAtoi_Negative_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.Atoi(intNegativeStr)
	}
}
func BenchmarkAtoi_Negative_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.Atoi(intNegative)
	}
}

func BenchmarkAtoi_Underscore_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.Atoi(intUnderscoreStr)
	}
}
func BenchmarkAtoi_Underscore_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.Atoi(intUnderscore)
	}
}

// ==========================================
// PARSE INT / UINT (Mixed Bases)
// ==========================================

func BenchmarkParseInt_HexLower_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseInt(intHexLowerStr, 0, 64)
	}
}
func BenchmarkParseInt_HexLower_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseInt(intHexLower, 0, 64)
	}
}

func BenchmarkParseInt_HexUpper_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseInt(intHexUpperStr, 0, 64)
	}
}
func BenchmarkParseInt_HexUpper_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseInt(intHexUpper, 0, 64)
	}
}

func BenchmarkParseInt_Octal_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseInt(intOctStr, 0, 64)
	}
}
func BenchmarkParseInt_Octal_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseInt(intOct, 0, 64)
	}
}

func BenchmarkParseInt_Binary_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseInt(intBinStr, 0, 64)
	}
}
func BenchmarkParseInt_Binary_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseInt(intBin, 0, 64)
	}
}

func BenchmarkParseInt_Base36_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseInt(intBase36Str, 36, 64)
	}
}
func BenchmarkParseInt_Base36_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseInt(intBase36, 36, 64)
	}
}

func BenchmarkParseUint_Large_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseUint(intLargeStr, 10, 64)
	}
}
func BenchmarkParseUint_Large_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseUint(intLarge, 10, 64)
	}
}

// ==========================================
// PARSE FLOAT
// ==========================================

func BenchmarkParseFloat_Small_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseFloat(floatSmallStr, 64)
	}
}
func BenchmarkParseFloat_Small_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseFloat(floatSmall, 64)
	}
}

func BenchmarkParseFloat_Large_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseFloat(floatLargeStr, 64)
	}
}
func BenchmarkParseFloat_Large_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseFloat(floatLarge, 64)
	}
}

func BenchmarkParseFloat_Sci_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseFloat(floatSciStr, 64)
	}
}
func BenchmarkParseFloat_Sci_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseFloat(floatSci, 64)
	}
}

func BenchmarkParseFloat_Hex_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseFloat(floatHexStr, 64)
	}
}
func BenchmarkParseFloat_Hex_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseFloat(floatHex, 64)
	}
}

func BenchmarkParseFloat_SlowPath_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseFloat(floatSlowPathStr, 64)
	}
}
func BenchmarkParseFloat_SlowPath_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseFloat(floatSlowPath, 64)
	}
}

func BenchmarkParseFloat_NaN_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseFloat(floatNaNStr, 64)
	}
}
func BenchmarkParseFloat_NaN_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseFloat(floatNaN, 64)
	}
}

func BenchmarkParseFloat_Inf_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseFloat(floatInfStr, 64)
	}
}
func BenchmarkParseFloat_Inf_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseFloat(floatInf, 64)
	}
}

// ==========================================
// PARSE COMPLEX
// ==========================================

func BenchmarkParseComplex_Standard_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseComplex(complexStandardStr, 128)
	}
}
func BenchmarkParseComplex_Standard_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseComplex(complexStandard, 128)
	}
}

func BenchmarkParseComplex_NoParens_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseComplex(complexNoParensStr, 128)
	}
}
func BenchmarkParseComplex_NoParens_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseComplex(complexNoParens, 128)
	}
}

func BenchmarkParseComplex_Imag_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseComplex(complexImagStr, 128)
	}
}
func BenchmarkParseComplex_Imag_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseComplex(complexImag, 128)
	}
}

// ==========================================
// REAL-WORLD OVERHEAD (FromBytes)
// Tests the latency of forcing an allocation
// `string(bytes)` to use the standard library.
// ==========================================

func BenchmarkAtoi_FromBytes_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.Atoi(string(intLarge))
	}
}
func BenchmarkAtoi_FromBytes_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.Atoi(intLarge)
	}
}

func BenchmarkParseFloat_FromBytes_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseFloat(string(floatLarge), 64)
	}
}
func BenchmarkParseFloat_FromBytes_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseFloat(floatLarge, 64)
	}
}

func BenchmarkParseComplex_FromBytes_Strconv(b *testing.B) {
	for b.Loop() {
		strconv.ParseComplex(string(complexStandard), 128)
	}
}
func BenchmarkParseComplex_FromBytes_Byteconv(b *testing.B) {
	for b.Loop() {
		byteconv.ParseComplex(complexStandard, 128)
	}
}
