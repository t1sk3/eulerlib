// Package fibonacci provides the Fibonacci numbers F(0) = 0, F(1) = 1,
// F(n) = F(n-1) + F(n-2), as streams and as single terms, with
// arbitrary-precision variants backed by math/big.
package fibonacci

import (
	"math/big"

	"github.com/t1sk3/eulerlib/v4/internal/mathx"
	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/constraints"
)

// All returns the Fibonacci numbers 0, 1, 1, 2, 3, 5, ... that fit in E,
// and ends at the largest one: All[int64]() has 93 elements, F(0) to
// F(92), and All[uint8]() ends at 233. Bound it with Limit or TakeWhile:
//
//	fibonacci.All[int]().TakeWhile(stream.AtMost(4_000_000))
//
// The stream can be consumed any number of times.
func All[E constraints.Integer]() stream.Stream[E] {
	return stream.FromSeq(func(yield func(E) bool) {
		max := mathx.MaxOf[E]()
		a, b := E(0), E(1)
		for {
			if !yield(a) {
				return
			}
			if a > max-b { // a+b would overflow: b is the last one that fits
				yield(b)
				return
			}
			a, b = b, a+b
		}
	})
}

// AllBig returns every Fibonacci number, 0, 1, 1, 2, ..., as a *big.Int.
// The stream is infinite, so bound it with Limit or TakeWhile. Each element
// is a new *big.Int that the caller may keep or modify.
func AllBig() stream.Stream[*big.Int] {
	return stream.FromSeq(func(yield func(*big.Int) bool) {
		a, b := big.NewInt(0), big.NewInt(1)
		for {
			if !yield(new(big.Int).Set(a)) {
				return
			}
			a.Add(a, b)
			a, b = b, a
		}
	})
}

// First returns F(0), ..., F(n-1). If some of those do not fit in E, it
// returns only the ones that do.
func First[E constraints.Integer](n int) []E {
	if n <= 0 {
		return []E{}
	}
	return All[E]().Limit(n).ToSlice()
}

// Nth returns F(n). It panics if n < 0 or F(n) does not fit in E (use
// NthBig for those).
func Nth[E constraints.Integer](n int) E {
	if n < 0 {
		panic("fibonacci: Nth with negative n")
	}
	v, ok := All[E]().Skip(n).FindFirst().Value()
	if !ok {
		panic("fibonacci: Nth overflows the element type")
	}
	return v
}

// NthBig returns F(n) as a *big.Int, using fast doubling in O(log n)
// multiplications. It panics if n < 0.
func NthBig(n int64) *big.Int {
	if n < 0 {
		panic("fibonacci: NthBig with negative n")
	}
	// Invariant: a = F(k), b = F(k+1), walking k through the bits of n.
	a, b := big.NewInt(0), big.NewInt(1)
	t1, t2 := new(big.Int), new(big.Int)
	for i := 62; i >= 0; i-- {
		// F(2k) = F(k) * (2F(k+1) - F(k)); F(2k+1) = F(k)^2 + F(k+1)^2
		t1.Lsh(b, 1).Sub(t1, a).Mul(t1, a)
		t2.Mul(a, a)
		b.Mul(b, b).Add(b, t2)
		a.Set(t1)
		if n>>i&1 == 1 {
			a.Add(a, b)
			a, b = b, a
		}
	}
	return a
}

// IndexOfFirstWithDigits returns the smallest n >= 1 such that F(n) has at
// least d decimal digits (Project Euler #25). It panics if d < 1.
func IndexOfFirstWithDigits(d int) int {
	if d < 1 {
		panic("fibonacci: IndexOfFirstWithDigits needs d >= 1")
	}
	threshold := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(d-1)), nil)
	a, b := big.NewInt(1), big.NewInt(1) // F(1), F(2)
	n := 1
	for a.Cmp(threshold) < 0 {
		a.Add(a, b)
		a, b = b, a
		n++
	}
	return n
}
