package numtheory

import (
	"maps"
	"math"
	"math/big"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/t1sk3/eulerlib/v4/internal/testutil"
	"github.com/t1sk3/eulerlib/v4/primes"
	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/pair"
)

func bruteDivisors(n int) []int {
	out := []int{}
	for d := 1; d <= n; d++ {
		if n%d == 0 {
			out = append(out, d)
		}
	}
	return out
}

func bruteFactorize(n int) map[int]int {
	f := map[int]int{}
	for p := 2; n > 1; p++ {
		for n%p == 0 {
			f[p]++
			n /= p
		}
	}
	return f
}

func TestDivisors(t *testing.T) {
	if got := Divisors(28); !slices.Equal(got, []int{1, 2, 4, 7, 14, 28}) {
		t.Errorf("Divisors(28) = %v", got)
	}
	if got := Divisors(1); !slices.Equal(got, []int{1}) {
		t.Errorf("Divisors(1) = %v", got)
	}
	if len(Divisors(0)) != 0 || len(Divisors(-4)) != 0 {
		t.Error("Divisors of n < 1 should be empty")
	}
	for n := 1; n <= 3000; n++ {
		want := bruteDivisors(n)
		if got := Divisors(n); !slices.Equal(got, want) {
			t.Fatalf("Divisors(%d) = %v, want %v", n, got, want)
		}
		if got := CountDivisors(n); got != len(want) {
			t.Fatalf("CountDivisors(%d) = %d, want %d", n, got, len(want))
		}
		sum := 0
		for _, d := range want[:len(want)-1] {
			sum += d
		}
		if got := SumDivisors(n); got != sum {
			t.Fatalf("SumDivisors(%d) = %d, want %d", n, got, sum)
		}
	}
	if CountDivisors(0) != 0 || SumDivisors(0) != 0 || SumDivisors(1) != 0 {
		t.Error("edge cases")
	}
	// A square near the top of int64, where math.Sqrt is inexact.
	big := int64(3037000499) * 3037000499
	if got := Divisors(big); got[len(got)-1] != big || len(got) != CountDivisors(big) {
		t.Errorf("Divisors(%d) = %v", big, got)
	}
}

func TestSigmaSieve(t *testing.T) {
	s := SigmaSieve(10_000)
	for i := range 10_001 {
		if s[i] != SumDivisors(i) {
			t.Fatalf("SigmaSieve[%d] = %d, want %d", i, s[i], SumDivisors(i))
		}
	}
	if got := SigmaSieve(uint8(255)); got[255] != uint8(SumDivisors(255)) {
		t.Error("SigmaSieve(uint8(255)) should not overflow its loop")
	}
}

func TestPredicates(t *testing.T) {
	if !IsPerfect(28) || !IsPerfect(496) || IsPerfect(27) || IsPerfect(0) {
		t.Error("IsPerfect")
	}
	if !IsAbundant(12) || IsAbundant(11) || !IsDeficient(11) {
		t.Error("IsAbundant/IsDeficient")
	}
	if !IsAmicable(220) || !IsAmicable(284) || IsAmicable(6) {
		t.Error("IsAmicable")
	}
}

func TestAmicablePairs(t *testing.T) {
	got := AmicablePairs(10_000).ToSlice()
	want := []pair.Pair[int, int]{pair.Of(220, 284), pair.Of(1184, 1210), pair.Of(2620, 2924), pair.Of(5020, 5564), pair.Of(6232, 6368)}
	if !slices.Equal(got, want) {
		t.Fatalf("AmicablePairs(10000) = %v", got)
	}
	// Both members must be below the limit.
	if got := AmicablePairs(284).ToSlice(); len(got) != 0 {
		t.Errorf("AmicablePairs(284) = %v, want none (284 is not < 284)", got)
	}
	if got := AmicablePairs(285).ToSlice(); !slices.Equal(got, want[:1]) {
		t.Errorf("AmicablePairs(285) = %v", got)
	}
	testutil.Reusable(t, AmicablePairs(10_000), -1)
	testutil.EarlyExit(t, AmicablePairs(10_000), 2)
}

