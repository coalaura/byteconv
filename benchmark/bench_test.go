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
