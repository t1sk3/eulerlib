package numtheory

import (
	"math"
	"math/bits"

	"github.com/t1sk3/eulerlib/v4/internal/mathx"
	"github.com/t1sk3/go-utils/types/constraints"
	"github.com/t1sk3/go-utils/types/optional"
	"github.com/t1sk3/go-utils/types/pair"
)

// powLE reports whether b^k <= n, without overflowing.
func powLE(b uint64, k int, n uint64) bool {
	r := uint64(1)
	for range k {
		hi, lo := bits.Mul64(r, b)
		if hi != 0 || lo > n {
			return false
		}
		r = lo
	}
	return true
}

// iroot returns floor(n^(1/k)) for k >= 1.
func iroot(n uint64, k int) uint64 {
	switch {
	case k == 1 || n < 2:
		return n
	case k == 2:
		return mathx.Isqrt(n)
	case k >= 64:
		return 1 // 2^64 > n
	}
	// The float estimate is within a step or two of the answer.
	r := uint64(math.Pow(float64(n), 1/float64(k)))
	for r > 0 && !powLE(r, k, n) {
		r--
	}
	for powLE(r+1, k, n) {
		r++
	}
	return r
}

// Iroot returns the integer k-th root of n, floor(n^(1/k)), exactly, for
// every value of E. It panics if n < 0 or k < 1.
func Iroot[E constraints.Integer](n E, k int) E {
	if n < 0 {
		panic("numtheory: Iroot of a negative number")
	}
	if k < 1 {
		panic("numtheory: Iroot needs k >= 1")
	}
	return E(iroot(uint64(n), k))
}

// PerfectPower writes n as base^exp with the largest possible exp >= 2,
// and reports whether that is possible: PerfectPower(64) is (2, 6, true)
// and PerfectPower(63) is (0, 0, false). It is false for n < 4 (0 and 1
// are every power of themselves, so they have no largest exponent).
func PerfectPower[E constraints.Integer](n E) (base E, exp int, ok bool) {
	if n < 4 {
		return 0, 0, false
	}
	u := uint64(n)
	for k := bits.Len64(u) - 1; k >= 2; k-- { // 2^k <= u
		r := iroot(u, k)
		if p := powOrZero(r, k); p == u {
			return E(r), k, true
		}
	}
	return 0, 0, false
}

// powOrZero returns b^k, or 0 if it does not fit in a uint64.
func powOrZero(b uint64, k int) uint64 {
	r := uint64(1)
	for range k {
		hi, lo := bits.Mul64(r, b)
		if hi != 0 {
			return 0
		}
		r = lo
	}
	return r
}

// SumOfTwoSquares returns a pair (a, b) with 0 <= a <= b and a² + b² = n,
// the one with the smallest a, or an empty Optional if n is not a sum of
// two squares (or is negative). It tries every a up to sqrt(n/2), so it
// takes O(√n).
func SumOfTwoSquares[E constraints.Integer](n E) optional.Optional[pair.Pair[E, E]] {
	if n < 0 {
		return optional.Empty[pair.Pair[E, E]]()
	}
	u := uint64(n)
	for a := uint64(0); a <= math.MaxUint32 && a*a <= u/2; a++ { // a² <= n - a²
		rest := u - a*a
		if b := mathx.Isqrt(rest); b*b == rest {
			return optional.Of(pair.Of(E(a), E(b)))
		}
	}
	return optional.Empty[pair.Pair[E, E]]()
}

// DivisorCountSieve returns a table d of length n+1 with d[i] =
// CountDivisors(i), the number of divisors of i, for every i in [0, n]
// (d[0] = 0), in O(n log n). It panics if n < 0.
func DivisorCountSieve[E constraints.Integer](n E) []E {
	if n < 0 {
		panic("numtheory: DivisorCountSieve of a negative number")
	}
	d := make([]E, uint64(n)+1)
	if n == 0 {
		return d
	}
	for i := E(1); ; i++ {
		for m := i; ; m += i {
			d[m]++
			if m > n-i {
				break
			}
		}
		if i == n {
			break
		}
	}
	return d
}
