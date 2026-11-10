// Package figurate provides polygonal (figurate) numbers: triangular,
// square, pentagonal, hexagonal and general s-gonal numbers, exact tests
// for them, and streams of them.
//
// The n-th s-gonal number is P(s, n) = ((s-2)n² - (s-4)n) / 2, so the
// triangular numbers are P(3, n), the squares P(4, n), and so on. The Is*
// tests count 0 as P(s, 0).
package figurate

import (
	"math/big"

	"github.com/t1sk3/eulerlib/v4/internal/mathx"
	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/constraints"
)

// NthTriangular returns n(n+1)/2.
func NthTriangular[E constraints.Integer](n E) E {
	if n%2 == 0 {
		return n / 2 * (n + 1)
	}
	return (n + 1) / 2 * n
}

// NthSquare returns n².
func NthSquare[E constraints.Integer](n E) E {
	return n * n
}

// NthPentagonal returns n(3n-1)/2.
func NthPentagonal[E constraints.Integer](n E) E {
	return NthPolygonal(5, n)
}

// NthHexagonal returns n(2n-1).
func NthHexagonal[E constraints.Integer](n E) E {
	return n * (2*n - 1)
}

// NthPolygonal returns the n-th s-gonal number ((s-2)n² - (s-4)n)/2
// (Project Euler #61). It panics if s < 3.
func NthPolygonal[E constraints.Integer](s, n E) E {
	if s < 3 {
		panic("figurate: polygonal numbers need s >= 3")
	}
	// n((s-2)n - (s-4)) is always even; halve whichever factor is even so
	// the intermediate value is no larger than the result.
	f := (s-2)*n - (s - 4)
	if n%2 == 0 {
		return n / 2 * f
	}
	return n * (f / 2)
}

// IsTriangular reports whether n = k(k+1)/2 for some k >= 0, i.e. whether
// 8n+1 is a perfect square. It is exact for every value of E.
func IsTriangular[E constraints.Integer](n E) bool {
	r, ok := discRoot(n, 8)
	return ok && r%2 == 1
}

// IsPentagonal reports whether n = k(3k-1)/2 for some k >= 0, i.e. whether
// 24n+1 is a perfect square r² with r ≡ 5 (mod 6) (or n = 0). It is exact
// for every value of E.
func IsPentagonal[E constraints.Integer](n E) bool {
	if n == 0 {
		return true
	}
	r, ok := discRoot(n, 24)
	return ok && r%6 == 5
}

// IsHexagonal reports whether n = k(2k-1) for some k >= 0, i.e. whether
// 8n+1 is a perfect square r² with r ≡ 3 (mod 4) (or n = 0). It is exact
// for every value of E.
func IsHexagonal[E constraints.Integer](n E) bool {
	if n == 0 {
		return true
	}
	r, ok := discRoot(n, 8)
	return ok && r%4 == 3
}

// discRoot returns sqrt(k*n+1) and true if k*n+1 is a perfect square, or
// false otherwise (also for n < 0). It falls back to math/big when k*n+1
// does not fit in a uint64.
func discRoot[E constraints.Integer](n E, k uint64) (uint64, bool) {
	if n < 0 {
		return 0, false
	}
	u := uint64(n)
	if u <= (^uint64(0)-1)/k {
		d := k*u + 1
		r := mathx.Isqrt(d)
		return r, r*r == d
	}
	d := new(big.Int).SetUint64(u)
	d.Mul(d, new(big.Int).SetUint64(k)).Add(d, big.NewInt(1))
	r := new(big.Int).Sqrt(d)
	return r.Uint64(), new(big.Int).Mul(r, r).Cmp(d) == 0 // r < 2^34, fits
}

// Triangulars returns 1, 3, 6, 10, ..., ending at the largest triangular
// number that fits in E.
func Triangulars[E constraints.Integer]() stream.Stream[E] { return Polygonals(E(3)) }

// Squares returns 1, 4, 9, 16, ..., ending at the largest square that fits
// in E.
func Squares[E constraints.Integer]() stream.Stream[E] { return Polygonals(E(4)) }

// Pentagonals returns 1, 5, 12, 22, ..., ending at the largest pentagonal
// number that fits in E.
func Pentagonals[E constraints.Integer]() stream.Stream[E] { return Polygonals(E(5)) }

// Hexagonals returns 1, 6, 15, 28, ..., ending at the largest hexagonal
// number that fits in E.
func Hexagonals[E constraints.Integer]() stream.Stream[E] { return Polygonals(E(6)) }

// Polygonals returns the s-gonal numbers P(s, 1), P(s, 2), ... in order,
// ending at the largest one that fits in E. It panics if s < 3.
func Polygonals[E constraints.Integer](s E) stream.Stream[E] {
	if s < 3 {
		panic("figurate: polygonal numbers need s >= 3")
	}
	return stream.FromSeq(func(yield func(E) bool) {
		max := mathx.MaxOf[E]()
		// P(s, n+1) = P(s, n) + (s-2)n + 1.
		v, n := E(1), E(1)
		for {
			if !yield(v) {
				return
			}
			if n > (max-1)/(s-2) {
				return
			}
			d := (s-2)*n + 1
			if v > max-d {
				return
			}
			v += d
			n++
		}
	})
}
