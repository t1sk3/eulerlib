package primes

import (
	"math"
	"math/big"
	"slices"
	"testing"

	"github.com/t1sk3/eulerlib/v4/internal/testutil"
	"github.com/t1sk3/go-utils/stream"
)

// trialDivision is the brute-force reference everything is checked against.
func trialDivision(n int) bool {
	if n < 2 {
		return false
	}
	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return false
		}
	}
	return true
}

func referencePrimes(lo, hi int) []int {
	out := []int{}
	for n := lo; n <= hi; n++ {
		if trialDivision(n) {
			out = append(out, n)
		}
	}
	return out
}

func TestIsPrimeSmall(t *testing.T) {
	for n := -100; n <= 200_000; n++ {
		if got, want := IsPrime(n), trialDivision(n); got != want {
			t.Fatalf("IsPrime(%d) = %v, want %v", n, got, want)
		}
	}
}

func TestIsPrimeLarge(t *testing.T) {
	primes := []uint64{
		1_000_000_007, 2_147_483_647, 4_294_967_291, 1_000_000_000_000_000_003,
		9_223_372_036_854_775_783,  // largest prime < 2^63
		18_446_744_073_709_551_557, // largest prime < 2^64
	}
	composites := []uint64{
		561, 41041, 825265, // Carmichael numbers
		3_215_031_751,                               // strong pseudoprime to bases 2, 3, 5, 7
		3_825_123_056_546_413_051,                   // strong pseudoprime to bases 2..23
		318_665_857_834_031_151_167_461 % (1 << 63), // arbitrary large value, checked against big below
		4_294_967_297,                               // 641 * 6700417
		18_446_744_073_709_551_615,                  // 2^64-1
		18_446_744_073_709_551_557 - 2,
		4_294_967_291 * 4_294_967_279,
	}
	for _, p := range primes {
		if !IsPrime(p) {
			t.Errorf("IsPrime(%d) = false, want true", p)
		}
	}
	for _, c := range composites {
		want := new(big.Int).SetUint64(c).ProbablyPrime(30)
		if got := IsPrime(c); got != want {
			t.Errorf("IsPrime(%d) = %v, want %v", c, got, want)
		}
	}
	if !IsPrime(int64(9_223_372_036_854_775_783)) || IsPrime(int64(math.MaxInt64)) || IsPrime(int64(math.MinInt64)) {
		t.Error("IsPrime on int64 extremes")
	}
}

func TestIsPrimeAgainstBig(t *testing.T) {
	// Every odd number in windows near 2^32, 2^53 and 2^64.
	for _, start := range []uint64{1<<32 - 1001, 1<<53 - 1001, math.MaxUint64 - 2001} {
		for n := start; n < start+2000; n += 2 {
			if got, want := IsPrime(n), new(big.Int).SetUint64(n).ProbablyPrime(30); got != want {
				t.Fatalf("IsPrime(%d) = %v, want %v", n, got, want)
			}
		}
	}
}

func TestNextPrevPrime(t *testing.T) {
	ref := referencePrimes(0, 10_000)
	for n := -5; n < 9_900; n++ {
		i, _ := slices.BinarySearch(ref, n+1)
		if got := NextPrime(n); got != ref[i] {
			t.Fatalf("NextPrime(%d) = %d, want %d", n, got, ref[i])
		}
		j, _ := slices.BinarySearch(ref, n)
		prev := PrevPrime(n)
		if j == 0 {
			if prev.IsPresent() {
				t.Fatalf("PrevPrime(%d) = %v, want empty", n, prev)
			}
		} else if prev.Get() != ref[j-1] {
			t.Fatalf("PrevPrime(%d) = %v, want %d", n, prev, ref[j-1])
		}
	}
	if got := NextPrime(uint64(18_446_744_073_709_551_556)); got != 18_446_744_073_709_551_557 {
		t.Errorf("NextPrime near 2^64 = %d", got)
	}
	if got := PrevPrime(uint64(math.MaxUint64)).Get(); got != 18_446_744_073_709_551_557 {
		t.Errorf("PrevPrime(MaxUint64) = %d", got)
	}
	mustPanic(t, "NextPrime(uint8(251))", func() { NextPrime(uint8(251)) })
	mustPanic(t, "NextPrime(int8(127))", func() { NextPrime(int8(127)) })
	mustPanic(t, "NextPrime(largest 64-bit prime)", func() { NextPrime(uint64(18_446_744_073_709_551_557)) })
}

