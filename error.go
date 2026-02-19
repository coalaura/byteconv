// Copyright 2026 coalaura. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package byteconv

// Error represents an error during a conversion.
type Error int

const (
	_ Error = iota
	// ErrRange indicates that a value is out of range for the target type.
	ErrRange
	// ErrSyntax indicates that a value does not have the right syntax for the target type.
	ErrSyntax
	// ErrBase indicates that a base is invalid.
	ErrBase
	// ErrBitSize indicates that a bit size is invalid.
	ErrBitSize
)

func (e Error) Error() string {
	switch e {
	case ErrRange:
		return "value out of range"
	case ErrSyntax:
		return "invalid syntax"
	case ErrBase:
		return "invalid base"
	case ErrBitSize:
		return "invalid bit size"
	}

	return "unknown error"
}
