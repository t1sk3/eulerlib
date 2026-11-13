package combinatorics

import (
	"math/big"
	"slices"

	"github.com/t1sk3/eulerlib/v4/internal/mathx"
	"github.com/t1sk3/go-utils/types/constraints"
)

// Catalan returns the n-th Catalan number C(2n, n) / (n+1): 1, 1, 2, 5,
// 14, 42, ... It counts, among much else, the balanced strings of n pairs
// of brackets and the lattice paths that stay below the diagonal. It
// panics if n < 0.
func Catalan(n int) *big.Int {
	if n < 0 {
		panic("combinatorics: Catalan of a negative number")
	}
	c := new(big.Int).Binomial(2*int64(n), int64(n))
	return c.Quo(c, big.NewInt(int64(n)+1))
}

// Derangements returns the number of permutations of n elements that
// leave no element in its place (the subfactorial !n): 1, 0, 1, 2, 9, 44,
// ... It panics if n < 0.
func Derangements(n int) *big.Int {
	if n < 0 {
		panic("combinatorics: Derangements of a negative number")
	}
	// D(n) = (n-1)(D(n-1) + D(n-2)), with D(0) = 1 and D(1) = 0.
	a, b := big.NewInt(1), big.NewInt(0)
	if n == 0 {
		return a
	}
	for i := 2; i <= n; i++ {
		a.Add(a, b).Mul(a, big.NewInt(int64(i-1)))
		a, b = b, a
	}
	return b
}

// Stirling1 returns the unsigned Stirling number of the first kind: the
// number of permutations of n elements with exactly k cycles. It is 0
// when k > n, and when exactly one of n and k is 0. It panics if n or k is
// negative.
func Stirling1(n, k int) *big.Int {
	// c(i, j) = (i-1)·c(i-1, j) + c(i-1, j-1)
	return stirling(n, k, func(i, j int) int64 { return int64(i - 1) })
}

// Stirling2 returns the Stirling number of the second kind: the number of
// ways to split a set of n elements into exactly k non-empty subsets. It
// is 0 when k > n, and when exactly one of n and k is 0. It panics if n or
// k is negative.
func Stirling2(n, k int) *big.Int {
	// S(i, j) = j·S(i-1, j) + S(i-1, j-1)
	return stirling(n, k, func(i, j int) int64 { return int64(j) })
}

// stirling runs the recurrence t(i, j) = factor(i, j)·t(i-1, j) +
// t(i-1, j-1) from t(0, 0) = 1 up to t(n, k).
func stirling(n, k int, factor func(i, j int) int64) *big.Int {
	if n < 0 || k < 0 {
		panic("combinatorics: Stirling number with a negative argument")
	}
	if k > n {
		return new(big.Int)
	}
	row := make([]*big.Int, k+1)
	for j := range row {
		row[j] = new(big.Int)
	}
	row[0].SetInt64(1)
	f := new(big.Int)
	for i := 1; i <= n; i++ {
		// Right to left, so row[j-1] is still row i-1's value.
		for j := min(i, k); j >= 1; j-- {
			row[j].Mul(row[j], f.SetInt64(factor(i, j))).Add(row[j], row[j-1])
		}
		row[0].SetInt64(0)
	}
	return row[k]
}

// Bell returns the n-th Bell number, the number of ways to split a set of
// n elements into non-empty subsets: 1, 1, 2, 5, 15, 52, ... It panics if
// n < 0.
func Bell(n int) *big.Int {
	if n < 0 {
		panic("combinatorics: Bell of a negative number")
	}
	// The Bell triangle: each row starts with the last entry of the row
	// before, and every entry adds its left neighbour to the entry above
	// that neighbour. The first entry of row n is B(n).
	row := []*big.Int{big.NewInt(1)}
	for range n {
		next := make([]*big.Int, len(row)+1)
		next[0] = row[len(row)-1]
		for j := range row {
			next[j+1] = new(big.Int).Add(next[j], row[j])
		}
		row = next
	}
	return new(big.Int).Set(row[0])
}

