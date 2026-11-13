package primes

import (
	"math/big"
	"testing"

	"github.com/t1sk3/go-utils/stream"
)

func TestPiAndSumUpTo(t *testing.T) {
	// Every n up to a bound, against the sieve.
	sv := NewSieve(20000)
	count, sum := 0, int64(0)
	for n := -2; n <= 20000; n++ {
		if n >= 0 && sv.Has(n) {
			count++
			sum += int64(n)
		}
		if got := Pi(n); got != count {
			t.Fatalf("Pi(%d) = %d, want %d", n, got, count)
		}
		if got := SumUpTo(n); got.Int64() != sum {
			t.Fatalf("SumUpTo(%d) = %s, want %d", n, got, sum)
		}
	}
	for n, want := range map[uint64]uint64{1_000_000: 78498, 1_000_000_000: 50847534, 10_000_000_000: 455052511} {
		if testing.Short() && n > 1_000_000_000 {
			continue
		}
		if got := Pi(n); got != want {
			t.Errorf("Pi(%d) = %d, want %d", n, got, want)
		}
	}
	if got := SumUpTo(2_000_000).String(); got != "142913828922" {
		t.Errorf("SumUpTo(2e6) = %s (Euler #10)", got)
	}
	const n = 30_000_000
	want := new(big.Int).SetUint64(stream.Sum(Between[uint64](2, n)))
	if got := SumUpTo(n); got.Cmp(want) != 0 {
		t.Errorf("SumUpTo(%d) = %s, want %s", n, got, want)
	}
	if got := Pi(uint8(255)); got != 54 {
		t.Errorf("Pi[uint8](255) = %d", got)
	}
	if got := SumUpTo(int8(127)).Int64(); got != 1720 {
		t.Errorf("SumUpTo[int8](127) = %d", got)
	}
}

// The 128-bit arithmetic behind SumUpTo.
func TestU128(t *testing.T) {
	toBig := func(a u128) *big.Int {
		b := new(big.Int).SetUint64(a.hi)
		return b.Lsh(b, 64).Add(b, new(big.Int).SetUint64(a.lo))
	}
	for _, v := range []uint64{1, 2, 3, 10, 1 << 32, 1<<63 - 1, 1 << 63, 1<<64 - 2, 1<<64 - 1} {
		want := new(big.Int).SetUint64(v)
		want.Mul(want, new(big.Int).Add(want, big.NewInt(1))).Rsh(want, 1).Sub(want, big.NewInt(1))
		if got := toBig(sumFrom2(v)); got.Cmp(want) != 0 {
			t.Errorf("sumFrom2(%d) = %s, want %s", v, got, want)
		}
	}
	a, b := u128{5, 3}, u128{2, 9}
	if got := a.sub(b); got != (u128{2, 1<<64 - 6}) {
		t.Errorf("sub = %v", got)
	}
	if got := (u128{1, 1 << 63}).mul(6); got != (u128{9, 0}) {
		t.Errorf("mul = %v", got)
	}
}

func TestNth(t *testing.T) {
	if Nth[int](1) != 2 || Nth[int](6) != 13 || Nth[int](10_001) != 104743 {
		t.Error("Nth")
	}
	if Nth[uint8](54) != 251 {
		t.Error("Nth[uint8](54)")
	}
	for name, f := range map[string]func(){
		"Nth(0)":         func() { Nth[int](0) },
		"Nth[uint8](55)": func() { Nth[uint8](55) },
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s should panic", name)
				}
			}()
			f()
		}()
	}
}

func TestIsPrimeBig(t *testing.T) {
	for n := int64(-3); n < 2000; n++ {
		if IsPrimeBig(big.NewInt(n)) != IsPrime(n) {
			t.Fatalf("IsPrimeBig(%d)", n)
		}
	}
	m127 := new(big.Int).Lsh(big.NewInt(1), 127)
	m127.Sub(m127, big.NewInt(1))
	if !IsPrimeBig(m127) || IsPrimeBig(new(big.Int).Mul(m127, m127)) {
		t.Error("IsPrimeBig around 2^127-1")
	}
}

func TestCircularAndTruncatable(t *testing.T) {
	// Euler #35: 55 circular primes below one million.
	if got := Between(2, 999_999).Filter(IsCircular[int]).Count(); got != 55 {
		t.Errorf("Euler #35 = %d, want 55", got)
	}
	if !IsCircular(197) || IsCircular(19) || IsCircular(101) || IsCircular(1) || IsCircular(-7) {
		t.Error("IsCircular")
	}
	// Euler #37: the eleven truncatable primes above 7 sum to 748317.
	n, sum := 0, 0
	for p := range Between(10, 1_000_000).Filter(IsTruncatable[int]).All() {
		n++
		sum += p
	}
	if n != 11 || sum != 748317 {
		t.Errorf("Euler #37: %d primes summing to %d", n, sum)
	}
	if !IsTruncatable(3797) || !IsTruncatable(7) || IsTruncatable(1) || IsTruncatable(-3797) {
		t.Error("IsTruncatable")
	}
	// 103 is prime, and so is 3, but "03" is not a truncation.
	if IsLeftTruncatable(103) || !IsLeftTruncatable(113) || IsLeftTruncatable(19) {
		t.Error("IsLeftTruncatable")
	}
	if !IsRightTruncatable(7393) || IsRightTruncatable(19) {
		t.Error("IsRightTruncatable")
	}
	// The largest left-truncatable prime has 24 digits; the largest that
	// fits here is checked without overflowing.
	if !IsLeftTruncatable(uint64(6_918_997_653_319_693_967)) {
		t.Error("IsLeftTruncatable(6918997653319693967)")
	}
}
