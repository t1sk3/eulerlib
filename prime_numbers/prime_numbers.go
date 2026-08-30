// Package prime_numbers provides prime-related functions: primality testing,
// listing and counting primes, prime generation, and a stepping iterator.
package prime_numbers

import (
	"math"

	"github.com/t1sk3/eulerlib/v2/etc"
	"github.com/t1sk3/eulerlib/v2/utils"
)

// PrimeNumberIterator steps through the primes one at a time, starting from
// its current value. Use NewPrimeNumberIterator to create one.
type PrimeNumberIterator[E utils.Integer] struct {
	current E
}

// Proceed advances the iterator to the next prime after its current value.
func (p *PrimeNumberIterator[E]) Proceed() {
	p.current = NextPrime(p.current)
}

// Next advances the iterator to the next prime and returns it.
func (p *PrimeNumberIterator[E]) Next() E {
	p.Proceed()
	return p.current
}

// Current returns the iterator's current value without advancing it.
func (p *PrimeNumberIterator[E]) Current() E {
	return p.current
}

// Reset returns the iterator to its initial (zero) state.
func (p *PrimeNumberIterator[E]) Reset() {
	p.current = 0
}

// NewPrimeNumberIterator returns a new PrimeNumberIterator, optionally
// starting just before the given value so the first call to Next returns the
// next prime after it. At most one starting value may be given; with none,
// the iterator starts at 0.
func NewPrimeNumberIterator[E utils.Integer](params ...E) *PrimeNumberIterator[E] {
	if len(params) > 1 {
		panic("Too many arguments")
	}
	if len(params) == 1 {
		return &PrimeNumberIterator[E]{NextPrime(params[0] - 1)}
	}
	return &PrimeNumberIterator[E]{0}
}

// checks to see if the given number is a prime
func IsPrime[E utils.Integer](p E) bool {
	end := E(math.Sqrt(float64(p)))
	if end*end == p {
		return false
	}
	end++
	if p%2 == 0 && p != 2 {
		return false
	}
	if p < 2 {
		return false
	}
	for i := E(3); i < end; i += 2 {
		if p%i == 0 {
			return false
		}
	}
	return true
}

// Returns a slice where at every index the boolean in that place indicates whether or not the index is a prime number
func ListPrimality[E utils.Integer](n E) []bool {
	if n < 0 {
		panic("n must be positive")
	}
	if n == 0 {
		return []bool{}
	}
	if n == 1 {
		return []bool{false}
	}

	res := etc.GenerateSlice(n+1, true)
	res[0] = false
	res[1] = false

	for i := E(2); i*i <= n; i++ {
		if res[i] {
			for j := i * i; j <= n; j += i {
				res[j] = false
			}
		}
	}
	return res
}

// Lists all primes up to n
func ListPrimes[E utils.Integer](n E) (res []E) {
	for i, p := range ListPrimality(n) {
		if p {
			res = append(res, E(i))
		}
	}
	return res
}

// Returns a generator that generates prime numbers
func PrimeGenerator[E utils.Integer](limit E) <-chan E {
	chnl := make(chan E)
	p := NewPrimeNumberIterator[E]()
	go func() {
		for p.Next() <= limit {
			chnl <- p.Current()
		}
		close(chnl)
	}()
	return chnl
}

// Returns the next prime after n
func NextPrime[E utils.Integer](n E) E {
	if n < 2 {
		return 2
	}
	if n == 2 {
		return 3
	}
	var res E
	if n%2 == 0 {
		res = n + 1
	} else {
		res = n + 2
	}
	for !IsPrime(res) {
		res += 2
	}
	return res
}

// Sums primes between s and e
func SumPrimes[E utils.Integer](s E, e E) (res E) {
	current := NextPrime(s - 1)
	if current > e {
		return 0
	}
	// if s is prime, add it to the sum
	res = current
	for i := s; i <= e; i++ {
		current = NextPrime(current)
		if current > e {
			break
		}
		res += current
	}
	return
}

// Counts how many primes exist between s and e
func PrimeCount[E utils.Integer](s E, e E) (res E) {
	for i := s; i <= e; i++ {
		if IsPrime(i) {
			res += 1
		}
	}
	return
}

// Returns the first n prime numbers
func FirstNPrimes[E utils.Integer](n E) []E {
	res := make([]E, n)
	p := NewPrimeNumberIterator[E]()
	for i := E(0); i < n; i++ {
		res[i] = p.Next()
	}
	return res
}
