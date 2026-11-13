package numtheory

import (
	"iter"
	"math/big"

	"github.com/t1sk3/eulerlib/v4/internal/mathx"
	"github.com/t1sk3/go-utils/types/constraints"
	"github.com/t1sk3/go-utils/types/optional"
	"github.com/t1sk3/go-utils/types/pair"
)

// SqrtContinuedFraction returns the continued fraction of √n as its
// integer part a0 and the block of terms that repeats forever after it:
// √23 = [4; (1, 3, 1, 8)] gives (4, [1 3 1 8]) (Project Euler #64). The
// period is empty when n is a perfect square. It panics if n < 0.
func SqrtContinuedFraction[E constraints.Integer](n E) (a0 E, period []E) {
	if n < 0 {
		panic("numtheory: SqrtContinuedFraction of a negative number")
	}
	u := uint64(n)
	r := mathx.Isqrt(u)
	period = []E{}
	if r*r == u {
		return E(r), period
	}
	// Each complete quotient is (√n + m) / d, with 0 <= m <= r and
	// 0 < d <= 2r, so nothing overflows. The period ends at the term 2·a0.
	m, d, a := uint64(0), uint64(1), r
	for a != 2*r {
		m = d*a - m
		d = (u - m*m) / d
		a = (r + m) / d
		period = append(period, E(a))
	}
	return E(r), period
}

// SqrtCFTerms returns the terms of the continued fraction of √n as an
// endless sequence: a0, then the period over and over (just a0, once, if n
// is a perfect square). Pass it to Convergents. It panics if n < 0.
func SqrtCFTerms[E constraints.Integer](n E) iter.Seq[E] {
	a0, period := SqrtContinuedFraction(n)
	return func(yield func(E) bool) {
		if !yield(a0) {
			return
		}
		for len(period) > 0 {
			for _, a := range period {
				if !yield(a) {
					return
				}
			}
		}
	}
}

// ECFTerms returns the terms of the continued fraction of e, endlessly:
// 2, 1, 2, 1, 1, 4, 1, 1, 6, 1, ... (Project Euler #65).
func ECFTerms() iter.Seq[int] {
	return func(yield func(int) bool) {
		if !yield(2) {
			return
		}
		for k := 1; ; k++ {
			if !yield(1) || !yield(2*k) || !yield(1) {
				return
			}
		}
	}
}

// Convergents returns the convergents h/k of the continued fraction with
// the given terms, as (numerator, denominator) pairs in lowest terms, one
// per term:
//
//	for h, k := range numtheory.Convergents(numtheory.SqrtCFTerms(2)) {
//	    // 1/1, 3/2, 7/5, 17/12, ...  (Project Euler #57)
//	}
//
// Every yielded *big.Int is new, so the caller may keep or modify it. The
// sequence is as long as terms, so bound an endless one by breaking out of
// the loop.
func Convergents[E constraints.Integer](terms iter.Seq[E]) iter.Seq2[*big.Int, *big.Int] {
	return func(yield func(*big.Int, *big.Int) bool) {
		// h(-1) = 1, h(-2) = 0; k(-1) = 0, k(-2) = 1.
		h1, h2 := big.NewInt(1), big.NewInt(0)
		k1, k2 := big.NewInt(0), big.NewInt(1)
		a := new(big.Int)
		for t := range terms {
			if mathx.IsSigned[E]() {
				a.SetInt64(int64(t))
			} else {
				a.SetUint64(uint64(t))
			}
			// h = a*h1 + h2, reusing h2's storage for the new value.
			h2.Add(h2, new(big.Int).Mul(a, h1))
			k2.Add(k2, new(big.Int).Mul(a, k1))
			h1, h2 = h2, h1
			k1, k2 = k2, k1
			if !yield(new(big.Int).Set(h1), new(big.Int).Set(k1)) {
				return
			}
		}
	}
}

// Pell returns the fundamental solution of Pell's equation x² - d·y² = 1:
// the solution (x, y) in positive integers with the smallest x (Project
// Euler #66). It returns an empty Optional if there is none, which is when
// d <= 0 or d is a perfect square.
//
// The solution is found among the convergents of √d. Every other solution
// (x', y') follows from x' + y'√d = (x + y√d)^k.
func Pell[E constraints.Integer](d E) optional.Optional[pair.Pair[*big.Int, *big.Int]] {
	if d <= 0 || IsSquare(d) {
		return optional.Empty[pair.Pair[*big.Int, *big.Int]]()
	}
	bd := new(big.Int).SetUint64(uint64(d))
	one := big.NewInt(1)
	lhs, t := new(big.Int), new(big.Int)
	for x, y := range Convergents(SqrtCFTerms(d)) {
		lhs.Mul(x, x)
		lhs.Sub(lhs, t.Mul(y, y).Mul(t, bd))
		if lhs.Cmp(one) == 0 {
			return optional.Of(pair.Of(x, y))
		}
	}
	panic("unreachable: √d has an infinite continued fraction")
}
