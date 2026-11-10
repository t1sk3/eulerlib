package fibonacci

import (
	"math/big"
	"slices"
	"testing"

	"github.com/t1sk3/eulerlib/v4/internal/testutil"
)

func TestAll(t *testing.T) {
	if got := All[int]().Limit(10).ToSlice(); !slices.Equal(got, []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}) {
		t.Errorf("All().Limit(10) = %v", got)
	}
	if got := All[int64]().Count(); got != 93 {
		t.Errorf("All[int64]().Count() = %d, want 93", got)
	}
	if got := All[int64]().FindLast().Get(); got != 7540113804746346429 {
		t.Errorf("last int64 Fibonacci = %d", got)
	}
	if got := All[uint64]().Count(); got != 94 {
		t.Errorf("All[uint64]().Count() = %d, want 94", got)
	}
	if got := All[uint8]().ToSlice(); got[len(got)-1] != 233 || len(got) != 14 {
		t.Errorf("All[uint8]() = %v", got)
	}
	if got := All[int8]().FindLast().Get(); got != 89 {
		t.Errorf("All[int8]() ends at %d, want 89", got)
	}
	// Agrees with AllBig on every term that fits.
	bigs := AllBig().Limit(93).ToSlice()
	for i, v := range All[int64]().ToSlice() {
		if bigs[i].Int64() != v {
			t.Fatalf("F(%d): %d vs %s", i, v, bigs[i])
		}
	}
	testutil.Reusable(t, All[int](), -1)
	testutil.EarlyExit(t, All[int](), 5)
}

func TestAllBig(t *testing.T) {
	s := AllBig().Limit(200).ToSlice()
	for i := 2; i < len(s); i++ {
		if new(big.Int).Add(s[i-1], s[i-2]).Cmp(s[i]) != 0 {
			t.Fatalf("F(%d) != F(%d) + F(%d)", i, i-1, i-2)
		}
		if s[i] == s[i-1] {
			t.Fatal("AllBig reuses *big.Int pointers")
		}
	}
	if s[100].String() != "354224848179261915075" {
		t.Errorf("F(100) = %s", s[100])
	}
	s[5].SetInt64(-1) // modifying a yielded value must not affect the stream
	if AllBig().Skip(6).FindFirst().Get().Int64() != 8 {
		t.Error("AllBig affected by caller modification")
	}
	testutil.Reusable(t, AllBig(), 50)
	testutil.EarlyExit(t, AllBig(), 5)
}

func TestFirstNth(t *testing.T) {
	if got := First[int](5); !slices.Equal(got, []int{0, 1, 1, 2, 3}) {
		t.Errorf("First(5) = %v", got)
	}
	if len(First[int](0)) != 0 || len(First[uint8](100)) != 14 {
		t.Error("First edge cases")
	}
	all := All[int64]().ToSlice()
	for n, want := range all {
		if Nth[int64](n) != want {
			t.Fatalf("Nth(%d) = %d, want %d", n, Nth[int64](n), want)
		}
	}
	mustPanic(t, "Nth overflow", func() { Nth[int64](93) })
	mustPanic(t, "Nth negative", func() { Nth[int](-1) })
}

func TestNthBig(t *testing.T) {
	s := AllBig().Limit(1000).ToSlice()
	for n, want := range s {
		if got := NthBig(int64(n)); got.Cmp(want) != 0 {
			t.Fatalf("NthBig(%d) = %s, want %s", n, got, want)
		}
	}
	if NthBig(100).String() != "354224848179261915075" {
		t.Error("NthBig(100)")
	}
	// F(2^20) has 219,140 digits; just check a known identity cheaply:
	// gcd(F(m), F(n)) = F(gcd(m, n)).
	g := new(big.Int).GCD(nil, nil, NthBig(6000), NthBig(4500))
	if g.Cmp(NthBig(1500)) != 0 {
		t.Error("gcd(F(6000), F(4500)) != F(1500)")
	}
	mustPanic(t, "NthBig negative", func() { NthBig(-1) })
}

func TestIndexOfFirstWithDigits(t *testing.T) {
	if got := IndexOfFirstWithDigits(3); got != 12 { // F(12) = 144
		t.Errorf("IndexOfFirstWithDigits(3) = %d, want 12", got)
	}
	if got := IndexOfFirstWithDigits(1); got != 1 {
		t.Errorf("IndexOfFirstWithDigits(1) = %d, want 1", got)
	}
	// Brute force against AllBig for small d.
	for d := 1; d <= 50; d++ {
		want := AllBig().Map(func(b *big.Int) int { return len(b.String()) }).
			Skip(1).
			TakeWhile(func(l int) bool { return l < d }).
			Count() + 1
		if got := IndexOfFirstWithDigits(d); got != want {
			t.Fatalf("IndexOfFirstWithDigits(%d) = %d, want %d", d, got, want)
		}
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
