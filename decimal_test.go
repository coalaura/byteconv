// Copyright 2026 coalaura. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package byteconv

import (
	"testing"
)

type roundIntTest struct {
	i     uint64
	shift int
	int   uint64
}

var roundinttests = []roundIntTest{
	{0, 100, 0},
	{512, -8, 2},
	{513, -8, 2},
	{640, -8, 2},
	{641, -8, 3},
	{384, -8, 2},
	{385, -8, 2},
	{383, -8, 1},
	{1, 100, 1<<64 - 1},
	{1000, 0, 1000},
}

func TestDecimalRoundedInteger(t *testing.T) {
	for i := 0; i < len(roundinttests); i++ {
		test := roundinttests[i]
		var d decimal
		d.Assign(test.i)
		d.Shift(test.shift)
		intVal := d.RoundedInteger()
		if intVal != test.int {
			t.Errorf("Decimal %v >> %v RoundedInteger = %v, want %v",
				test.i, test.shift, intVal, test.int)
		}
	}
}
