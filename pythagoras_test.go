package eulerlib

import "testing"

func TestIsTriplet(t *testing.T) {
	testCases := []struct {
		a, b, c int64
		want    bool
	}{
		{3, 4, 5, true},
		{5, 12, 13, true},
		{8, 15, 17, true},
		{1, 2, 3, false},
		{6, 8, 11, false},
	}

	for _, tc := range testCases {
		got := IsTriplet(tc.a, tc.b, tc.c)
		if got != tc.want {
			t.Errorf("IsTriplet(%d, %d, %d) = %t, want %t", tc.a, tc.b, tc.c, got, tc.want)
		}
	}
}
