package eulerlib

import (
	"math"
)

type PrimeNumberIterator[E Integer] struct {
	current E
}

func (p *PrimeNumberIterator[E]) Proceed() {
	p.current = NextPrime(p.current)
}

func (p *PrimeNumberIterator[E]) Next() E {
	p.Proceed()
	return p.current
}

func (p *PrimeNumberIterator[E]) Current() E {
	return p.current
}

func (p *PrimeNumberIterator[E]) Reset() {
	p.current = 0
}

// function that returns a new prime number iterator with an optional starting point, only one argument is allowed
// if no arguments are given, the iterator starts at 0 and a type must be specified
func NewPrimeNumberIterator[E Integer](params ...E) *PrimeNumberIterator[E] {
	if len(params) > 1 {
		panic("Too many arguments")
	}
	if len(params) == 1 {
		return &PrimeNumberIterator[E]{NextPrime(params[0] - 1)}
	}
	return &PrimeNumberIterator[E]{0}
}

// checks to see if the given number is a prime
func IsPrime[E Integer](p E) bool {
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
func ListPrimality[E Integer](n E) []bool {
	if n < 0 {
		panic("n must be positive")
	}
	if n == 0 {
		return []bool{}
	}
	if n == 1 {
		return []bool{false}
	}

	res := GenerateSlice(n+1, true)
	res[0] = false
	res[1] = false

	for _, e := range ListPrimes(n) {
		res[E(e)] = true
	}
	return res
}

// Lists all primes up to n
func ListPrimes[E Integer](n E) []E {
	res := []E{}
	p := NewPrimeNumberIterator[E]()
	for p.Next() <= n {
		res = append(res, p.Current())
	}
	return res
}

// Returns a generator that generates prime numbers
func PrimeGenerator[E Integer](limit E) <-chan E {
	chnl := make(chan E)
	p := NewPrimeNumberIterator[E]()
	go func() {
		for {
			next := p.Next()
			if next > limit {
				break
			}
			chnl <- next
		}
		close(chnl)
	}()
	return chnl
}

// Returns the next prime after n
func NextPrime[E Integer](n E) E {
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
func SumPrimes[E Integer](s E, e E) (res E) {
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
func PrimeCount[E Integer](s E, e E) (res E) {
	for i := s; i <= e; i++ {
		if IsPrime(i) {
			res += 1
		}
	}
	return
}

// Returns the first n prime numbers
func FirstNPrimes[E Integer](n E) []E {
	res := make([]E, n)
	p := NewPrimeNumberIterator[E]()
	for i := E(0); i < n; i++ {
		res[i] = p.Next()
	}
	return res
}
