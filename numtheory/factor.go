package numtheory

import (
	"iter"
	"maps"
	"math/big"

	"github.com/t1sk3/eulerlib/v4/internal/mathx"

	"github.com/t1sk3/go-utils/types/constraints"
	"github.com/t1sk3/go-utils/types/pair"
)

// Factors returns the prime factorization of n as (prime, exponent) pairs
// in ascending order of the prime. It yields nothing for n < 2. Collect it
// into a map with maps.Collect(Factors(n)), or range over it directly:
//
//	for p, e := range numtheory.Factors(360) { ... } // (2, 3), (3, 2), (5, 1)
//
// It trial-divides by small primes and splits larger cofactors with
// Pollard's rho, so it is fast for every 64-bit value, including the
// product of two large primes. The sequence can be ranged over any number
// of times.
func Factors[E constraints.Integer](n E) iter.Seq2[E, int] {
	return func(yield func(E, int) bool) {
		if n < 2 {
			return
		}
		fs := mathx.Factor(uint64(n))
		for i := 0; i < len(fs); {
			j := i + 1
			for j < len(fs) && fs[j] == fs[i] {
				j++
			}
			if !yield(E(fs[i]), j-i) {
				return
			}
			i = j
		}
	}
}

// Factorize returns the prime factorization of n as a map of prime to
// exponent. It returns an empty map for n < 2. It is
// maps.Collect(Factors(n)).
func Factorize[E constraints.Integer](n E) map[E]int {
	return maps.Collect(Factors(n))
}

// FactorsSPF is Factors using a smallest-prime-factor table spf, as built
// by primes.SmallestFactorSieve(m) for some m >= n. It runs in O(log n), so
// build the table once and reuse it when you factorize many numbers. It
// panics if n is beyond the table, or the table is not a valid
// smallest-prime-factor table.
func FactorsSPF[E constraints.Integer](n E, spf []E) iter.Seq2[E, int] {
	return func(yield func(E, int) bool) {
		m := n
		if m < 2 {
			return
		}
		if uint64(m) >= uint64(len(spf)) {
			panic("numtheory: spf table too small")
		}
		for m > 1 {
			p := spf[m]
			if p < 2 {
				panic("numtheory: spf table missing factor")
			}
			e := 0
			for m%p == 0 {
				m /= p
				e++
			}
			if !yield(p, e) {
				return
			}
		}
	}
}

// FactorizeSPF is Factorize using a smallest-prime-factor table; see
// FactorsSPF.
func FactorizeSPF[E constraints.Integer](n E, spf []E) map[E]int {
	return maps.Collect(FactorsSPF(n, spf))
}

// PrimeFactors returns the prime factors of n with repeats, in ascending
// order: PrimeFactors(12) is [2 2 3]. It returns an empty slice for n < 2.
func PrimeFactors[E constraints.Integer](n E) []E {
	out := []E{}
	for p, e := range Factors(n) {
		for range e {
			out = append(out, p)
		}
	}
	return out
}

// DistinctPrimeFactors returns the distinct prime factors of n in
// ascending order: DistinctPrimeFactors(12) is [2 3] (Project Euler #47).
func DistinctPrimeFactors[E constraints.Integer](n E) []E {
	out := []E{}
	for p := range Factors(n) {
		out = append(out, p)
	}
	return out
}

// FactorizeBig returns the prime factorization of n as (prime, exponent)
// pairs in ascending order of the prime. It returns an empty slice for
// n < 2, and does not modify n. It uses trial division and stops as soon
// as the remaining cofactor is prime (big.Int.ProbablyPrime, which is exact
// below 2^64 and has no known counterexample above), so it is fast whenever
// n has at most one large prime factor, and slow when n has two or more.
func FactorizeBig(n *big.Int) []pair.Pair[*big.Int, int] {
	out := []pair.Pair[*big.Int, int]{}
	one := big.NewInt(1)
	if n.Cmp(one) <= 0 {
		return out
	}
	m := new(big.Int).Set(n)
	q, r := new(big.Int), new(big.Int)
	divide := func(p *big.Int) (divided bool) {
		e := 0
		for {
			q.QuoRem(m, p, r)
			if r.Sign() != 0 {
				break
			}
			m.Set(q)
			e++
		}
		if e > 0 {
			out = append(out, pair.Of(new(big.Int).Set(p), e))
		}
		return e > 0
	}
	divide(big.NewInt(2))
	p := big.NewInt(3)
	two := big.NewInt(2)
	sq := new(big.Int)
	// The cofactor only needs a primality test when it has changed.
	check := true
	for m.Cmp(one) > 0 {
		if (check && m.ProbablyPrime(20)) || sq.Mul(p, p).Cmp(m) > 0 {
			out = append(out, pair.Of(new(big.Int).Set(m), 1))
			break
		}
		check = divide(p)
		p.Add(p, two)
	}
	return out
}