func TestFactors(t *testing.T) {
	if got := maps.Collect(Factors(360)); !maps.Equal(got, map[int]int{2: 3, 3: 2, 5: 1}) {
		t.Errorf("Factors(360) = %v", got)
	}
	for n := -5; n <= 10_000; n++ {
		want := map[int]int{}
		if n >= 2 {
			want = bruteFactorize(n)
		}
		if got := Factorize(n); !maps.Equal(got, want) {
			t.Fatalf("Factorize(%d) = %v, want %v", n, got, want)
		}
		var prev int
		for p, e := range Factors(n) {
			if p <= prev || e < 1 || !primes.IsPrime(p) {
				t.Fatalf("Factors(%d) not ascending primes: %d^%d after %d", n, p, e, prev)
			}
			prev = p
		}
	}
	// Reusable and early exit.
	seq := Factors(2 * 3 * 5 * 7 * 11)
	if a, b := maps.Collect(seq), maps.Collect(seq); !maps.Equal(a, b) {
		t.Error("Factors not reusable")
	}
	for range seq {
		break
	}
	if got := Factorize(int64(math.MaxInt64)); !maps.Equal(got, map[int64]int{7: 2, 73: 1, 127: 1, 337: 1, 92737: 1, 649657: 1}) {
		t.Errorf("Factorize(MaxInt64) = %v", got)
	}
	if got := Factorize(uint64(18_446_744_073_709_551_557)); !maps.Equal(got, map[uint64]int{18_446_744_073_709_551_557: 1}) {
		t.Errorf("Factorize(largest 64-bit prime) = %v", got)
	}
}

func TestFactorsSPF(t *testing.T) {
	spf := primes.SmallestFactorSieve(10_000)
	for n := 0; n <= 10_000; n++ {
		if got, want := FactorizeSPF(n, spf), Factorize(n); !maps.Equal(got, want) {
			t.Fatalf("FactorizeSPF(%d) = %v, want %v", n, got, want)
		}
		if got, want := slices.Collect(func(y func(int) bool) {
			for p := range FactorsSPF(n, spf) {
				if !y(p) {
					return
				}
			}
		}), DistinctPrimeFactors(n); !slices.Equal(got, want) {
			t.Fatalf("FactorsSPF(%d) order = %v, want %v", n, got, want)
		}
	}
	mustPanic(t, "spf too small", func() { FactorizeSPF(10_001, spf) })
	mustPanic(t, "bad spf", func() { FactorizeSPF(5, make([]int, 10)) })
}

func TestPrimeFactors(t *testing.T) {
	if got := PrimeFactors(360); !slices.Equal(got, []int{2, 2, 2, 3, 3, 5}) {
		t.Errorf("PrimeFactors(360) = %v", got)
	}
	if got := DistinctPrimeFactors(644); !slices.Equal(got, []int{2, 7, 23}) {
		t.Errorf("DistinctPrimeFactors(644) = %v", got)
	}
	if len(PrimeFactors(1)) != 0 || len(DistinctPrimeFactors(0)) != 0 {
		t.Error("edge cases")
	}
}

func TestFactorizeBig(t *testing.T) {
	for n := int64(0); n <= 3000; n++ {
		got := FactorizeBig(big.NewInt(n))
		want := Factorize(n)
		if len(got) != len(want) {
			t.Fatalf("FactorizeBig(%d) = %v, want %v", n, got, want)
		}
		for i, f := range got {
			if !f.First.IsInt64() || want[f.First.Int64()] != f.Second || (i > 0 && got[i-1].First.Cmp(f.First) >= 0) {
				t.Fatalf("FactorizeBig(%d) = %v, want %v", n, got, want)
			}
		}
	}
	// 2^10 * (2^89 - 1), a Mersenne prime cofactor far above 64 bits.
	m89 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 89), big.NewInt(1))
	n := new(big.Int).Lsh(m89, 10)
	got := FactorizeBig(n)
	if len(got) != 2 || got[0].First.Int64() != 2 || got[0].Second != 10 || got[1].First.Cmp(m89) != 0 || got[1].Second != 1 {
		t.Errorf("FactorizeBig(2^10 * M89) = %v", got)
	}
	if n.Cmp(new(big.Int).Lsh(m89, 10)) != 0 {
		t.Error("FactorizeBig modified its argument")
	}
}

