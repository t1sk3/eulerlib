package pythagoras

import (
	"testing"

	"github.com/t1sk3/eulerlib/v4/internal/testutil"
)

func gcd3(a, b, c int) int {
	return gcd(gcd(a, b), c)
}

func TestIsTriplet(t *testing.T) {
	if !IsTriplet(3, 4, 5) || !IsTriplet(5, 12, 13) || IsTriplet(1, 2, 3) {
		t.Error("IsTriplet")
	}
	tr := Triplet[int]{3, 4, 5}
	if tr.Perimeter() != 12 || tr.Scale(2) != (Triplet[int]{6, 8, 10}) {
		t.Error("Triplet methods")
	}
}

// bruteTriplets returns every triplet a < b < c <= maxC, keeping only
// primitive ones if primitive is set.
func bruteTriplets(maxC int, primitive bool) map[Triplet[int]]bool {
	out := map[Triplet[int]]bool{}
	for c := 1; c <= maxC; c++ {
		for b := 1; b < c; b++ {
			for a := 1; a < b; a++ {
				if a*a+b*b == c*c && (!primitive || gcd3(a, b, c) == 1) {
					out[Triplet[int]{a, b, c}] = true
				}
			}
		}
	}
	return out
}

func TestPrimitiveTriplets(t *testing.T) {
	want := bruteTriplets(100, true)
	got := PrimitiveTriplets(100).ToSlice()
	if len(want) != 16 {
		t.Fatalf("brute force found %d primitive triplets with c <= 100, the plan says 16", len(want))
	}
	if len(got) != len(want) {
		t.Fatalf("PrimitiveTriplets(100) has %d, want %d", len(got), len(want))
	}
	for _, tr := range got {
		if !want[tr] || !IsTriplet(tr.A, tr.B, tr.C) || tr.A >= tr.B {
			t.Errorf("unexpected triplet %v", tr)
		}
	}
	testutil.Reusable(t, PrimitiveTriplets(1000), -1)
	testutil.EarlyExit(t, PrimitiveTriplets(1000), 3)
	if PrimitiveTriplets(4).Count() != 0 || PrimitiveTriplets(5).Count() != 1 {
		t.Error("small bounds")
	}
}

func TestTripletsWithPerimeterAtMost(t *testing.T) {
	const p = 300
	want := map[Triplet[int]]bool{}
	for tr := range bruteTriplets(p, false) {
		if tr.Perimeter() <= p {
			want[tr] = true
		}
	}
	got := TripletsWithPerimeterAtMost(p).ToSlice()
	seen := map[Triplet[int]]bool{}
	for _, tr := range got {
		if seen[tr] {
			t.Fatalf("duplicate %v", tr)
		}
		seen[tr] = true
		if !want[tr] {
			t.Fatalf("unexpected %v", tr)
		}
	}
	if len(seen) != len(want) {
		t.Fatalf("got %d triplets, want %d", len(seen), len(want))
	}
	testutil.Reusable(t, TripletsWithPerimeterAtMost(1000), -1)
	testutil.EarlyExit(t, TripletsWithPerimeterAtMost(1000), 3)
}

func TestEuler39(t *testing.T) {
	// The perimeter p <= 1000 with the most right-triangle solutions.
	counts := TripletsWithPerimeterAtMost(1000).CountBy(Triplet[int].Perimeter)
	best := 0
	for per, c := range counts {
		if c > counts[best] || (c == counts[best] && per < best) {
			best = per
		}
	}
	// Brute force for comparison.
	bruteBest, bruteMax := 0, 0
	for per := 1; per <= 1000; per++ {
		n := 0
		for a := 1; a < per; a++ {
			for b := a + 1; a+b < per; b++ {
				if c := per - a - b; c > b && a*a+b*b == c*c {
					n++
				}
			}
		}
		if n > bruteMax {
			bruteBest, bruteMax = per, n
		}
	}
	if best != bruteBest || counts[best] != bruteMax {
		t.Errorf("Euler #39: %d (%d), brute force %d (%d)", best, counts[best], bruteBest, bruteMax)
	}
}
