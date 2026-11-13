package numtheory

import (
	"github.com/t1sk3/eulerlib/v4/internal/mathx"
	"github.com/t1sk3/go-utils/types/constraints"
	"github.com/t1sk3/go-utils/types/optional"
)

// distinct returns the distinct values of the sorted slice fs.
func distinct(fs []uint64) []uint64 {
	var out []uint64
	for i, f := range fs {
		if i == 0 || f != fs[i-1] {
			out = append(out, f)
		}
	}
	return out
}

func gcd64(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// MultiplicativeOrder returns the multiplicative order of a modulo m: the
// smallest k >= 1 with a^k ≡ 1 (mod m). It returns an empty Optional if no
// such k exists, which is when a and m are not coprime, or m <= 0. The
// order modulo 1 is 1.
//
// It factorizes φ(m) and strips prime factors from it, so it is fast for
// every 64-bit modulus.
func MultiplicativeOrder[E constraints.Integer](a, m E) optional.Optional[E] {
	if m <= 0 {
		return optional.Empty[E]()
	}
	ord, ok := order64(reduce(a, m), uint64(m))
	if !ok {
		return optional.Empty[E]()
	}
	return optional.Of(E(ord))
}

// order64 returns the order of a modulo m >= 1, for a in [0, m).
func order64(a, m uint64) (uint64, bool) {
	if m == 1 {
		return 1, true
	}
	if gcd64(a, m) != 1 {
		return 0, false
	}
	// The order divides φ(m): start there and divide out every prime that
	// can go.
	phi := m
	for _, p := range distinct(mathx.Factor(m)) {
		phi -= phi / p
	}
	ord := phi
	for _, p := range distinct(mathx.Factor(phi)) {
		for ord%p == 0 && mathx.PowMod(a, ord/p, m) == 1 {
			ord /= p
		}
	}
	return ord, true
}

// RecurringCycleLength returns the length of the recurring cycle in the
// decimal expansion of 1/d (Project Euler #26): 6 for d = 7
// (0.(142857)), 1 for d = 3 or d = 6, and 0 if the expansion terminates
// (d has no prime factor other than 2 and 5). The sign of d is ignored. It
// panics if d is 0.
//
// The cycle length is the multiplicative order of 10 modulo d with its
// factors of 2 and 5 removed.
func RecurringCycleLength[E constraints.Integer](d E) E {
	if d == 0 {
		panic("numtheory: RecurringCycleLength(0)")
	}
	u := mathx.Abs64(d)
	for u%2 == 0 {
		u /= 2
	}
	for u%5 == 0 {
		u /= 5
	}
	if u == 1 {
		return 0
	}
	ord, _ := order64(10%u, u)
	return E(ord)
}

// Jacobi returns the Jacobi symbol (a/n), which is -1, 0 or +1. For a
// prime n it is the Legendre symbol: 0 if n divides a, +1 if a is a
// non-zero square modulo n, and -1 if it is not. It panics unless n is odd
// and positive.
func Jacobi[E constraints.Integer](a, n E) int {
	if n <= 0 || n%2 == 0 {
		panic("numtheory: Jacobi needs an odd positive modulus")
	}
	x, m, t := reduce(a, n), uint64(n), 1
	for x != 0 {
		for x%2 == 0 {
			x /= 2
			if r := m % 8; r == 3 || r == 5 {
				t = -t
			}
		}
		x, m = m, x
		if x%4 == 3 && m%4 == 3 {
			t = -t
		}
		x %= m
	}
	if m == 1 {
		return t
	}
	return 0
}

// ModSqrt returns a square root of a modulo the prime p: the smaller of
// the two x in [0, p) with x² ≡ a (mod p), found with the Tonelli–Shanks
// algorithm. It returns an empty Optional if a is not a square modulo p.
// It panics if p is not prime.
func ModSqrt[E constraints.Integer](a, p E) optional.Optional[E] {
	if p < 2 || !mathx.IsPrime(uint64(p)) {
		panic("numtheory: ModSqrt needs a prime modulus")
	}
	x, m := reduce(a, p), uint64(p)
	if m == 2 || x == 0 {
		return optional.Of(E(x))
	}
	if mathx.PowMod(x, (m-1)/2, m) != 1 {
		return optional.Empty[E]()
	}
	var r uint64
	if m%4 == 3 {
		r = mathx.PowMod(x, m/4+1, m) // (m+1)/4 without overflow
	} else {
		// m - 1 = q * 2^s with q odd.
		q, s := m-1, 0
		for q%2 == 0 {
			q /= 2
			s++
		}
		z := uint64(2) // a non-residue
		for mathx.PowMod(z, (m-1)/2, m) != m-1 {
			z++
		}
		c := mathx.PowMod(z, q, m)
		t := mathx.PowMod(x, q, m)
		r = mathx.PowMod(x, q/2+1, m) // (q+1)/2, q is odd
		for t != 1 {
			// Least i with t^(2^i) = 1.
			i, t2 := 0, t
			for t2 != 1 {
				t2 = mathx.MulMod(t2, t2, m)
				i++
			}
			b := c
			for range s - i - 1 {
				b = mathx.MulMod(b, b, m)
			}
			s = i
			c = mathx.MulMod(b, b, m)
			t = mathx.MulMod(t, c, m)
			r = mathx.MulMod(r, b, m)
		}
	}
	return optional.Of(E(min(r, m-r)))
}

// PrimitiveRoot returns the smallest primitive root modulo the prime p: the
// smallest g whose powers run through every non-zero residue, i.e. whose
// multiplicative order is p-1. PrimitiveRoot(2) is 1. It panics if p is
// not prime.
func PrimitiveRoot[E constraints.Integer](p E) E {
	if p < 2 || !mathx.IsPrime(uint64(p)) {
		panic("numtheory: PrimitiveRoot needs a prime modulus")
	}
	m := uint64(p)
	if m == 2 {
		return 1
	}
	qs := distinct(mathx.Factor(m - 1))
next:
	for g := uint64(2); ; g++ {
		for _, q := range qs {
			if mathx.PowMod(g, (m-1)/q, m) == 1 {
				continue next
			}
		}
		return E(g)
	}
}
