// Copyright 2026 coalaura. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package byteconv

// SetOptimize makes the unexported "optimize" flag available to tests.
func SetOptimize(opt bool) bool {
	old := optimize
	optimize = opt
	return old
}

// ParseFloatPrefix makes the unexported "parseFloatPrefix" available to tests.
func ParseFloatPrefix(s []byte, bitSize int) (float64, int, error) {
	return parseFloatPrefix(s, bitSize)
}
