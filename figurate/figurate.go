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
	"math/bits"

	"github.com/t1sk3/eulerlib/v4/internal/mathx"
	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/constraints"
	"github.com/t1sk3/go-utils/types/optional"
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

// polygonal returns P(s, k) for s >= 3, and whether it fits in a uint64.
func polygonal(s, k uint64) (uint64, bool) {
	if k == 0 {
		return 0, true
	}
	// P(s, k) = k·f/2 with f = (s-2)(k-1) + 2; halve whichever is even.
	hi, f := bits.Mul64(s-2, k-1)
	if hi != 0 || f > ^uint64(0)-2 {
		return 0, false
	}
	f += 2
	if k%2 == 0 {
		hi, f = bits.Mul64(k/2, f)
	} else {
		hi, f = bits.Mul64(k, f/2)
	}
	return f, hi == 0
}

// PolygonalIndex returns the k >= 0 with NthPolygonal(s, k) == n, or an
// empty Optional if n is not an s-gonal number (or is negative):
// PolygonalIndex(3, 10) is 4, because 10 is the fourth triangular number.
// The index of 0 is 0 and the index of 1 is 1, for every s. It is exact
// for every value of E, and panics if s < 3.
func PolygonalIndex[E constraints.Integer](s, n E) optional.Optional[E] {
	if s < 3 {
		panic("figurate: polygonal numbers need s >= 3")
	}
	if n < 0 {
		return optional.Empty[E]()
	}
	us, un := uint64(s), uint64(n)
	// P(s, k) is about (s-2)k²/2, so k is within a step or two of
	// sqrt(2n/(s-2)); try the candidates around it.
	var est uint64
	if hi, lo := bits.Mul64(un, 2); hi < us-2 {
		q, _ := bits.Div64(hi, lo, us-2)
		est = mathx.Isqrt(q)
	} else { // 2n/(s-2) needs more than 64 bits: s = 3 and n >= 2^63
		d := new(big.Int).SetUint64(un)
		est = d.Sqrt(d.Lsh(d, 1)).Uint64()
	}
	for k := max(est, 2) - 2; k <= est+2; k++ {
		if p, ok := polygonal(us, k); ok && p == un {
			return optional.Of(E(k))
		}
	}
	return optional.Empty[E]()
}

// IsPolygonal reports whether n = P(s, k) for some k >= 0: the general
// form of IsTriangular, IsPentagonal and IsHexagonal (Project Euler #61).
// It is exact for every value of E, and panics if s < 3.
func IsPolygonal[E constraints.Integer](s, n E) bool {
	return PolygonalIndex(s, n).IsPresent()
}
