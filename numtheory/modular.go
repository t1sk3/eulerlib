package numtheory

import (
	"github.com/t1sk3/eulerlib/v4/internal/mathx"
	"github.com/t1sk3/go-utils/types/constraints"
	"github.com/t1sk3/go-utils/types/optional"
	"github.com/t1sk3/go-utils/types/pair"
)

// Gcd returns the greatest common divisor of args, which is never
// negative. Gcd() is 0, Gcd(n) is |n|, and Gcd(0, n) is |n|.
func Gcd[E constraints.Integer](args ...E) E {
	var g E
	for _, v := range args {
		g = gcd(g, v)
	}
	return abs(g)
}

func gcd[E constraints.Integer](a, b E) E {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func abs[E constraints.Integer](n E) E {
	if n < 0 {
		return -n
	}
	return n
}

// Lcm returns the least common multiple of args, which is never negative.
// Lcm() is 1, and the Lcm of anything with 0 is 0. Overflows E silently if
// the result does not fit.
func Lcm[E constraints.Integer](args ...E) E {
	l := E(1)
	for _, v := range args {
		l = lcm(l, v)
	}
	return l
}

func lcm[E constraints.Integer](a, b E) E {
	if a == 0 || b == 0 {
		return 0
	}
	return abs(a / gcd(a, b) * b)
}

// IsCoprime reports whether a and b have no common factor other than 1.
func IsCoprime[E constraints.Integer](a, b E) bool {
	return Gcd(a, b) == 1
}

// ExtGcd returns g = gcd(a, b) together with x and y such that
// a*x + b*y = g (Bézout's identity), using the extended Euclidean
// algorithm. g is never negative.
func ExtGcd[E constraints.Signed](a, b E) (g, x, y E) {
	oldR, r := a, b
	oldS, s := E(1), E(0)
	oldT, t := E(0), E(1)
	for r != 0 {
		q := oldR / r
		oldR, r = r, oldR-q*r
		oldS, s = s, oldS-q*s
		oldT, t = t, oldT-q*t
	}
	if oldR < 0 {
		return -oldR, -oldS, -oldT
	}
	return oldR, oldS, oldT
}

// ModInverse returns the inverse of a modulo m: the x in [0, |m|) with
// a*x ≡ 1 (mod m). It returns an empty Optional if a and m are not coprime
// (so no inverse exists) or m is 0.
func ModInverse[E constraints.Signed](a, m E) optional.Optional[E] {
	if m == 0 {
		return optional.Empty[E]()
	}
	m = abs(m)
	g, x, _ := ExtGcd(a, m)
	if g != 1 {
		return optional.Empty[E]()
	}
	x %= m
	if x < 0 {
		x += m
	}
	return optional.Of(x)
}

// CRT solves the system x ≡ remainders[i] (mod moduli[i]) for every i
// (the Chinese remainder theorem; the moduli need not be pairwise coprime).
// It returns the solution x in [0, M) paired with M = lcm(moduli), or an
// empty Optional if the system has no solution, a modulus is 0, or the
// slices are empty or of different lengths. The result is exact as long
// as M fits in E.
func CRT[E constraints.Signed](remainders, moduli []E) optional.Optional[pair.Pair[E, E]] {
	none := optional.Empty[pair.Pair[E, E]]()
	if len(remainders) == 0 || len(remainders) != len(moduli) {
		return none
	}
	var x, m E = 0, 1
	for i := range remainders {
		r, n := remainders[i], abs(moduli[i])
		if n == 0 {
			return none
		}
		r %= n
		if r < 0 {
			r += n
		}
		g, p, _ := ExtGcd(m, n)
		if (r-x)%g != 0 {
			return none
		}
		l := m / g * n
		step := n / g
		// x + m*k with k = (r-x)/g * p mod step solves both congruences.
		// MulMod keeps the product exact, and x + m*k < l, so nothing
		// overflows as long as l itself fits in E.
		k := MulMod((r-x)/g, p, step)
		x += m * k
		m = l
	}
	return optional.Of(pair.Of(x, m))
}

// PowMod returns x^y mod m in [0, m). Negative x is reduced into [0, m)
// first. The result is exact for every m that fits in E: intermediate
// products use 128 bits. It panics if m <= 0 or y < 0.
func PowMod[E constraints.Integer](x, y, m E) E {
	if m <= 0 {
		panic("numtheory: PowMod modulus must be positive")
	}
	if y < 0 {
		panic("numtheory: PowMod negative exponent")
	}
	return E(mathx.PowMod(reduce(x, m), uint64(y), uint64(m)))
}

// MulMod returns a*b mod m in [0, m) without overflowing E. Negative
// operands are reduced into [0, m) first. It panics if m <= 0.
func MulMod[E constraints.Integer](a, b, m E) E {
	if m <= 0 {
		panic("numtheory: MulMod modulus must be positive")
	}
	return E(mathx.MulMod(reduce(a, m), reduce(b, m), uint64(m)))
}

// reduce returns x mod m in [0, m) as a uint64, for m > 0.
func reduce[E constraints.Integer](x, m E) uint64 {
	r := x % m
	if r < 0 {
		r += m
	}
	return uint64(r)
}