func TestGcdLcm(t *testing.T) {
	if Gcd(12, 18) != 6 || Gcd(12, 18, 8) != 2 || Gcd(-12, 18) != 6 || Gcd(0, 5) != 5 || Gcd(7) != 7 || Gcd(-7) != 7 || Gcd[int]() != 0 {
		t.Error("Gcd")
	}
	if Lcm(4, 6) != 12 || Lcm(1, 2, 3, 4, 5, 6, 7, 8, 9, 10) != 2520 || Lcm(-4, 6) != 12 || Lcm(0, 5) != 0 || Lcm[int]() != 1 {
		t.Error("Lcm")
	}
	if !IsCoprime(8, 15) || IsCoprime(8, 12) {
		t.Error("IsCoprime")
	}
	r := rand.New(rand.NewPCG(5, 6))
	for range 5000 {
		a, b := r.Int64N(1_000_000)-500_000, r.Int64N(1_000_000)-500_000
		g, x, y := ExtGcd(a, b)
		if g != Gcd(a, b) || a*x+b*y != g {
			t.Fatalf("ExtGcd(%d, %d) = (%d, %d, %d)", a, b, g, x, y)
		}
	}
}

func TestModInverse(t *testing.T) {
	cases := []struct {
		a, m, want int64
		ok         bool
	}{
		{3, 11, 4, true}, {10, 17, 12, true}, {6, 9, 0, false}, {1, 1, 0, true}, {3, 0, 0, false}, {3, -11, 4, true}, {-3, 11, 7, true},
	}
	for _, tc := range cases {
		got, ok := ModInverse(tc.a, tc.m).Value()
		if ok != tc.ok || (ok && got != tc.want) {
			t.Errorf("ModInverse(%d, %d) = (%d, %v), want (%d, %v)", tc.a, tc.m, got, ok, tc.want, tc.ok)
		}
	}
	// Optional.Map method form (go-utils v0.6).
	if got := ModInverse(int64(3), 11).Map(func(x int64) int64 { return x * 3 % 11 }).Get(); got != 1 {
		t.Errorf("ModInverse(3, 11).Map = %d", got)
	}
}

func TestCRT(t *testing.T) {
	check := func(rs, ms []int64, wantX, wantM int64, wantOK bool) {
		t.Helper()
		res, ok := CRT(rs, ms).Value()
		if ok != wantOK || (ok && (res.First != wantX || res.Second != wantM)) {
			t.Errorf("CRT(%v, %v) = (%v, %v), want (%d, %d, %v)", rs, ms, res, ok, wantX, wantM, wantOK)
		}
	}
	check([]int64{2, 3, 2}, []int64{3, 5, 7}, 23, 105, true)
	check([]int64{5, 3}, []int64{6, 4}, 11, 12, true)
	check([]int64{1, 2}, []int64{4, 4}, 0, 0, false)
	check([]int64{1, 2}, []int64{4}, 0, 0, false)
	check([]int64{1}, []int64{0}, 0, 0, false)
	check([]int64{1, 2}, []int64{3, 0}, 0, 0, false)
	check(nil, nil, 0, 0, false)
	check([]int64{-1}, []int64{5}, 4, 5, true)
	// Moduli whose product is close to MaxInt64: the v3 version overflowed.
	p, q := int64(3_037_000_493), int64(3_037_000_453)
	x := int64(9_000_000_000_000_000_007) % (p * q)
	check([]int64{x % p, x % q}, []int64{p, q}, x, p*q, true)
	// Brute force on small systems.
	r := rand.New(rand.NewPCG(7, 8))
	for range 2000 {
		m1, m2 := r.Int64N(30)+1, r.Int64N(30)+1
		r1, r2 := r.Int64N(m1), r.Int64N(m2)
		l := Lcm(m1, m2)
		want, found := int64(-1), false
		for v := range l {
			if v%m1 == r1 && v%m2 == r2 {
				want, found = v, true
				break
			}
		}
		res, ok := CRT([]int64{r1, r2}, []int64{m1, m2}).Value()
		if ok != found || (ok && (res.First != want || res.Second != l)) {
			t.Fatalf("CRT([%d %d], [%d %d]) = %v %v, want %d %v", r1, r2, m1, m2, res, ok, want, found)
		}
	}
}

