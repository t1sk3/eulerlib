package figurate

import (
	"math"
	"slices"
	"testing"

	"github.com/t1sk3/eulerlib/v4/internal/testutil"
)

func TestNth(t *testing.T) {
	for n := int64(0); n <= 1000; n++ {
		if NthTriangular(n) != n*(n+1)/2 || NthSquare(n) != n*n || NthPentagonal(n) != n*(3*n-1)/2 || NthHexagonal(n) != n*(2*n-1) {
			t.Fatalf("Nth*(%d)", n)
		}
		for s := int64(3); s <= 10; s++ {
			if got, want := NthPolygonal(s, n), ((s-2)*n*n-(s-4)*n)/2; got != want {
				t.Fatalf("NthPolygonal(%d, %d) = %d, want %d", s, n, got, want)
			}
		}
	}
	if NthPolygonal(3, int64(4)) != 10 || NthPolygonal(4, 5) != 25 || NthPolygonal(8, 3) != 21 {
		t.Error("NthPolygonal spot checks")
	}
	// No intermediate overflow: T(n) for the largest n whose T(n) fits.
	n := int64(4294967295) // T(n) = 9223372034707292160 < MaxInt64
	if NthTriangular(n) != 9223372034707292160 {
		t.Errorf("NthTriangular(%d) = %d", n, NthTriangular(n))
	}
	mustPanic(t, "NthPolygonal s=2", func() { NthPolygonal(2, 5) })
}

func TestIsAgainstBruteForce(t *testing.T) {
	const limit = 200_000
	tri, pent, hex := map[int]bool{0: true}, map[int]bool{0: true}, map[int]bool{0: true}
	for k := 1; NthTriangular(k) <= limit; k++ {
		tri[NthTriangular(k)] = true
	}
	for k := 1; NthPentagonal(k) <= limit; k++ {
		pent[NthPentagonal(k)] = true
	}
	for k := 1; NthHexagonal(k) <= limit; k++ {
		hex[NthHexagonal(k)] = true
	}
	for n := -10; n <= limit; n++ {
		if IsTriangular(n) != tri[n] || IsPentagonal(n) != pent[n] || IsHexagonal(n) != hex[n] {
			t.Fatalf("n=%d: tri %v/%v pent %v/%v hex %v/%v", n, IsTriangular(n), tri[n], IsPentagonal(n), pent[n], IsHexagonal(n), hex[n])
		}
	}
}

func TestIsNearTypeLimits(t *testing.T) {
	// Values near 2^64, where 8n+1 / 24n+1 no longer fit in a uint64 and
	// the math/big path is used.
	cases := []struct {
		k   uint64
		nth func(uint64) uint64
		is  func(uint64) bool
	}{
		{6_000_000_000, NthTriangular[uint64], IsTriangular[uint64]},
		{3_500_000_000, NthPentagonal[uint64], IsPentagonal[uint64]},
		{3_000_000_000, NthHexagonal[uint64], IsHexagonal[uint64]},
	}
	for _, c := range cases {
		v := c.nth(c.k)
		if v < 1<<63 {
			t.Fatalf("test value %d is not near 2^64", v)
		}
		if !c.is(v) || c.is(v+1) || c.is(v-1) {
			t.Errorf("Is* wrong around %d", v)
		}
	}
	if IsTriangular(uint64(math.MaxUint64)) || IsPentagonal(uint64(math.MaxUint64)) || IsHexagonal(uint64(math.MaxUint64)) {
		t.Error("MaxUint64 is not figurate")
	}
	// Near 2^63 for int64.
	v := NthTriangular(int64(4_294_967_295))
	if !IsTriangular(v) || IsTriangular(v+1) {
		t.Errorf("IsTriangular around %d", v)
	}
	if IsTriangular(int64(-1)) || IsPentagonal(int8(-1)) {
		t.Error("negative numbers are not figurate")
	}
}

func TestStreams(t *testing.T) {
	if got := Triangulars[int]().Limit(5).ToSlice(); !slices.Equal(got, []int{1, 3, 6, 10, 15}) {
		t.Errorf("Triangulars = %v", got)
	}
	if got := Squares[int]().Limit(5).ToSlice(); !slices.Equal(got, []int{1, 4, 9, 16, 25}) {
		t.Errorf("Squares = %v", got)
	}
	if got := Pentagonals[int]().Limit(5).ToSlice(); !slices.Equal(got, []int{1, 5, 12, 22, 35}) {
		t.Errorf("Pentagonals = %v", got)
	}
	if got := Hexagonals[int]().Limit(5).ToSlice(); !slices.Equal(got, []int{1, 6, 15, 28, 45}) {
		t.Errorf("Hexagonals = %v", got)
	}
	for s := 3; s <= 9; s++ {
		for i, v := range Polygonals(s).Limit(100).ToSlice() {
			if v != NthPolygonal(s, i+1) {
				t.Fatalf("Polygonals(%d)[%d] = %d", s, i, v)
			}
		}
	}
	// Streams stop before overflowing.
	got := Triangulars[uint8]().ToSlice()
	if got[len(got)-1] != 253 || len(got) != 22 {
		t.Errorf("Triangulars[uint8] = %v", got)
	}
	if got := Squares[int8]().FindLast().Get(); got != 121 {
		t.Errorf("last int8 square = %d", got)
	}
	if got := Squares[uint16]().FindLast().Get(); got != 255*255 {
		t.Errorf("last uint16 square = %d", got)
	}
	testutil.Reusable(t, Hexagonals[int](), 50)
	testutil.EarlyExit(t, Hexagonals[int](), 5)
	mustPanic(t, "Polygonals(2)", func() { Polygonals(2) })
}

func TestEuler45(t *testing.T) {
	// Brute force: walk T, P, H in lockstep and find the third common value.
	brute := func() int64 {
		ti, pi, hi := int64(1), int64(1), int64(1)
		found := 0
		for {
			tv, pv, hv := NthTriangular(ti), NthPentagonal(pi), NthHexagonal(hi)
			if tv == pv && pv == hv {
				found++
				if found == 3 {
					return tv
				}
				ti++
				continue
			}
			switch m := min(tv, pv, hv); m {
			case tv:
				ti++
			case pv:
				pi++
			default:
				hi++
			}
		}
	}
	want := brute()
	got := Hexagonals[int64]().Filter(IsPentagonal[int64]).Skip(2).FindFirst().Get()
	if got != want {
		t.Errorf("Euler #45: %d, brute force %d", got, want)
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