// WaysToSum returns the number of ways to write target as a sum of the
// given parts, each usable any number of times, ignoring order: the
// coin-change count. WaysToSum(200, []int{1, 2, 5, 10, 20, 50, 100, 200})
// is Project Euler #31, and with the primes as parts it is #77. Repeated
// parts count once. There is one way to make 0 (use nothing) and none to
// make a negative target. It panics if a part is not positive.
func WaysToSum(target int, parts []int) *big.Int {
	ps := slices.Clone(parts)
	slices.Sort(ps)
	ps = slices.Compact(ps)
	if len(ps) > 0 && ps[0] <= 0 {
		panic("combinatorics: WaysToSum needs positive parts")
	}
	if target < 0 {
		return new(big.Int)
	}
	ways := make([]*big.Int, target+1)
	for i := range ways {
		ways[i] = new(big.Int)
	}
	ways[0].SetInt64(1)
	for _, p := range ps {
		for s := p; s <= target; s++ {
			ways[s].Add(ways[s], ways[s-p])
		}
	}
	return ways[target]
}

// BinomialMod returns C(n, k) mod p for a prime p, using Lucas' theorem,
// so n and k may be far larger than p. It is 0 when k < 0 or k > n. Each
// base-p digit of k costs O(min(digit, p - digit)) multiplications, which
// makes it fast for a small p with any n, and for a large p (such as
// 1_000_000_007) with a small k or n - k; for many binomials below one
// bound, FactorialsMod is faster. It panics if n < 0 or p < 2; the result
// is meaningless if p is not prime.
func BinomialMod[E constraints.Integer](n, k, p E) E {
	if n < 0 {
		panic("combinatorics: BinomialMod with negative n")
	}
	if p < 2 {
		panic("combinatorics: BinomialMod needs a prime modulus")
	}
	if k < 0 || k > n {
		return 0
	}
	un, uk, up := uint64(n), uint64(k), uint64(p)
	r := uint64(1)
	for uk > 0 && r != 0 {
		a, b := un%up, uk%up
		un, uk = un/up, uk/up
		if b > a {
			return 0
		}
		b = min(b, a-b)
		// C(a, b) = Π (a-b+i)/i for i = 1..b, with a < p so every i is
		// invertible.
		num, den := uint64(1), uint64(1)
		for i := uint64(1); i <= b; i++ {
			num = mathx.MulMod(num, a-b+i, up)
			den = mathx.MulMod(den, i, up)
		}
		r = mathx.MulMod(r, mathx.MulMod(num, mathx.PowMod(den, up-2, up), up), up)
	}
	return E(r)
}

// FactorialsMod returns two tables of length n+1 for a prime p > n:
// fact[i] = i! mod p and inv[i] = the inverse of i! mod p. With them a
// binomial coefficient modulo p is three table lookups:
//
//	C(a, b) ≡ fact[a] · inv[b] · inv[a-b]  (mod p), for 0 <= b <= a <= n
//
// It panics if n < 0 or p <= n; the tables are meaningless if p is not
// prime.
func FactorialsMod[E constraints.Integer](n int, p E) (fact, inv []E) {
	if n < 0 {
		panic("combinatorics: FactorialsMod of a negative number")
	}
	if p < 2 || uint64(p) <= uint64(n) {
		panic("combinatorics: FactorialsMod needs a prime modulus above n")
	}
	up := uint64(p)
	fact, inv = make([]E, n+1), make([]E, n+1)
	f := uint64(1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		f = mathx.MulMod(f, uint64(i), up)
		fact[i] = E(f)
	}
	g := mathx.PowMod(f, up-2, up)
	for i := n; i >= 0; i-- {
		inv[i] = E(g)
		g = mathx.MulMod(g, uint64(i), up)
	}
	return fact, inv
}
