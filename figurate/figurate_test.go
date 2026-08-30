package figurate

import "testing"

func TestNthTriangular(t *testing.T) {
	cases := []struct{ n, want int64 }{
		{1, 1}, {2, 3}, {3, 6}, {4, 10}, {5, 15}, {10, 55},
	}
	for _, tc := range cases {
		if got := NthTriangular(tc.n); got != tc.want {
			t.Errorf("NthTriangular(%d) = %d, want %d", tc.n, got, tc.want)
		}
	}
}

func TestIsTriangular(t *testing.T) {
	triangulars := map[int64]bool{
		0: true, 1: true, 3: true, 6: true, 10: true, 15: true, 21: true, 28: true, 36: true, 45: true, 55: true,
		2: false, 4: false, 5: false, 7: false, 8: false, 9: false, 11: false,
	}
	for n, want := range triangulars {
		if got := IsTriangular(n); got != want {
			t.Errorf("IsTriangular(%d) = %t, want %t", n, got, want)
		}
	}
}

func TestNthPentagonal(t *testing.T) {
	cases := []struct{ n, want int64 }{
		{1, 1}, {2, 5}, {3, 12}, {4, 22}, {5, 35}, {10, 145},
	}
	for _, tc := range cases {
		if got := NthPentagonal(tc.n); got != tc.want {
			t.Errorf("NthPentagonal(%d) = %d, want %d", tc.n, got, tc.want)
		}
	}
}

func TestIsPentagonal(t *testing.T) {
	pentagonals := map[int64]bool{
		1: true, 5: true, 12: true, 22: true, 35: true, 51: true,
		2: false, 3: false, 4: false, 6: false, 10: false, 15: false,
	}
	for n, want := range pentagonals {
		if got := IsPentagonal(n); got != want {
			t.Errorf("IsPentagonal(%d) = %t, want %t", n, got, want)
		}
	}
}

func TestNthHexagonal(t *testing.T) {
	cases := []struct{ n, want int64 }{
		{1, 1}, {2, 6}, {3, 15}, {4, 28}, {5, 45}, {10, 190},
	}
	for _, tc := range cases {
		if got := NthHexagonal(tc.n); got != tc.want {
			t.Errorf("NthHexagonal(%d) = %d, want %d", tc.n, got, tc.want)
		}
	}
}

func TestIsHexagonal(t *testing.T) {
	hexagonals := map[int64]bool{
		1: true, 6: true, 15: true, 28: true, 45: true, 66: true,
		2: false, 3: false, 5: false, 7: false, 14: false,
	}
	for n, want := range hexagonals {
		if got := IsHexagonal(n); got != want {
			t.Errorf("IsHexagonal(%d) = %t, want %t", n, got, want)
		}
	}
}