func TestPowModMulMod(t *testing.T) {
	r := rand.New(rand.NewPCG(9, 10))
	for range 3000 {
		m := r.Int64N(math.MaxInt64-1) + 1
		if r.IntN(3) == 0 {
			m = math.MaxInt64 - r.Int64N(1000) // near 2^63
		}
		x, y := r.Int64()-math.MaxInt64/2, r.Int64N(1<<40)
		want := new(big.Int).Exp(big.NewInt(x), big.NewInt(y), big.NewInt(m))
		if got := PowMod(x, y, m); got != want.Int64() {
			t.Fatalf("PowMod(%d, %d, %d) = %d, want %s", x, y, m, got, want)
		}
		a, b := r.Int64()-math.MaxInt64/2, r.Int64()
		want = new(big.Int).Mul(big.NewInt(a), big.NewInt(b))
		want.Mod(want, big.NewInt(m))
		if got := MulMod(a, b, m); got != want.Int64() {
			t.Fatalf("MulMod(%d, %d, %d) = %d, want %s", a, b, m, got, want)
		}
	}
	if PowMod(uint64(math.MaxUint64), 2, math.MaxUint64-58) != 3364 {
		t.Error("PowMod uint64")
	}
	if PowMod(2, 10, 1) != 0 || PowMod(3, 0, 7) != 1 {
		t.Error("PowMod edge cases")
	}
	mustPanic(t, "PowMod m=0", func() { PowMod(2, 3, 0) })
	mustPanic(t, "PowMod y<0", func() { PowMod(2, -3, 5) })
}

func TestPowers(t *testing.T) {
	for b := -5; b <= 5; b++ {
		want := 1
		for n := range 12 {
			if got := Pow(b, n); got != want {
				t.Fatalf("Pow(%d, %d) = %d, want %d", b, n, got, want)
			}
			want *= b
		}
	}
	if PowBig(big.NewInt(2), 100).String() != "1267650600228229401496703205376" {
		t.Error("PowBig")
	}
	if f, _ := PowBigFloat(big.NewFloat(1.5), 3).Float64(); f != 3.375 {
		t.Error("PowBigFloat")
	}
	if !IsPowerOfTwo(1) || !IsPowerOfTwo(1024) || IsPowerOfTwo(0) || IsPowerOfTwo(12) || IsPowerOfTwo(-8) {
		t.Error("IsPowerOfTwo")
	}
}

func TestIsSquareIsqrt(t *testing.T) {
	for n := -10; n <= 100_000; n++ {
		r := 0
		for (r+1)*(r+1) <= n {
			r++
		}
		if got, want := IsSquare(n), n >= 0 && r*r == n; got != want {
			t.Fatalf("IsSquare(%d) = %v", n, got)
		}
		if n >= 0 && Isqrt(n) != r {
			t.Fatalf("Isqrt(%d) = %d, want %d", n, Isqrt(n), r)
		}
	}
	// Near 2^63 and 2^64, where math.Sqrt rounds.
	s := int64(3_037_000_499)
	if !IsSquare(s*s) || IsSquare(s*s+1) || IsSquare(s*s-1) || IsSquare(int64(math.MaxInt64)) {
		t.Error("IsSquare near 2^63")
	}
	u := uint64(4_294_967_295)
	if !IsSquare(u*u) || IsSquare(u*u+1) || IsSquare(uint64(math.MaxUint64)) || Isqrt(uint64(math.MaxUint64)) != u {
		t.Error("IsSquare near 2^64")
	}
	if !FloatIsInteger(3.0) || FloatIsInteger(3.5) || FloatIsInteger(math.Inf(1)) || FloatIsInteger(math.NaN()) || !FloatIsInteger(float32(-2)) {
		t.Error("FloatIsInteger")
	}
}

func TestTotientMobius(t *testing.T) {
	ts := TotientSieve(5000)
	for n := 0; n <= 5000; n++ {
		want := 0
		for k := 1; k <= n; k++ {
			if Gcd(k, n) == 1 {
				want++
			}
		}
		if Totient(n) != want || ts[n] != want {
			t.Fatalf("φ(%d): Totient %d, sieve %d, want %d", n, Totient(n), ts[n], want)
		}
	}
	mu := MobiusSieve(5000)
	for n := 0; n <= 5000; n++ {
		want := 0
		if n >= 1 {
			want = 1
			for _, e := range Factors(n) {
				if e > 1 {
					want = 0
					break
				}
				want = -want
			}
		}
		if mu[n] != want {
			t.Fatalf("μ(%d) = %d, want %d", n, mu[n], want)
		}
	}
	if got := TotientSieve(uint8(255)); got[255] != 128 {
		t.Errorf("TotientSieve(uint8(255))[255] = %d", got[255])
	}
	if got := MobiusSieve(int8(127)); got[127] != -1 {
		t.Errorf("MobiusSieve(int8(127))[127] = %d", got[127])
	}
}

