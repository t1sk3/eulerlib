package numtheory

import (
	"math"
	"math/big"

	"github.com/t1sk3/eulerlib/v4/internal/mathx"
	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/constraints"
)

// Pow returns b^n by binary exponentiation. It overflows E silently if the
// result does not fit. It panics if n < 0.
func Pow[E constraints.Integer](b E, n int) E {
	if n < 0 {
		panic("numtheory: Pow negative exponent")
	}
	r := E(1)
	for n > 0 {
		if n&1 == 1 {
			r *= b
		}
		b *= b
		n >>= 1
	}
	return r
}

// PowBig returns b^n as a new *big.Int. It does not modify b, and panics
// if n < 0.
func PowBig(b *big.Int, n int64) *big.Int {
	if n < 0 {
		panic("numtheory: PowBig negative exponent")
	}
	return new(big.Int).Exp(b, big.NewInt(n), nil)
}

// PowBigFloat returns b^n as a new *big.Float with b's precision, by binary
// exponentiation. It does not modify b, and panics if n < 0.
func PowBigFloat(b *big.Float, n int64) *big.Float {
	if n < 0 {
		panic("numtheory: PowBigFloat negative exponent")
	}
	result := new(big.Float).SetPrec(b.Prec()).SetInt64(1)
	base := new(big.Float).Copy(b)
	for n > 0 {
		if n&1 == 1 {
			result.Mul(result, base)
		}
		base.Mul(base, base)
		n >>= 1
	}
	return result
}

// IsPowerOfTwo reports whether n is a positive power of two (1, 2, 4, ...).
func IsPowerOfTwo[E constraints.Integer](n E) bool {
	return n > 0 && n&(n-1) == 0
}

// IsSquare reports whether n is a perfect square. It is exact for every
// value of E (it uses an integer square root, not math.Sqrt).
func IsSquare[E constraints.Integer](n E) bool {
	if n < 0 {
		return false
	}
	// Squares are 0, 1, 4 or 9 mod 16: rule out most others cheaply.
	if (0x0213>>(uint64(n)&15))&1 == 0 {
		return false
	}
	r := mathx.Isqrt(uint64(n))
	return r*r == uint64(n)
}

// Isqrt returns floor(sqrt(n)), exactly, for every value of E. It panics if
// n < 0.
func Isqrt[E constraints.Integer](n E) E {
	if n < 0 {
		panic("numtheory: Isqrt of a negative number")
	}
	return E(mathx.Isqrt(uint64(n)))
}

// FloatIsInteger reports whether n is a finite float with no fractional
// part. It is false for infinities and NaN.
func FloatIsInteger[E constraints.Float](n E) bool {
	f := float64(n)
	return !math.IsInf(f, 0) && math.Trunc(f) == f
}

// Totient returns Euler's totient φ(n): how many integers in [1, n] are
// coprime to n. It returns 0 for n <= 0.
func Totient[E constraints.Integer](n E) E {
	if n <= 0 {
		return 0
	}
	res := n
	for p := range Factors(n) {
		res -= res / p
	}
	return res
}

// TotientSieve returns a table t of length n+1 with t[i] = φ(i) for every
// i in [0, n] (t[0] = 0), in O(n log log n). It panics if n < 0.
func TotientSieve[E constraints.Integer](n E) []E {
	if n < 0 {
		panic("numtheory: TotientSieve of a negative number")
	}
	t := make([]E, uint64(n)+1)
	for i := range t {
		t[i] = E(i)
	}
	for i := E(2); i <= n && i > 0; i++ {
		if t[i] == i { // i is prime
			for j := i; ; j += i {
				t[j] -= t[j] / i
				if j > n-i {
					break
				}
			}
		}
		if i == n {
			break
		}
	}
	return t
}

// MobiusSieve returns a table mu of length n+1 with mu[i] = μ(i), the
// Möbius function, for every i in [0, n]: μ(0) = 0, μ(1) = 1, μ(i) = 0 if
// a square divides i, and (-1)^k if i is the product of k distinct primes.
// It uses a linear sieve, O(n). It panics if n < 0.
func MobiusSieve[E constraints.Signed](n E) []E {
	if n < 0 {
		panic("numtheory: MobiusSieve of a negative number")
	}
	mu := make([]E, uint64(n)+1)
	composite := make([]bool, uint64(n)+1)
	var ps []E
	if n >= 1 {
		mu[1] = 1
	}
	for i := E(2); i <= n && i > 0; i++ {
		if !composite[i] {
			ps = append(ps, i)
			mu[i] = -1
		}
		for _, p := range ps {
			if p > n/i {
				break
			}
			composite[i*p] = true
			if i%p == 0 {
				mu[i*p] = 0
				break
			}
			mu[i*p] = -mu[i]
		}
		if i == n {
			break
		}
	}
	return mu
}

// Abs returns |n|. Like ordinary Go arithmetic, Abs of the minimum value of
// a signed type overflows and returns that value unchanged.
func Abs[E constraints.Signed](n E) E {
	return abs(n)
}

// SumOfSquares returns 1² + 2² + ... + n² = n(n+1)(2n+1)/6.
func SumOfSquares[E constraints.Integer](n E) E {
	return n * (n + 1) * (2*n + 1) / 6
}

// SquareOfSum returns (1 + 2 + ... + n)² = (n(n+1)/2)².
func SquareOfSum[E constraints.Integer](n E) E {
	s := n * (n + 1) / 2
	return s * s
}

// ToRadians converts an angle in degrees to radians.
func ToRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

// DivisibleBy returns a predicate reporting whether its argument is a
// multiple of k. It panics if k is 0.
func DivisibleBy[E constraints.Integer](k E) func(E) bool {
	if k == 0 {
		panic("numtheory: DivisibleBy(0)")
	}
	return func(n E) bool { return n%k == 0 }
}

// DivisibleByAny returns a predicate reporting whether its argument is a
// multiple of at least one of ks. Project Euler #1:
//
//	stream.Sum(stream.Range(1, 1000).Filter(numtheory.DivisibleByAny(3, 5)))
//
// It panics if any of ks is 0.
func DivisibleByAny[E constraints.Integer](ks ...E) func(E) bool {
	for _, k := range ks {
		if k == 0 {
			panic("numtheory: DivisibleByAny(0)")
		}
	}
	return func(n E) bool {
		for _, k := range ks {
			if n%k == 0 {
				return true
			}
		}
		return false
	}
}

// CoprimeTo returns a predicate reporting whether its argument is coprime
// to k.
func CoprimeTo[E constraints.Integer](k E) func(E) bool {
	return func(n E) bool { return IsCoprime(n, k) }
}

// Product returns a collector that multiplies a stream's elements
// together. The product of an empty stream is 1.
func Product[E constraints.Number]() stream.Collector[E, E] {
	return stream.NewCollector(
		func() E { return 1 },
		func(acc, v E) E { return acc * v },
		stream.Identity[E],
	)
}

// GcdOf returns a collector computing the gcd of a stream's elements. The
// gcd of an empty stream is 0.
func GcdOf[E constraints.Integer]() stream.Collector[E, E] {
	return stream.NewCollector(
		func() E { return 0 },
		gcd[E],
		abs[E],
	)
}

// LcmOf returns a collector computing the lcm of a stream's elements. The
// lcm of an empty stream is 1. Project Euler #5:
//
//	stream.RangeClosed(1, 20).Collect(numtheory.LcmOf[int]())
func LcmOf[E constraints.Integer]() stream.Collector[E, E] {
	return stream.NewCollector(
		func() E { return 1 },
		lcm[E],
		stream.Identity[E],
	)
}