func TestSieve(t *testing.T) {
	for _, limit := range []int{0, 1, 2, 3, 4, 63, 64, 127, 128, 129, 1000, 10_000} {
		sv := NewSieve(limit)
		want := referencePrimes(0, limit)
		if sv.Limit() != limit {
			t.Fatalf("Limit() = %d", sv.Limit())
		}
		if sv.Count() != len(want) {
			t.Fatalf("NewSieve(%d).Count() = %d, want %d", limit, sv.Count(), len(want))
		}
		if got := slices.Collect(sv.All()); !slices.Equal(got, want) {
			t.Fatalf("NewSieve(%d).All() = %v, want %v", limit, got, want)
		}
		for n := -3; n <= limit; n++ {
			if sv.Has(n) != trialDivision(n) {
				t.Fatalf("NewSieve(%d).Has(%d) wrong", limit, n)
			}
		}
	}
	sv := NewSieve(100)
	mustPanic(t, "Has beyond limit", func() { sv.Has(101) })
	if !sv.IsPrime(101) || sv.IsPrime(1001) {
		t.Error("Sieve.IsPrime beyond the limit should fall back to IsPrime")
	}
	if got := stream.Range(0, 30).Filter(sv.IsPrime).ToSlice(); !slices.Equal(got, []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}) {
		t.Errorf("Filter(sv.IsPrime) = %v", got)
	}
	testutil.Reusable(t, sv.Stream(), -1)
	testutil.EarlyExit(t, sv.Stream(), 5)
	mustPanic(t, "NewSieve(-1)", func() { NewSieve(-1) })
	if NewSieve(uint8(255)).Count() != 54 {
		t.Error("NewSieve(uint8(255)).Count() != 54")
	}
}

func TestStreams(t *testing.T) {
	ref := referencePrimes(0, 200_000)
	if got := All[int]().Limit(len(ref)).ToSlice(); !slices.Equal(got, ref) {
		t.Fatal("All() disagrees with trial division")
	}
	if got := UpTo(200_000); !slices.Equal(got, ref) {
		t.Fatal("UpTo disagrees with trial division")
	}
	if got := First[int](len(ref)); !slices.Equal(got, ref) {
		t.Fatal("First disagrees with trial division")
	}
	for _, r := range [][2]int{{-10, 1}, {-10, 2}, {0, 0}, {2, 2}, {3, 3}, {4, 4}, {10, 5}, {90, 100}, {1000, 3000}, {131_000, 200_000}} {
		want := referencePrimes(r[0], r[1])
		if got := Between(r[0], r[1]).ToSlice(); !slices.Equal(got, want) {
			t.Fatalf("Between(%d, %d) = %v, want %v", r[0], r[1], got, want)
		}
		if got := Count(r[0], r[1]); got != len(want) {
			t.Fatalf("Count(%d, %d) = %d, want %d", r[0], r[1], got, len(want))
		}
		sum := 0
		for _, p := range want {
			sum += p
		}
		if got := Sum(r[0], r[1]); got != sum {
			t.Fatalf("Sum(%d, %d) = %d, want %d", r[0], r[1], got, sum)
		}
	}
	if got := From(100).Limit(3).ToSlice(); !slices.Equal(got, []int{101, 103, 107}) {
		t.Errorf("From(100) = %v", got)
	}
	if got := From(101).Limit(1).ToSlice(); !slices.Equal(got, []int{101}) {
		t.Errorf("From(101) = %v", got)
	}
	if got := All[int]().TakeWhile(stream.LessThan(100)).Count(); got != 25 {
		t.Errorf("primes below 100 = %d, want 25", got)
	}
	if len(First[int](0)) != 0 || len(UpTo(1)) != 0 {
		t.Error("empty results")
	}
}

