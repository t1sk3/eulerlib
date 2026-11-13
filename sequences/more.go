package sequences

import (
	"github.com/t1sk3/eulerlib/v4/internal/mathx"
	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/constraints"
	"github.com/t1sk3/go-utils/types/pair"
)

// Farey returns the Farey sequence of order n: every fraction a/b in
// lowest terms with 0 <= a/b <= 1 and b <= n, in ascending order, as
// (numerator, denominator) pairs from (0, 1) to (1, 1). Farey(3) yields
// 0/1, 1/3, 1/2, 2/3, 1/1 (Project Euler #71–73). Each term follows from
// the two before it, so nothing is sorted or stored. It is empty for
// n < 1.
func Farey[E constraints.Integer](n E) stream.Stream[pair.Pair[E, E]] {
	return stream.FromSeq(func(yield func(pair.Pair[E, E]) bool) {
		if n < 1 {
			return
		}
		// uint64 keeps the intermediate values (up to 2n) from
		// overflowing a narrow E.
		N := uint64(n)
		a, b, c, d := uint64(0), uint64(1), uint64(1), N
		if !yield(pair.Of(E(a), E(b))) {
			return
		}
		for c <= N {
			k := (N + b) / d
			a, b, c, d = c, d, k*c-a, k*d-b
			if !yield(pair.Of(E(a), E(b))) {
				return
			}
		}
	})
}

// IsHappy reports whether n is a happy number: replacing it by the sum of
// the squares of its digits, over and over, reaches 1 (7 → 49 → 97 → 130
// → 10 → 1). Every other positive number falls into the cycle through 4,
// which also contains 89 (Project Euler #92). It is false for n <= 0.
func IsHappy[E constraints.Integer](n E) bool {
	if n <= 0 {
		return false
	}
	u := uint64(n)
	for u != 1 && u != 4 {
		var s uint64
		for ; u > 0; u /= 10 {
			s += (u % 10) * (u % 10)
		}
		u = s
	}
	return u == 1
}

// ChainLength returns the number of distinct terms in the chain start,
// next(start), next(next(start)), ... before a term repeats (Project Euler
// #74: 69 → 363600 → 1454 → 169 → 363601 → 1454 has five). It remembers
// every term, and does not return if the chain never repeats.
func ChainLength[T comparable](start T, next func(T) T) int {
	seen := map[T]struct{}{}
	for v := start; ; v = next(v) {
		if _, ok := seen[v]; ok {
			return len(seen)
		}
		seen[v] = struct{}{}
	}
}

// LinearRecurrenceMod returns a(n) mod m for the sequence defined by its
// first k terms, initial = a(0), ..., a(k-1), and the recurrence
//
//	a(i) = coeffs[0]·a(i-1) + coeffs[1]·a(i-2) + ... + coeffs[k-1]·a(i-k)
//
// It raises the k×k companion matrix to the n-th power by repeated
// squaring, O(k³ log n), so n can be astronomically large: the Fibonacci
// numbers modulo m are LinearRecurrenceMod([]int{1, 1}, []int{0, 1}, n,
// m). Negative coefficients and terms are reduced into [0, m) first. It
// panics if m <= 0, or coeffs and initial are empty or of different
// lengths.
func LinearRecurrenceMod[E constraints.Integer](coeffs, initial []E, n uint64, m E) E {
	k := len(coeffs)
	if m <= 0 {
		panic("sequences: LinearRecurrenceMod modulus must be positive")
	}
	if k == 0 || len(initial) != k {
		panic("sequences: LinearRecurrenceMod needs as many initial terms as coefficients")
	}
	mod := uint64(m)
	reduce := func(x E) uint64 {
		r := x % m
		if r < 0 {
			r += m
		}
		return uint64(r)
	}
	if n < uint64(k) {
		return E(reduce(initial[n]))
	}
	addMod := func(x, y uint64) uint64 {
		s := x + y
		if s >= mod || s < x { // also if the addition wrapped
			s -= mod
		}
		return s
	}
	mul := func(x, y [][]uint64) [][]uint64 {
		z := make([][]uint64, k)
		for i := range z {
			z[i] = make([]uint64, k)
			for l := range k {
				if x[i][l] == 0 {
					continue
				}
				for j := range k {
					z[i][j] = addMod(z[i][j], mathx.MulMod(x[i][l], y[l][j], mod))
				}
			}
		}
		return z
	}
	// The companion matrix maps (a(i-1), ..., a(i-k)) to (a(i), ...,
	// a(i-k+1)).
	base := make([][]uint64, k)
	result := make([][]uint64, k)
	for i := range k {
		base[i] = make([]uint64, k)
		result[i] = make([]uint64, k)
		result[i][i] = 1 % mod
		if i > 0 {
			base[i][i-1] = 1 % mod
		}
	}
	for j, c := range coeffs {
		base[0][j] = reduce(c)
	}
	for e := n - uint64(k-1); e > 0; e >>= 1 {
		if e&1 == 1 {
			result = mul(result, base)
		}
		base = mul(base, base)
	}
	// a(n) is the first entry of result · (a(k-1), ..., a(0)).
	var s uint64
	for j := range k {
		s = addMod(s, mathx.MulMod(result[0][j], reduce(initial[k-1-j]), mod))
	}
	return E(s)
}
