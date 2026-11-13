package primes

import (
	"math/big"
	"math/bits"

	"github.com/t1sk3/eulerlib/v4/digits"
	"github.com/t1sk3/eulerlib/v4/internal/mathx"
	"github.com/t1sk3/go-utils/types/constraints"
)

// Pi returns π(n), the number of primes <= n, without enumerating them. It
// uses the Lucy_Hedgehog method: O(n^(3/4)) time and O(√n) memory, so
// π(10^11) takes about a second, where Count(2, n) would sieve the whole
// range. It returns 0 for n < 2.
func Pi[E constraints.Integer](n E) E {
	if n < 2 {
		return 0
	}
	u := uint64(n)
	r := mathx.Isqrt(u)
	// lo[v] counts the numbers in [2, v] that survive the sieving so far,
	// for v <= r; hi[i] does the same for v = n/i.
	lo := make([]uint64, r+1)
	hi := make([]uint64, r+1)
	for i := uint64(1); i <= r; i++ {
		lo[i] = i - 1
		hi[i] = u/i - 1
	}
	for p := uint64(2); p <= r; p++ {
		if lo[p] == lo[p-1] {
			continue // p is not prime
		}
		below, p2 := lo[p-1], p*p
		for i := uint64(1); i <= r && u/i >= p2; i++ {
			if d := i * p; d <= r {
				hi[i] -= hi[d] - below
			} else {
				hi[i] -= lo[u/d] - below
			}
		}
		for v := r; v >= p2; v-- {
			lo[v] -= lo[v/p] - below
		}
	}
	return E(hi[1])
}

// u128 is an unsigned 128-bit integer, enough for the sum of the primes
// below 2^64.
type u128 struct{ hi, lo uint64 }

func (a u128) sub(b u128) u128 {
	lo, borrow := bits.Sub64(a.lo, b.lo, 0)
	hi, _ := bits.Sub64(a.hi, b.hi, borrow)
	return u128{hi, lo}
}

func (a u128) mul(m uint64) u128 {
	hi, lo := bits.Mul64(a.lo, m)
	return u128{hi + a.hi*m, lo}
}

// sumFrom2 returns 2 + 3 + ... + v = v(v+1)/2 - 1 for v >= 1.
func sumFrom2(v uint64) u128 {
	a, b := v/2, v|1 // v even: (v/2)(v+1); v odd: ((v+1)/2)·v
	if v%2 == 1 {
		a++
	}
	hi, lo := bits.Mul64(a, b)
	return u128{hi, lo}.sub(u128{0, 1})
}

// SumUpTo returns the sum of the primes <= n without enumerating them
// (Project Euler #10 for much larger limits), using the same method, time
// and memory as Pi. The result is a *big.Int because it outgrows a uint64
// near n = 10^11. It returns 0 for n < 2.
func SumUpTo[E constraints.Integer](n E) *big.Int {
	if n < 2 {
		return new(big.Int)
	}
	u := uint64(n)
	r := mathx.Isqrt(u)
	lo := make([]u128, r+1)
	hi := make([]u128, r+1)
	for i := uint64(1); i <= r; i++ {
		lo[i] = sumFrom2(i)
		hi[i] = sumFrom2(u / i)
	}
	for p := uint64(2); p <= r; p++ {
		if lo[p] == lo[p-1] {
			continue
		}
		below, p2 := lo[p-1], p*p
		for i := uint64(1); i <= r && u/i >= p2; i++ {
			if d := i * p; d <= r {
				hi[i] = hi[i].sub(hi[d].sub(below).mul(p))
			} else {
				hi[i] = hi[i].sub(lo[u/d].sub(below).mul(p))
			}
		}
		for v := r; v >= p2; v-- {
			lo[v] = lo[v].sub(lo[v/p].sub(below).mul(p))
		}
	}
	s := new(big.Int).SetUint64(hi[1].hi)
	return s.Lsh(s, 64).Add(s, new(big.Int).SetUint64(hi[1].lo))
}

// Nth returns the n-th prime, counting from Nth(1) = 2 (Project Euler #7).
// It panics if n < 1 or the n-th prime does not fit in E.
func Nth[E constraints.Integer](n int) E {
	if n < 1 {
		panic("primes: Nth needs n >= 1")
	}
	p, ok := All[E]().Skip(n - 1).FindFirst().Value()
	if !ok {
		panic("primes: Nth overflows the element type")
	}
	return p
}

// IsPrimeBig reports whether n is prime. It is exact for n < 2^64; above
// that it is the Baillie–PSW test plus 20 Miller–Rabin rounds of
// big.Int.ProbablyPrime, for which no composite is known to pass.
func IsPrimeBig(n *big.Int) bool {
	return n.ProbablyPrime(20)
}

// IsCircular reports whether every rotation of the decimal digits of n is
// prime, n included (Project Euler #35): 197, 971 and 719 are all prime.
func IsCircular[E constraints.Integer](n E) bool {
	if n < 2 {
		return false
	}
	for _, r := range digits.Rotations(uint64(n)) {
		if !isPrime64(r) {
			return false
		}
	}
	return true
}

// IsLeftTruncatable reports whether n stays prime as digits are removed
// from the left, one at a time, n included: 3797, 797, 97 and 7. It is
// false for a number with a zero digit.
func IsLeftTruncatable[E constraints.Integer](n E) bool {
	if n < 2 {
		return false
	}
	u := uint64(n)
	// A zero digit would leave a truncation with a leading zero.
	for t := u; t > 0; t /= 10 {
		if t%10 == 0 {
			return false
		}
	}
	pow := uint64(1) // the place value of m's leading digit
	for u/pow >= 10 {
		pow *= 10
	}
	for m := u; ; pow /= 10 {
		if !isPrime64(m) {
			return false
		}
		if pow == 1 {
			return true
		}
		m %= pow
	}
}

// IsRightTruncatable reports whether n stays prime as digits are removed
// from the right, one at a time, n included: 3797, 379, 37 and 3.
func IsRightTruncatable[E constraints.Integer](n E) bool {
	if n < 2 {
		return false
	}
	for u := uint64(n); u > 0; u /= 10 {
		if !isPrime64(u) {
			return false
		}
	}
	return true
}

// IsTruncatable reports whether n is both left- and right-truncatable
// (Project Euler #37). The single-digit primes 2, 3, 5 and 7 pass, though
// that problem does not count them.
func IsTruncatable[E constraints.Integer](n E) bool {
	return IsRightTruncatable(n) && IsLeftTruncatable(n)
}
