package mathx

import (
	"math/bits"
	"slices"
)

// mrBases are the first 12 primes. Using all of them as Miller–Rabin
// witnesses is deterministic for every n < 3.3e24, which covers uint64.
var mrBases = [...]uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}

// IsPrime reports whether n is prime, using deterministic Miller–Rabin.
// It is exact for every uint64.
func IsPrime(n uint64) bool {
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
	d := n - 1
	s := bits.TrailingZeros64(d)
	d >>= s
outer:
	for _, a := range mrBases {
		x := PowMod(a, d, n)
		if x == 1 || x == n-1 {
			continue
		}
		for r := 1; r < s; r++ {
			x = MulMod(x, x, n)
			if x == n-1 {
				continue outer
			}
		}
		return false
	}
	return true
}

// Factor returns the prime factors of n with repeats, in ascending order
// (nil for n < 2). It trial-divides by small primes, then splits what is
// left with Pollard's rho (Brent's variant), so it is fast for every
// uint64, including products of two large primes.
func Factor(n uint64) []uint64 {
	if n < 2 {
		return nil
	}
	var out []uint64
	for n%2 == 0 {
		out = append(out, 2)
		n /= 2
	}
	for p := uint64(3); p < 1000 && p <= n/p; p += 2 {
		for n%p == 0 {
			out = append(out, p)
			n /= p
		}
	}
	if n > 1 {
		out = splitInto(out, n)
	}
	slices.Sort(out)
	return out
}

// splitInto appends the prime factors of n > 1, which has no factor below
// 1000, to out.
func splitInto(out []uint64, n uint64) []uint64 {
	if n < 1000*1000 || IsPrime(n) {
		return append(out, n)
	}
	d := rho(n)
	return splitInto(splitInto(out, d), n/d)
}

// rho returns a non-trivial factor of the odd composite n, using Pollard's
// rho with Brent's cycle detection and batched gcds.
func rho(n uint64) uint64 {
	const batch = 128
	for c := uint64(1); ; c++ {
		f := func(x uint64) uint64 {
			s := MulMod(x, x, n) + c
			if s >= n || s < c { // reduce mod n, also if the addition wrapped
				s -= n
			}
			return s
		}
		var x, ys uint64
		y, q, g := uint64(2), uint64(1), uint64(1)
		for r := uint64(1); g == 1; r *= 2 {
			x = y
			for range r {
				y = f(y)
			}
			for k := uint64(0); k < r && g == 1; k += batch {
				ys = y
				for range min(batch, r-k) {
					y = f(y)
					q = MulMod(q, diff(x, y), n)
				}
				g = gcd(q, n)
			}
		}
		if g == n {
			// The batch overshot: redo it one step at a time.
			for g = 1; g == 1; {
				ys = f(ys)
				g = gcd(diff(x, ys), n)
			}
		}
		if g != n {
			return g
		}
	}
}

func diff(a, b uint64) uint64 {
	if a > b {
		return a - b
	}
	return b - a
}

func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
