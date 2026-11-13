package numtheory

import (
	"math"
	"math/big"
	"testing"
)

// SumOfSquares must be exact whenever the result fits, even though
// n(n+1)(2n+1) itself overflows long before that.
func TestSumOfSquaresNoIntermediateOverflow(t *testing.T) {
	want := func(n int64) *big.Int {
		r := new(big.Int).Mul(big.NewInt(n), big.NewInt(n+1))
		return r.Mul(r, big.NewInt(2*n+1)).Quo(r, big.NewInt(6))
	}
	for _, n := range []int64{0, 1, 2, 3, 4, 5, 6, 7, 100, 1_664_510, 1_664_511, 2_000_000, 2_000_001, 2_000_002, 3_024_616} {
		w := want(n)
		if !w.IsInt64() {
			t.Fatalf("test value n=%d does not fit", n)
		}
		if got := SumOfSquares(n); got != w.Int64() {
			t.Errorf("SumOfSquares(%d) = %d, want %s", n, got, w)
		}
	}
	// Every n whose result fits in the type, for the small types.
	for n := int64(0); want(n).Int64() <= math.MaxInt16; n++ {
		if got := SumOfSquares(int16(n)); int64(got) != want(n).Int64() {
			t.Errorf("SumOfSquares[int16](%d) = %d, want %s", n, got, want(n))
		}
	}
	for n := int64(0); want(n).Int64() <= math.MaxUint8; n++ {
		if got := SumOfSquares(uint8(n)); int64(got) != want(n).Int64() {
			t.Errorf("SumOfSquares[uint8](%d) = %d, want %s", n, got, want(n))
		}
	}
	if SumOfSquares(-3) != 0 {
		t.Error("SumOfSquares of a negative number should be 0")
	}
}

// The abundant/deficient/perfect tests must not depend on the divisor sum
// fitting in E.
func TestAbundanceIsExact(t *testing.T) {
	for n := -5; n <= 0; n++ {
		if IsAbundant(n) || IsDeficient(n) || IsPerfect(n) {
			t.Errorf("n = %d should be neither abundant, deficient nor perfect", n)
		}
	}
	for n := 1; n <= 255; n++ {
		s := 0
		for d := 1; d < n; d++ {
			if n%d == 0 {
				s += d
			}
		}
		a, d, p := s > n, s < n, s == n
		if IsAbundant(uint8(n)) != a || IsDeficient(uint8(n)) != d || IsPerfect(uint8(n)) != p {
			t.Errorf("uint8 %d: abundant/deficient/perfect = %v/%v/%v, want %v/%v/%v",
				n, IsAbundant(uint8(n)), IsDeficient(uint8(n)), IsPerfect(uint8(n)), a, d, p)
		}
		if n <= 127 && (IsAbundant(int8(n)) != a || IsDeficient(int8(n)) != d || IsPerfect(int8(n)) != p) {
			t.Errorf("int8 %d: wrong classification", n)
		}
	}
	// σ(n) - n does not fit in int64, and σ(n) does not fit in uint64.
	if n := int64(3) << 61; !IsAbundant(n) || IsDeficient(n) {
		t.Errorf("3·2^61 is abundant")
	}
	if n := uint64(3) << 62; !IsAbundant(n) || IsDeficient(n) || IsPerfect(n) {
		t.Errorf("3·2^62 is abundant")
	}
	if n := uint64(18446744073709551557); IsAbundant(n) || !IsDeficient(n) { // prime
		t.Errorf("a prime is deficient")
	}
	if n := uint64(2305843008139952128); !IsPerfect(n) || IsAbundant(n) || IsDeficient(n) { // 2^30(2^31-1)
		t.Errorf("2^30(2^31-1) is perfect")
	}
}

// CRT must report a modulus that does not fit instead of wrapping.
func TestCRTModulusOverflow(t *testing.T) {
	p1, p2 := int64(4294967311), int64(4294967357) // p1*p2 > 2^63
	if res, ok := CRT([]int64{5, 7}, []int64{p1, p2}).Value(); ok {
		t.Errorf("CRT with an overflowing modulus = %v, want empty", res)
	}
	if _, ok := CRT([]int8{1, 2}, []int8{16, 15}).Value(); ok { // lcm 240
		t.Error("CRT[int8] with lcm 240 should be empty")
	}
	// The largest modulus that still fits keeps working.
	res, ok := CRT([]int8{1, 2}, []int8{9, 14}).Value() // lcm 126
	if !ok || res.Second != 126 || res.First%9 != 1 || res.First%14 != 2 {
		t.Errorf("CRT[int8]([1 2], [9 14]) = %v %v", res, ok)
	}
	q1, q2 := int64(3037000493), int64(3037000453) // product just below 2^63
	res64, ok := CRT([]int64{5, 7}, []int64{q1, q2}).Value()
	if !ok || res64.Second != q1*q2 || res64.First%q1 != 5 || res64.First%q2 != 7 {
		t.Errorf("CRT near MaxInt64 = %v %v", res64, ok)
	}
}
