package mathx

import (
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// isqrtOK reports whether r == floor(sqrt(n)), comparing in 128 bits so
// nothing overflows.
func isqrtOK(n, r uint64) bool {
	bn := new(big.Int).SetUint64(n)
	br := new(big.Int).SetUint64(r)
	if new(big.Int).Mul(br, br).Cmp(bn) > 0 {
		return false
	}
	br.Add(br, big.NewInt(1))
	return new(big.Int).Mul(br, br).Cmp(bn) > 0
}

func TestIsqrt(t *testing.T) {
	cases := map[uint64]uint64{
		0: 0, 1: 1, 2: 1, 3: 1, 4: 2, 15: 3, 16: 4, 17: 4,
		math.MaxUint32: 65535,
		1 << 52:        1 << 26,
		(1<<53 + 1):    94906265,
		math.MaxUint64: math.MaxUint32,
		// (2^32-1)^2 and its neighbours
		18446744065119617025: 4294967295,
		18446744065119617024: 4294967294,
	}
	for n, want := range cases {
		if got := Isqrt(n); got != want {
			t.Errorf("Isqrt(%d) = %d, want %d", n, got, want)
		}
	}
	// Squares and their neighbours near the float64 precision limit.
	for _, r := range []uint64{1 << 26, 94906265, 1<<32 - 1, 3037000499, 3037000500} {
		for _, n := range []uint64{r*r - 1, r * r, r*r + 1} {
			if got := Isqrt(n); !isqrtOK(n, got) {
				t.Errorf("Isqrt(%d) = %d is not floor(sqrt)", n, got)
			}
		}
	}
}

func FuzzIsqrt(f *testing.F) {
	for _, seed := range []uint64{0, 1, 2, 1 << 53, 1<<53 + 1, math.MaxUint64, math.MaxUint64 - 1} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, n uint64) {
		if r := Isqrt(n); !isqrtOK(n, r) {
			t.Fatalf("Isqrt(%d) = %d is not floor(sqrt)", n, r)
		}
	})
}

func TestMulModAgainstBig(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 10_000 {
		a, b, m := r.Uint64(), r.Uint64(), r.Uint64()|1
		want := new(big.Int).Mul(new(big.Int).SetUint64(a), new(big.Int).SetUint64(b))
		want.Mod(want, new(big.Int).SetUint64(m))
		if got := MulMod(a, b, m); got != want.Uint64() {
			t.Fatalf("MulMod(%d, %d, %d) = %d, want %s", a, b, m, got, want)
		}
	}
}

func TestPowModAgainstBig(t *testing.T) {
	r := rand.New(rand.NewPCG(3, 4))
	moduli := []uint64{1, 2, 3, 1_000_000_007, 1<<63 - 25, 1<<63 + 29, math.MaxUint64}
	for range 2_000 {
		moduli = append(moduli, r.Uint64()|1)
	}
	for _, m := range moduli {
		b, e := r.Uint64(), r.Uint64()
		want := new(big.Int).Exp(new(big.Int).SetUint64(b), new(big.Int).SetUint64(e), new(big.Int).SetUint64(m))
		if got := PowMod(b, e, m); got != want.Uint64() {
			t.Fatalf("PowMod(%d, %d, %d) = %d, want %s", b, e, m, got, want)
		}
	}
	if got := PowMod(5, 0, 1); got != 0 {
		t.Errorf("PowMod(5, 0, 1) = %d, want 0", got)
	}
	if got := PowMod(0, 0, 7); got != 1 {
		t.Errorf("PowMod(0, 0, 7) = %d, want 1", got)
	}
}

func TestLimits(t *testing.T) {
	if MaxOf[int8]() != math.MaxInt8 || MinOf[int8]() != math.MinInt8 {
		t.Error("int8 limits")
	}
	if MaxOf[int16]() != math.MaxInt16 || MinOf[int16]() != math.MinInt16 {
		t.Error("int16 limits")
	}
	if MaxOf[int32]() != math.MaxInt32 || MinOf[int32]() != math.MinInt32 {
		t.Error("int32 limits")
	}
	if MaxOf[int64]() != math.MaxInt64 || MinOf[int64]() != math.MinInt64 {
		t.Error("int64 limits")
	}
	if MaxOf[int]() != math.MaxInt || MinOf[int]() != math.MinInt {
		t.Error("int limits")
	}
	if MaxOf[uint8]() != math.MaxUint8 || MaxOf[uint16]() != math.MaxUint16 ||
		MaxOf[uint32]() != math.MaxUint32 || MaxOf[uint64]() != math.MaxUint64 || MaxOf[uint]() != math.MaxUint {
		t.Error("unsigned max")
	}
	if MinOf[uint8]() != 0 || MinOf[uint64]() != 0 {
		t.Error("unsigned min")
	}
	type myInt int32
	if MaxOf[myInt]() != math.MaxInt32 || !IsSigned[myInt]() || IsSigned[uint]() {
		t.Error("named types")
	}
}

func TestAbs64(t *testing.T) {
	if Abs64(int64(math.MinInt64)) != 1<<63 {
		t.Error("Abs64(MinInt64)")
	}
	if Abs64(int8(-128)) != 128 || Abs64(int8(5)) != 5 || Abs64(uint64(math.MaxUint64)) != math.MaxUint64 {
		t.Error("Abs64")
	}
}

func TestIsPrimeAgainstBig(t *testing.T) {
	r := rand.New(rand.NewPCG(11, 12))
	for n := range uint64(100_000) {
		if IsPrime(n) != new(big.Int).SetUint64(n).ProbablyPrime(20) {
			t.Fatalf("IsPrime(%d) wrong", n)
		}
	}
	for range 20_000 {
		n := r.Uint64() | 1
		if IsPrime(n) != new(big.Int).SetUint64(n).ProbablyPrime(20) {
			t.Fatalf("IsPrime(%d) wrong", n)
		}
	}
}

func TestFactor(t *testing.T) {
	check := func(n uint64) {
		t.Helper()
		fs := Factor(n)
		prod := uint64(1)
		for i, p := range fs {
			if !IsPrime(p) || (i > 0 && fs[i-1] > p) {
				t.Fatalf("Factor(%d) = %v: not ascending primes", n, fs)
			}
			prod *= p
		}
		if (n >= 2 && prod != n) || (n < 2 && len(fs) != 0) {
			t.Fatalf("Factor(%d) = %v: product %d", n, fs, prod)
		}
	}
	for n := range uint64(20_000) {
		check(n)
	}
	for _, n := range []uint64{
		math.MaxUint64,
		18_446_744_073_709_551_557,    // prime
		4_294_967_291 * 4_294_967_279, // two 32-bit primes
		1_000_000_007 * 998_244_353,   // two 30-bit primes
		999_983 * 999_983 * 17,        // square of a prime times a small one
		1 << 63,
		3_825_123_056_546_413_051, // strong pseudoprime to many bases
	} {
		check(n)
	}
	r := rand.New(rand.NewPCG(13, 14))
	for range 2_000 {
		check(r.Uint64())
	}
}
