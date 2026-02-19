// Copyright 2026 coalaura. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package byteconv_test

import (
	"testing"

	. "github.com/coalaura/byteconv"
)

type atobTest struct {
	in  string
	out bool
	err error
}

var atobtests = []atobTest{
	{"", false, ErrSyntax},
	{"asdf", false, ErrSyntax},
	{"0", false, nil},
	{"f", false, nil},
	{"F", false, nil},
	{"FALSE", false, nil},
	{"false", false, nil},
	{"False", false, nil},
	{"1", true, nil},
	{"t", true, nil},
	{"T", true, nil},
	{"TRUE", true, nil},
	{"true", true, nil},
	{"True", true, nil},
}

func TestParseBool(t *testing.T) {
	for _, test := range atobtests {
		b, e := ParseBool([]byte(test.in))
		if b != test.out || e != test.err {
			t.Errorf("ParseBool(%s) = %v, %v, want %v, %v", test.in, b, e, test.out, test.err)
		}
	}
}