func TestStreamsStopAtTypeLimit(t *testing.T) {
	got := All[uint8]().ToSlice()
	if len(got) != 54 || got[len(got)-1] != 251 {
		t.Errorf("All[uint8]() has %d primes ending at %d, want 54 ending at 251", len(got), got[len(got)-1])
	}
	if got := All[int8]().ToSlice(); got[len(got)-1] != 127 {
		t.Errorf("All[int8]() ends at %d, want 127", got[len(got)-1])
	}
	if got := len(First[uint8](1000)); got != 54 {
		t.Errorf("First[uint8](1000) has %d primes, want 54", got)
	}
	if got := From(uint64(18_446_744_073_709_551_500)).ToSlice(); !slices.Equal(got, []uint64{18_446_744_073_709_551_521, 18_446_744_073_709_551_533, 18_446_744_073_709_551_557}) {
		t.Errorf("primes from 2^64-116 = %v", got)
	}
	if got := From(int64(math.MaxInt64 - 100)).ToSlice(); !slices.Equal(got, []int64{9_223_372_036_854_775_783}) {
		t.Errorf("primes from MaxInt64-100 = %v", got)
	}
}

func TestStreamsReusableAndEarlyExit(t *testing.T) {
	for name, s := range map[string]stream.Stream[int]{
		"All":     All[int](),
		"From":    From(1000),
		"Between": Between(10, 10_000),
	} {
		t.Run(name, func(t *testing.T) {
			testutil.Reusable(t, s, 100)
			testutil.EarlyExit(t, s, 10)
		})
	}
}

func TestSmallestFactorSieve(t *testing.T) {
	spf := SmallestFactorSieve(10_000)
	if len(spf) != 10_001 || spf[0] != 0 || spf[1] != 0 {
		t.Fatal("SmallestFactorSieve shape")
	}
	for n := 2; n <= 10_000; n++ {
		want := n
		for d := 2; d*d <= n; d++ {
			if n%d == 0 {
				want = d
				break
			}
		}
		if spf[n] != want {
			t.Fatalf("spf[%d] = %d, want %d", n, spf[n], want)
		}
	}
	if got := SmallestFactorSieve(uint8(255)); got[255] != 3 || got[251] != 251 {
		t.Error("SmallestFactorSieve(uint8(255)) should not overflow")
	}
	if got := SmallestFactorSieve(1); len(got) != 2 {
		t.Error("SmallestFactorSieve(1)")
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

func BenchmarkAllFirst1M(b *testing.B) {
	for b.Loop() {
		All[int]().Limit(1_000_000).Count()
	}
}

func BenchmarkBetween2M(b *testing.B) {
	for b.Loop() {
		Sum(2, 2_000_000)
	}
}

func BenchmarkIsPrime64bit(b *testing.B) {
	for b.Loop() {
		IsPrime(uint64(18_446_744_073_709_551_557))
	}
}

func BenchmarkNewSieve10M(b *testing.B) {
	for b.Loop() {
		NewSieve(10_000_000)
	}
}

// referenceBetween lists the primes in [lo, hi] by testing each number with
// isPrime64 (Miller–Rabin), which is checked against math/big above and
// shares no code with the sieve.
func referenceBetween(lo, hi uint64) []uint64 {
	out := []uint64{}
	for n := lo; ; n++ {
		if isPrime64(n) {
			out = append(out, n)
		}
		if n == hi {
			return out
		}
	}
}

func TestSegmentedLargeRanges(t *testing.T) {
	ranges := [][2]uint64{
		// spans three segments
		{1_000_000_000_000, 1_000_000_000_000 + 3*2*segOdds + 13},
		// crosses the point where the sieve switches to Miller–Rabin
		{maxBase*maxBase - 3000, maxBase*maxBase + 3000},
		{math.MaxUint64 - 5000, math.MaxUint64},
	}
	for _, r := range ranges {
		want := referenceBetween(r[0], r[1])
		if got := Between(r[0], r[1]).ToSlice(); !slices.Equal(got, want) {
			t.Errorf("Between(%d, %d): got %d primes, want %d", r[0], r[1], len(got), len(want))
		}
	}
}
