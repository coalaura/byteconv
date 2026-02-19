// Copyright 2026 coalaura. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package byteconv implements conversions to and from byte representations
// of basic data types. It aims to be a zero-allocation, highly performant
// 1-to-1 equivalent of the standard library's strconv package, but operates
// directly on byte slices ([]byte) to avoid string conversions and memory allocations.
package byteconv
