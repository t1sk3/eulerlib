// Package mathx holds exact integer helpers shared by eulerlib's packages:
// integer square roots, overflow-free modular multiplication and
// exponentiation, and the limits of a generic integer type.
package mathx

import (
	"math"
	"math/bits"
	"unsafe"

	"github.com/t1sk3/go-utils/types/constraints"
)

// Isqrt returns floor(sqrt(n)), exactly, for every uint64.
//
// math.Sqrt on a float64 is only exact up to 2^53, so the float estimate is
// corrected by at most a few steps in each direction.
func Isqrt(n uint64) uint64 {
	r := min(uint64(math.Sqrt(float64(n))), math.MaxUint32)
	for r*r > n {
		r--
	}
	for r < math.MaxUint32 && (r+1)*(r+1) <= n {
		r++
	}
	return r
}

// MulMod returns (a*b) mod m without overflow, using a 128-bit product.
// It panics if m is 0.
func MulMod(a, b, m uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	return bits.Rem64(hi, lo, m)
}

// PowMod returns b^e mod m by binary exponentiation. Every intermediate
// product goes through MulMod, so the result is exact for any m that fits
// in a uint64. It panics if m is 0.
func PowMod(b, e, m uint64) uint64 {
	r := uint64(1) % m
	b %= m
	for e > 0 {
		if e&1 == 1 {
			r = MulMod(r, b, m)
		}
		b = MulMod(b, b, m)
		e >>= 1
	}
	return r
}

// IsSigned reports whether E is a signed integer type.
func IsSigned[E constraints.Integer]() bool {
	return ^E(0) < 0
}

// MaxOf returns the largest value of the integer type E.
func MaxOf[E constraints.Integer]() E {
	ones := ^E(0)
	if ones > 0 { // unsigned: all bits set
		return ones
	}
	return E(^uint64(0) >> (65 - unsafe.Sizeof(ones)*8))
}

// MinOf returns the smallest value of the integer type E: 0 for unsigned
// types, -MaxOf[E]()-1 for signed ones.
func MinOf[E constraints.Integer]() E {
	if !IsSigned[E]() {
		return 0
	}
	return -MaxOf[E]() - 1
}

// Abs64 returns |n| as a uint64. It is exact for every value of E,
// including the minimum of a signed type (whose negation would overflow E).
func Abs64[E constraints.Integer](n E) uint64 {
	if n < 0 {
		return -uint64(int64(n))
	}
	return uint64(n)
}
