package pythagoras

import (
	"math"
	"testing"
)

// The bound checks used to overflow narrow element types and let through
// "triplets" with negative or wrapped members. Every type must yield
// exactly what a wide type yields for the same bound.
func TestNarrowTypesMatchWideOnes(t *testing.T) {
	wantPrim := map[Triplet[int64]]bool{}
	for tr := range PrimitiveTriplets[int64](math.MaxInt16).All() {
		wantPrim[tr] = true
	}
	n := 0
	for tr := range PrimitiveTriplets[int16](math.MaxInt16).All() {
		n++
		if !wantPrim[Triplet[int64]{int64(tr.A), int64(tr.B), int64(tr.C)}] {
			t.Fatalf("PrimitiveTriplets[int16] yielded %v", tr)
		}
	}
	if n != len(wantPrim) {
		t.Errorf("PrimitiveTriplets[int16]: %d triplets, want %d", n, len(wantPrim))
	}

	wantPer := map[Triplet[int64]]bool{}
	for tr := range TripletsWithPerimeterAtMost[int64](math.MaxInt16).All() {
		wantPer[tr] = true
	}
	n = 0
	for tr := range TripletsWithPerimeterAtMost[int16](math.MaxInt16).All() {
		n++
		if !wantPer[Triplet[int64]{int64(tr.A), int64(tr.B), int64(tr.C)}] {
			t.Fatalf("TripletsWithPerimeterAtMost[int16] yielded %v", tr)
		}
	}
	if n != len(wantPer) {
		t.Errorf("TripletsWithPerimeterAtMost[int16]: %d triplets, want %d", n, len(wantPer))
	}

	if got, want := PrimitiveTriplets[uint8](255).Count(), PrimitiveTriplets[int](255).Count(); got != want {
		t.Errorf("PrimitiveTriplets[uint8](255): %d, want %d", got, want)
	}
	if got, want := TripletsWithPerimeterAtMost[int8](127).Count(), TripletsWithPerimeterAtMost[int](127).Count(); got != want {
		t.Errorf("TripletsWithPerimeterAtMost[int8](127): %d, want %d", got, want)
	}
}

// The rewritten bounds must agree with the plain formulas for every bound.
func TestBoundsMatchBruteForce(t *testing.T) {
	for limit := 0; limit <= 400; limit++ {
		prim, per := 0, 0
		for m := 2; m*m <= limit; m++ {
			for n := 1; n < m; n++ {
				if (m-n)%2 == 0 || gcd(m, n) != 1 {
					continue
				}
				if m*m+n*n <= limit {
					prim++
				}
				if p := 2 * m * (m + n); p <= limit {
					per += limit / p
				}
			}
		}
		if got := PrimitiveTriplets(limit).Count(); got != prim {
			t.Fatalf("PrimitiveTriplets(%d): %d, want %d", limit, got, prim)
		}
		if got := TripletsWithPerimeterAtMost(limit).Count(); got != per {
			t.Fatalf("TripletsWithPerimeterAtMost(%d): %d, want %d", limit, got, per)
		}
	}
}
