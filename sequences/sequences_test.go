package sequences

import (
	"slices"
	"testing"

	"github.com/t1sk3/eulerlib/v4/internal/testutil"
)

func TestCollatz(t *testing.T) {
	if got := Collatz(int64(6)).ToSlice(); !slices.Equal(got, []int64{6, 3, 10, 5, 16, 8, 4, 2, 1}) {
		t.Errorf("Collatz(6) = %v", got)
	}
	if got := Collatz(1).ToSlice(); !slices.Equal(got, []int{1}) {
		t.Errorf("Collatz(1) = %v", got)
	}
	if Collatz(0).Count() != 0 || Collatz(-3).Count() != 0 {
		t.Error("Collatz of n <= 0 should be empty")
	}
	testutil.Reusable(t, Collatz(27), -1)
	testutil.EarlyExit(t, Collatz(27), 5)
}

func TestCollatzLength(t *testing.T) {
	cases := map[int64]int{1: 1, 6: 9, 27: 112, 0: 0, -5: 0, 837799: 525}
	for n, want := range cases {
		if got := CollatzLength(n); got != want {
			t.Errorf("CollatzLength(%d) = %d, want %d", n, got, want)
		}
	}
	for n := int64(1); n <= 10_000; n++ {
		if got, want := Collatz(n).Count(), CollatzLength(n); got != want {
			t.Fatalf("Collatz(%d).Count() = %d, CollatzLength = %d", n, got, want)
		}
	}
}

func TestCollatzOverflow(t *testing.T) {
	// 113383 is the smallest start whose sequence leaves int32.
	if got := Collatz(int32(113382)).Count(); got != CollatzLength(int32(113382)) {
		t.Errorf("Collatz(int32(113382)) = %d terms", got)
	}
	defer func() {
		if recover() == nil {
			t.Error("Collatz(int32(113383)) should panic on overflow")
		}
	}()
	if CollatzLength(int32(113383)) != CollatzLength(int64(113383)) {
		t.Error("CollatzLength must not depend on E")
	}
	Collatz(int32(113383)).Count()
}

func TestCollatzLengths(t *testing.T) {
	l := CollatzLengths(100_000)
	if len(l) != 100_000 || l[0] != 0 {
		t.Fatal("CollatzLengths shape")
	}
	for i := 1; i < len(l); i++ {
		if l[i] != CollatzLength(i) {
			t.Fatalf("CollatzLengths[%d] = %d, want %d", i, l[i], CollatzLength(i))
		}
	}
	if len(CollatzLengths(0)) != 0 {
		t.Error("CollatzLengths(0)")
	}
}

func TestLongestCollatz(t *testing.T) {
	if got := LongestCollatz(int64(1)); got.First != 0 || got.Second != 0 {
		t.Errorf("LongestCollatz(1) = %v", got)
	}
	for _, limit := range []int{2, 10, 100, 1000, 10_000} {
		var bestN, bestL int
		for n := 1; n < limit; n++ {
			if l := CollatzLength(n); l > bestL {
				bestN, bestL = n, l
			}
		}
		if got := LongestCollatz(limit); got.First != bestN || got.Second != bestL {
			t.Errorf("LongestCollatz(%d) = %v, want (%d, %d)", limit, got, bestN, bestL)
		}
	}
}

func BenchmarkLongestCollatz1M(b *testing.B) {
	for b.Loop() {
		LongestCollatz(1_000_000)
	}
}
