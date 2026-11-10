// Package numtheory provides number-theory functions: divisors and
// divisor sums, prime factorization, gcd/lcm, modular arithmetic, powers
// and roots, arithmetic functions (totient, Möbius), predicate builders for
// Stream.Filter, and collectors for Stream.Collect.
//
// Functions document what they do with zero and negative inputs. As in the
// rest of eulerlib, invalid arguments that are programmer errors (a zero
// modulus, a negative exponent) panic; arithmetic that overflows E wraps
// silently, like ordinary Go arithmetic, unless the documentation says
// otherwise.
package numtheory

import (
	"slices"

	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/constraints"
	"github.com/t1sk3/go-utils/types/pair"
)

// Divisors returns every positive divisor of n in ascending order,
// including 1 and n. It returns an empty slice for n < 1. The divisors are
// built from n's prime factorization, so this is fast even for large n.
func Divisors[E constraints.Integer](n E) []E {
	if n < 1 {
		return []E{}
	}
	ds := []E{1}
	for p, e := range Factors(n) {
		k := len(ds)
		pk := E(1)
		for range e {
			pk *= p
			for _, d := range ds[:k] {
				ds = append(ds, d*pk)
			}
		}
	}
	slices.Sort(ds)
	return ds
}

// CountDivisors returns the number of positive divisors of n, computed from
// its prime factorization. It returns 0 for n < 1.
func CountDivisors[E constraints.Integer](n E) int {
	if n < 1 {
		return 0
	}
	c := 1
	for _, e := range Factors(n) {
		c *= e + 1
	}
	return c
}

// SumDivisors returns the sum of the proper divisors of n (every positive
// divisor except n itself). It returns 0 for n <= 1. It is computed from
// n's prime factorization as σ(n) - n in wrapping arithmetic, so the result
// is correct whenever it fits in E, even if σ(n) does not.
func SumDivisors[E constraints.Integer](n E) E {
	if n <= 1 {
		return 0
	}
	sigma := E(1)
	for p, e := range Factors(n) {
		term, pk := E(1), E(1)
		for range e {
			pk *= p
			term += pk
		}
		sigma *= term
	}
	return sigma - n
}

// SigmaSieve returns a table s of length n+1 where s[i] = SumDivisors(i),
// the sum of the proper divisors of i, for every i in [0, n]. It runs in
// O(n log n), much faster than calling SumDivisors for every i when you
// need all of them (Project Euler #21, #23). It panics if n < 0.
func SigmaSieve[E constraints.Integer](n E) []E {
	if n < 0 {
		panic("numtheory: SigmaSieve of a negative number")
	}
	s := make([]E, uint64(n)+1)
	for d := E(1); d <= n/2; d++ {
		for m := 2 * d; ; m += d {
			s[m] += d
			if m > n-d {
				break
			}
		}
	}
	return s
}

// IsAbundant reports whether the proper divisors of n add up to more than n.
func IsAbundant[E constraints.Integer](n E) bool {
	return SumDivisors(n) > n
}

// IsDeficient reports whether the proper divisors of n add up to less than
// n.
func IsDeficient[E constraints.Integer](n E) bool {
	return SumDivisors(n) < n
}

// IsPerfect reports whether the proper divisors of n add up to exactly n.
func IsPerfect[E constraints.Integer](n E) bool {
	return n > 0 && SumDivisors(n) == n
}

// IsAmicable reports whether n is part of an amicable pair: m =
// SumDivisors(n) differs from n and SumDivisors(m) == n.
func IsAmicable[E constraints.Integer](n E) bool {
	m := SumDivisors(n)
	return m != n && SumDivisors(m) == n
}

// AmicablePairs returns every amicable pair (a, b) with a < b < limit, in
// ascending order of a (Project Euler #21). It uses SigmaSieve, so it needs
// O(limit) memory.
func AmicablePairs[E constraints.Integer](limit E) stream.Stream[pair.Pair[E, E]] {
	return stream.FromSeq(func(yield func(pair.Pair[E, E]) bool) {
		if limit <= 2 {
			return
		}
		s := SigmaSieve(limit - 1)
		for a := E(2); a < limit; a++ {
			if b := s[a]; b > a && b < limit && s[b] == a {
				if !yield(pair.Of(a, b)) {
					return
				}
			}
		}
	})
}