func TestMisc(t *testing.T) {
	if Abs(-5) != 5 || Abs(5) != 5 {
		t.Error("Abs")
	}
	if SumOfSquares(10) != 385 || SquareOfSum(10) != 3025 {
		t.Error("SumOfSquares/SquareOfSum")
	}
	if math.Abs(ToRadians(180)-math.Pi) > 1e-15 {
		t.Error("ToRadians")
	}
}

func TestPredicateBuilders(t *testing.T) {
	if got := stream.Sum(stream.Range(1, 1000).Filter(DivisibleByAny(3, 5))); got != 233168 {
		t.Errorf("Euler #1 = %d", got)
	}
	if got := stream.Range(1, 20).Filter(DivisibleBy(7)).ToSlice(); !slices.Equal(got, []int{7, 14}) {
		t.Errorf("DivisibleBy(7) = %v", got)
	}
	if got := stream.RangeClosed(1, 10).Filter(CoprimeTo(10)).ToSlice(); !slices.Equal(got, []int{1, 3, 7, 9}) {
		t.Errorf("CoprimeTo(10) = %v", got)
	}
	if got := slices.IndexFunc([]int{4, 9, 10, 21}, DivisibleBy(5)); got != 2 {
		t.Errorf("IndexFunc(DivisibleBy(5)) = %d", got)
	}
	mustPanic(t, "DivisibleBy(0)", func() { DivisibleBy(0) })
	mustPanic(t, "DivisibleByAny(0)", func() { DivisibleByAny(3, 0) })
}

func TestCollectors(t *testing.T) {
	if got := stream.RangeClosed(1, 20).Collect(LcmOf[int]()); got != 232792560 {
		t.Errorf("Euler #5 = %d", got)
	}
	if got := stream.Of(12, -18, 30).Collect(GcdOf[int]()); got != 6 {
		t.Errorf("GcdOf = %d", got)
	}
	if got := stream.RangeClosed(1, 10).Collect(Product[int]()); got != 3628800 {
		t.Errorf("Product = %d", got)
	}
	if got := stream.Of(0.5, 4.0).Collect(Product[float64]()); got != 2 {
		t.Errorf("Product float = %v", got)
	}
	// Empty streams give the identity element.
	if stream.Empty[int]().Collect(Product[int]()) != 1 || stream.Empty[int]().Collect(LcmOf[int]()) != 1 || stream.Empty[int]().Collect(GcdOf[int]()) != 0 {
		t.Error("collector identities")
	}
}

func mustPanic(t *testing.T, name string, f func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s: expected a panic", name)
		}
	}()
	f()
}

func TestLargeInputs(t *testing.T) {
	// The product of two 32-bit primes: trial division would need ~4e9 steps.
	p, q := uint64(4_294_967_291), uint64(4_294_967_279)
	if got := Divisors(p * q); !slices.Equal(got, []uint64{1, q, p, p * q}) {
		t.Errorf("Divisors(p*q) = %v", got)
	}
	if got := SumDivisors(p * q); got != 1+p+q {
		t.Errorf("SumDivisors(p*q) = %d", got)
	}
	// 2^62 has 63 divisors; σ(2^62) - 2^62 = 2^62 - 1.
	n := int64(1) << 62
	if CountDivisors(n) != 63 || SumDivisors(n) != n-1 || len(Divisors(n)) != 63 {
		t.Error("2^62")
	}
	// A highly composite number: compare against summing Divisors.
	hc := int64(963761198400) // 6720 divisors
	ds := Divisors(hc)
	var sum int64
	for _, d := range ds[:len(ds)-1] {
		sum += d
	}
	if len(ds) != 6720 || CountDivisors(hc) != 6720 || SumDivisors(hc) != sum || !slices.IsSorted(ds) {
		t.Errorf("highly composite: %d divisors, sum %d vs %d", len(ds), SumDivisors(hc), sum)
	}
}
