package combinatorics

import (
	"math"
	"math/big"
	"testing"
)

// Binomial used to convert n and k to int64, returning 0 for a uint64 n
// above MaxInt64.
func TestBinomialLargeUint64(t *testing.T) {
	n := uint64(math.MaxUint64)
	bn := new(big.Int).SetUint64(n)
	if got := Binomial(n, 1); got.Cmp(bn) != 0 {
		t.Errorf("Binomial(MaxUint64, 1) = %s", got)
	}
	if got := Binomial(n, n-1); got.Cmp(bn) != 0 {
		t.Errorf("Binomial(MaxUint64, MaxUint64-1) = %s", got)
	}
	if got := Binomial(n, n); got.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Binomial(MaxUint64, MaxUint64) = %s", got)
	}
	// C(n, 3) = n(n-1)(n-2)/6
	want := new(big.Int).Mul(bn, new(big.Int).SetUint64(n-1))
	want.Mul(want, new(big.Int).SetUint64(n-2)).Quo(want, big.NewInt(6))
	if got := Binomial(n, 3); got.Cmp(want) != 0 {
		t.Errorf("Binomial(MaxUint64, 3) = %s, want %s", got, want)
	}
	m := uint64(1<<63 + 10)
	if got := Binomial(m, m-1); got.Cmp(new(big.Int).SetUint64(m)) != 0 {
		t.Errorf("Binomial(2^63+10, 2^63+9) = %s", got)
	}
	// Both paths agree at the int64 boundary (Pascal's rule).
	lo, hi := uint64(math.MaxInt64), uint64(math.MaxInt64)+1
	sum := new(big.Int).Add(Binomial(lo, 4), Binomial(lo, 5))
	if got := Binomial(hi, 5); got.Cmp(sum) != 0 {
		t.Errorf("Binomial(2^63, 5) = %s, want %s", got, sum)
	}
}
