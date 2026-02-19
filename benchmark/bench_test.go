package byteconv_bench

import (
	"strconv"
	"testing"
	"github.com/coalaura/byteconv"
)

var (
	boolTrueStr  = "true"
	boolTrue     = []byte(boolTrueStr)
	boolFalseStr = "false"
	boolFalse    = []byte(boolFalseStr)

	intSmallStr = "123"
	intSmall    = []byte(intSmallStr)
	intLargeStr = "1234567890123456789"
	intLarge    = []byte(intLargeStr)

	floatSmallStr = "3.14"
	floatSmall    = []byte(floatSmallStr)
	floatLargeStr = "1.7976931348623157e+308"
	floatLarge    = []byte(floatLargeStr)
)

func BenchmarkParseBool_Strconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseBool(boolTrueStr)
		strconv.ParseBool(boolFalseStr)
	}
}

func BenchmarkParseBool_Byteconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		byteconv.ParseBool(boolTrue)
		byteconv.ParseBool(boolFalse)
	}
}

func BenchmarkAtoi_Small_Strconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.Atoi(intSmallStr)
	}
}

func BenchmarkAtoi_Small_Byteconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		byteconv.Atoi(intSmall)
	}
}

func BenchmarkAtoi_Large_Strconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.Atoi(intLargeStr)
	}
}

func BenchmarkAtoi_Large_Byteconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		byteconv.Atoi(intLarge)
	}
}

func BenchmarkParseFloat_Small_Strconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseFloat(floatSmallStr, 64)
	}
}

func BenchmarkParseFloat_Small_Byteconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		byteconv.ParseFloat(floatSmall, 64)
	}
}

func BenchmarkParseFloat_Large_Strconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseFloat(floatLargeStr, 64)
	}
}

func BenchmarkParseFloat_Large_Byteconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		byteconv.ParseFloat(floatLarge, 64)
	}
}

// To accurately compare string casting overhead vs byteconv
func BenchmarkParseFloat_FromBytes_Strconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseFloat(string(floatLarge), 64)
	}
}

func BenchmarkParseFloat_FromBytes_Byteconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		byteconv.ParseFloat(floatLarge, 64)
	}
}

var (
	intHexStr = "0x1234abcd"
	intHex    = []byte(intHexStr)
	intOctStr = "0o1234567"
	intOct    = []byte(intOctStr)
	intBinStr = "0b101010101010"
	intBin    = []byte(intBinStr)

	floatSciStr = "1.2345e-67"
	floatSci    = []byte(floatSciStr)
	floatNaNStr = "NaN"
	floatNaN    = []byte(floatNaNStr)
	floatInfStr = "+Inf"
	floatInf    = []byte(floatInfStr)
	
	complexStr = "(1.234+5.678i)"
	complexBytes = []byte(complexStr)
)

func BenchmarkParseInt_Hex_Strconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseInt(intHexStr, 0, 64)
	}
}

func BenchmarkParseInt_Hex_Byteconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		byteconv.ParseInt(intHex, 0, 64)
	}
}

func BenchmarkParseInt_Oct_Strconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseInt(intOctStr, 0, 64)
	}
}

func BenchmarkParseInt_Oct_Byteconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		byteconv.ParseInt(intOct, 0, 64)
	}
}

func BenchmarkParseInt_Bin_Strconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseInt(intBinStr, 0, 64)
	}
}

func BenchmarkParseInt_Bin_Byteconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		byteconv.ParseInt(intBin, 0, 64)
	}
}

func BenchmarkParseFloat_Sci_Strconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseFloat(floatSciStr, 64)
	}
}

func BenchmarkParseFloat_Sci_Byteconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		byteconv.ParseFloat(floatSci, 64)
	}
}

func BenchmarkParseFloat_NaN_Strconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseFloat(floatNaNStr, 64)
	}
}

func BenchmarkParseFloat_NaN_Byteconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		byteconv.ParseFloat(floatNaN, 64)
	}
}

func BenchmarkParseFloat_Inf_Strconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseFloat(floatInfStr, 64)
	}
}

func BenchmarkParseFloat_Inf_Byteconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		byteconv.ParseFloat(floatInf, 64)
	}
}

func BenchmarkParseComplex_Strconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseComplex(complexStr, 128)
	}
}

func BenchmarkParseComplex_Byteconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		byteconv.ParseComplex(complexBytes, 128)
	}
}
