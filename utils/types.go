// Package utils provides the generic type constraints and reflection-based
// type checks shared by the rest of eulerlib.
package utils

// SignedInteger is satisfied by any signed integer type.
type SignedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// UnsignedInteger is satisfied by any unsigned integer type.
type UnsignedInteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is satisfied by any signed or unsigned integer type.
type Integer interface {
	SignedInteger | UnsignedInteger
}

// Float is satisfied by any floating-point type.
type Float interface {
	~float32 | ~float64
}

// RealNumber is satisfied by any integer or floating-point type.
type RealNumber interface {
	Integer | Float
}

// Complex is satisfied by any complex number type.
type Complex interface {
	~complex64 | ~complex128
}

// Number is satisfied by any real or complex number type.
type Number interface {
	RealNumber | Complex
}

// Comparable is satisfied by any number, string, or bool type.
type Comparable interface {
	Number | ~string | ~bool
}

// SignedNumber is satisfied by any signed integer, float, or complex type.
type SignedNumber interface {
	SignedInteger | Float | Complex
}
