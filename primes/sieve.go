package primes

import (
	"iter"
	"math/bits"

	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/constraints"
)

// Sieve is a sieve of Eratosthenes over [0, Limit()]. It stores one bit per
// odd number, so a sieve up to n uses about n/16 bytes. A Sieve is
// read-only once built and safe for concurrent use.
type Sieve[E constraints.Integer] struct {
	limit E
	// composite has bit i set when 2i+1 is not prime (bit 0, for 1, is set).
	composite []uint64
	count     int
}

// NewSieve sieves every number up to and including limit. It panics if
// limit < 0.
func NewSieve[E constraints.Integer](limit E) *Sieve[E] {
	if limit < 0 {
		panic("primes: NewSieve with a negative limit")
	}
	lim := uint64(limit)
	odds := lim/2 + 1 // odd numbers 1, 3, ..., <= lim (plus one spare when lim is even)
	c := make([]uint64, (odds+63)/64)
	c[0] |= 1 // 1 is not prime
	for p := uint64(3); p <= lim/p; p += 2 {
		if c[p/2/64]&(1<<(p/2%64)) != 0 {
			continue
		}
		for j := p * p / 2; j < odds; j += p {
			c[j/64] |= 1 << (j % 64)
		}
	}
	s := &Sieve[E]{limit: limit, composite: c}
	// Count the odd primes <= lim, then add 2.
	n := 0
	full := (lim + 1) / 2 // odd numbers <= lim
	for i, w := range c {
		valid := uint64(64)
		if rem := full - uint64(i)*64; rem < 64 {
			valid = rem
		}
		mask := ^uint64(0)
		if valid < 64 {
			mask = 1<<valid - 1
		}
		n += bits.OnesCount64(^w & mask)
	}
	if lim >= 2 {
		n++
	}
	s.count = n
	return s
}

// Limit returns the largest number the sieve covers.
func (s *Sieve[E]) Limit() E { return s.limit }

// Has reports whether n is prime. It panics if n > Limit(); use IsPrime for
// a check that works for any n.
func (s *Sieve[E]) Has(n E) bool {
	if n > s.limit {
		panic("primes: Sieve.Has beyond the sieve's limit")
	}
	return s.has(n)
}

func (s *Sieve[E]) has(n E) bool {
	if n < 2 {
		return false
	}
	u := uint64(n)
	if u%2 == 0 {
		return u == 2
	}
	i := u / 2
	return s.composite[i/64]&(1<<(i%64)) == 0
}

// IsPrime reports whether n is prime, using the sieve for n <= Limit() and
// the package-level IsPrime above it. Its method value is a predicate, so it
// passes straight to Stream.Filter:
//
//	sv := primes.NewSieve(1_000_000)
//	stream.Range(1, 1_000_000).Filter(sv.IsPrime)
func (s *Sieve[E]) IsPrime(n E) bool {
	if n > s.limit {
		return IsPrime(n)
	}
	return s.has(n)
}

// Count returns the number of primes <= Limit().
func (s *Sieve[E]) Count() int { return s.count }

// All returns the primes <= Limit() in ascending order.
func (s *Sieve[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		lim := uint64(s.limit)
		if lim < 2 || !yield(2) {
			return
		}
		for wi, w := range s.composite {
			free := ^w
			for free != 0 {
				b := uint64(bits.TrailingZeros64(free))
				free &= free - 1
				p := (uint64(wi)*64+b)*2 + 1
				if p > lim {
					return
				}
				if !yield(E(p)) {
					return
				}
			}
		}
	}
}

// Stream returns the primes <= Limit() as a Stream.
func (s *Sieve[E]) Stream() stream.Stream[E] {
	return stream.FromSeq(s.All())
}

// segmented calls yield with every prime p, lo <= p <= hi, in ascending
// order, until yield returns false.
func segmented(lo, hi uint64, yield func(uint64) bool) {
	if lo <= 2 && hi >= 2 {
		if !yield(2) {
			return
		}
	}
	if lo < 3 {
		lo = 3
	}
	if lo%2 == 0 {
		lo++
	}
	for c := lo; c <= hi && c >= lo; c += 2 {
		if isPrime64(c) && !yield(c) {
			return
		}
		if c > hi-2 {
			return
		}
	}
}
