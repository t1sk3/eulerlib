// Package primes provides primality testing, sieves and prime streams.
//
// IsPrime is a deterministic Miller–Rabin test, exact for every 64-bit
// value. Sieve is a bitset sieve of Eratosthenes over a fixed range, and
// All, From and Between are lazy streams backed by a segmented sieve, so
// they need no upper bound up front:
//
//	primes.All[int]().Limit(10_001).FindLast()            // the 10,001st prime
//	stream.Sum(primes.Between(2, 2_000_000))               // sum of primes up to 2e6
//	primes.All[int]().Filter(digits.IsPalindrome[int])     // palindromic primes
//
// Every stream in this package can be consumed any number of times, and
// stops at the largest value of its element type instead of overflowing.
package primes

import (
	"github.com/t1sk3/eulerlib/v4/internal/mathx"
	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/constraints"
	"github.com/t1sk3/go-utils/types/optional"
)

// mrBases are the first 12 primes. Using all of them as Miller–Rabin
// witnesses is deterministic for every n < 3.3e24, which covers uint64.
var mrBases = [...]uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}

// IsPrime reports whether n is prime. It is exact for every value of every
// integer type (deterministic Miller–Rabin), and false for n < 2.
func IsPrime[E constraints.Integer](n E) bool {
	if n < 2 {
		return false
	}
	return isPrime64(uint64(n))
}

func isPrime64(n uint64) bool {
	if n < 2 {
		return false
	}
	for _, p := range mrBases {
		if n%p == 0 {
			return n == p
		}
	}
	if n < 41*41 {
		return true
	}
	d, s := n-1, 0
	for d%2 == 0 {
		d /= 2
		s++
	}
outer:
	for _, a := range mrBases {
		x := mathx.PowMod(a, d, n)
		if x == 1 || x == n-1 {
			continue
		}
		for r := 1; r < s; r++ {
			x = mathx.MulMod(x, x, n)
			if x == n-1 {
				continue outer
			}
		}
		return false
	}
	return true
}

// NextPrime returns the smallest prime greater than n. It panics if that
// prime does not fit in E (for example NextPrime(uint8(251))).
func NextPrime[E constraints.Integer](n E) E {
	if n < 2 {
		return 2
	}
	max := uint64(mathx.MaxOf[E]())
	c := uint64(n) + 1
	if c > 2 && c%2 == 0 {
		c++
	}
	for ; c <= max && c >= uint64(n); c += 2 {
		if isPrime64(c) {
			return E(c)
		}
		if c > max-2 {
			break
		}
	}
	panic("primes: NextPrime overflows the element type")
}

// PrevPrime returns the largest prime smaller than n, or an empty Optional
// if there is none (n <= 2).
func PrevPrime[E constraints.Integer](n E) optional.Optional[E] {
	if n <= 2 {
		return optional.Empty[E]()
	}
	if n == 3 {
		return optional.Of(E(2))
	}
	c := uint64(n) - 1
	if c%2 == 0 {
		c--
	}
	for ; c >= 3; c -= 2 {
		if isPrime64(c) {
			return optional.Of(E(c))
		}
	}
	return optional.Of(E(2))
}

// All returns every prime in ascending order: 2, 3, 5, 7, ... The stream is
// infinite in practice (it ends at the largest prime that fits in E), so
// bound it with Limit or TakeWhile before a terminal operation that reads
// every element.
func All[E constraints.Integer]() stream.Stream[E] {
	return From(E(2))
}

// From returns every prime >= n in ascending order, ending at the largest
// prime that fits in E.
func From[E constraints.Integer](n E) stream.Stream[E] {
	return Between(n, mathx.MaxOf[E]())
}

// Between returns every prime p with lo <= p <= hi, in ascending order.
func Between[E constraints.Integer](lo, hi E) stream.Stream[E] {
	if hi < 2 || lo > hi {
		return stream.Empty[E]()
	}
	l := uint64(0)
	if lo > 0 {
		l = uint64(lo)
	}
	h := uint64(hi)
	return stream.FromSeq(func(yield func(E) bool) {
		segmented(l, h, func(p uint64) bool { return yield(E(p)) })
	})
}

// UpTo returns every prime <= n, in ascending order.
func UpTo[E constraints.Integer](n E) []E {
	if n < 2 {
		return []E{}
	}
	sv := NewSieve(n)
	out := make([]E, 0, sv.Count())
	for p := range sv.All() {
		out = append(out, p)
	}
	return out
}

// First returns the first n primes. If fewer than n primes fit in E, it
// returns all of those that do.
func First[E constraints.Integer](n int) []E {
	if n <= 0 {
		return []E{}
	}
	return All[E]().Limit(n).ToSlice()
}

// Count returns the number of primes p with lo <= p <= hi.
func Count[E constraints.Integer](lo, hi E) int {
	return Between(lo, hi).Count()
}

// Sum returns the sum of the primes p with lo <= p <= hi. Like ordinary Go
// arithmetic, it overflows E silently if the sum does not fit.
func Sum[E constraints.Integer](lo, hi E) E {
	return stream.Sum(Between(lo, hi))
}

// SmallestFactorSieve returns a table spf of length n+1 where spf[i] is the
// smallest prime factor of i, for i in [0, n]; spf[0] and spf[1] are 0.
// Building it takes O(n log log n); with it, numtheory.FactorizeSPF
// factorizes any number <= n in O(log n). It panics if n < 0.
func SmallestFactorSieve[E constraints.Integer](n E) []E {
	if n < 0 {
		panic("primes: SmallestFactorSieve of a negative number")
	}
	spf := make([]E, uint64(n)+1)
	if n < 2 {
		return spf
	}
	for i := E(2); ; i++ {
		if spf[i] == 0 {
			spf[i] = i
			if i <= n/i {
				for j := i * i; ; j += i {
					if spf[j] == 0 {
						spf[j] = i
					}
					if j > n-i {
						break
					}
				}
			}
		}
		if i == n {
			break
		}
	}
	return spf
}
