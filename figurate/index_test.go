package figurate

import (
	"math"
	"math/big"
	"testing"
)

func TestPolygonalIndex(t *testing.T) {
	for s := 3; s <= 40; s++ {
		index := map[int]int{}
		for k := 0; ; k++ {
			p := NthPolygonal(s, k)
			if p > 30000 {
				break
			}
			index[p] = k
		}
		for n := -3; n <= 30000; n++ {
			want, ok := index[n]
			got, gotOK := PolygonalIndex(s, n).Value()
			if gotOK != ok || got != want {
				t.Fatalf("PolygonalIndex(%d, %d) = %d %v, want %d %v", s, n, got, gotOK, want, ok)
			}
			if IsPolygonal(s, n) != ok {
				t.Fatalf("IsPolygonal(%d, %d) = %v", s, n, !ok)
			}
		}
	}
	// It agrees with the dedicated tests.
	for n := range 20000 {
		if IsPolygonal(3, n) != IsTriangular(n) || IsPolygonal(5, n) != IsPentagonal(n) || IsPolygonal(6, n) != IsHexagonal(n) {
			t.Fatalf("IsPolygonal disagrees at %d", n)
		}
	}
	// Large values: every s-gonal number near the top of uint64, and its
	// neighbours.
	for _, s := range []uint64{3, 4, 5, 6, 8, 1000, 1 << 40, 1<<63 + 5} {
		// The largest k with P(s, k) <= MaxUint64, by bisection.
		lo, hi := uint64(1), uint64(1)<<33
		for lo < hi {
			mid := lo + (hi-lo+1)/2
			if _, ok := polygonal(s, mid); ok {
				lo = mid
			} else {
				hi = mid - 1
			}
		}
		for k := lo; k > lo-min(lo, 50); k-- {
			p, _ := polygonal(s, k)
			// Check polygonal itself against math/big.
			want := new(big.Int).SetUint64(s - 2)
			want.Mul(want, new(big.Int).SetUint64(k)).Mul(want, new(big.Int).SetUint64(k))
			sub := new(big.Int).Mul(new(big.Int).SetInt64(0).SetUint64(s), new(big.Int).SetUint64(k))
			sub.Sub(sub, new(big.Int).Mul(big.NewInt(4), new(big.Int).SetUint64(k)))
			want.Sub(want, sub).Rsh(want, 1)
			if !want.IsUint64() || want.Uint64() != p {
				t.Fatalf("polygonal(%d, %d) = %d, want %s", s, k, p, want)
			}
			if got, ok := PolygonalIndex(s, p).Value(); !ok || got != k {
				t.Fatalf("PolygonalIndex(%d, %d) = %d %v, want %d", s, p, got, ok, k)
			}
			if k > 2 && (IsPolygonal(s, p-1) || (p < math.MaxUint64 && IsPolygonal(s, p+1))) {
				t.Fatalf("neighbours of P(%d, %d) should not be %d-gonal", s, k, s)
			}
		}
	}
	if !IsPolygonal(uint8(3), 253) || IsPolygonal(uint8(3), 254) || PolygonalIndex(int8(4), 121).Get() != 11 {
		t.Error("PolygonalIndex on narrow types")
	}
	mustPanic(t, "PolygonalIndex(2, 5)", func() { PolygonalIndex(2, 5) })
}
