// Package figurate provides functions for generating and testing triangular,
// pentagonal, and hexagonal numbers.
package figurate

import (
	"math"

	"github.com/t1sk3/eulerlib/v2/utils"
)

// NthTriangular returns the nth triangular number: n*(n+1)/2.
func NthTriangular[E utils.Integer](n E) E {
	return n * (n + 1) / 2
}

// IsTriangular returns true if n is a triangular number.
// A number t is triangular iff 8t+1 is a perfect square.
func IsTriangular[E utils.Integer](n E) bool {
	if n < 0 {
		return false
	}
	n64 := int64(n)
	if n64 < 0 { // large uint64 that overflowed int64
		return false
	}
	disc := 8*n64 + 1
	sq := int64(math.Sqrt(float64(disc)))
	return sq*sq == disc
}

// NthPentagonal returns the nth pentagonal number: n*(3n-1)/2.
func NthPentagonal[E utils.Integer](n E) E {
	return n * (3*n - 1) / 2
}

// IsPentagonal returns true if n is a pentagonal number.
// A number p is pentagonal iff (1+sqrt(1+24p))/6 is a positive integer.
func IsPentagonal[E utils.Integer](n E) bool {
	if n < 0 {
		return false
	}
	n64 := int64(n)
	if n64 < 0 { // large uint64 that overflowed int64
		return false
	}
	disc := 1 + 24*n64
	sq := int64(math.Sqrt(float64(disc)))
	if sq*sq != disc {
		return false
	}
	return (1+sq)%6 == 0
}

// NthHexagonal returns the nth hexagonal number: n*(2n-1).
func NthHexagonal[E utils.Integer](n E) E {
	return n * (2*n - 1)
}

// IsHexagonal returns true if n is a hexagonal number.
// A number h is hexagonal iff (1+sqrt(1+8h))/4 is a positive integer.
func IsHexagonal[E utils.Integer](n E) bool {
	if n < 0 {
		return false
	}
	n64 := int64(n)
	if n64 < 0 { // large uint64 that overflowed int64
		return false
	}
	disc := 1 + 8*n64
	sq := int64(math.Sqrt(float64(disc)))
	if sq*sq != disc {
		return false
	}
	return (1+sq)%4 == 0
}
