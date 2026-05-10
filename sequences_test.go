package eulerlib

import "testing"

func TestCollatz(t *testing.T) {
	cases := []struct {
		n    int64
		want []int64
	}{
		{1, []int64{1}},
		{6, []int64{6, 3, 10, 5, 16, 8, 4, 2, 1}},
		{0, []int64{}},
		{-1, []int64{}},
	}
	for _, tc := range cases {
		got := Collatz(tc.n)
		if len(got) != len(tc.want) {
			t.Errorf("Collatz(%d) length = %d, want %d", tc.n, len(got), len(tc.want))
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("Collatz(%d)[%d] = %d, want %d", tc.n, i, got[i], tc.want[i])
			}
		}
	}
}

func TestCollatzLength(t *testing.T) {
	cases := []struct {
		n    int64
		want int
	}{
		{1, 1},
		{6, 9},
		{27, 112},
		{0, 0},
		{-5, 0},
	}
	for _, tc := range cases {
		if got := CollatzLength(tc.n); got != tc.want {
			t.Errorf("CollatzLength(%d) = %d, want %d", tc.n, got, tc.want)
		}
	}
}
