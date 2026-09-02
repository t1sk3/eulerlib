package etc

import "testing"

func TestReverseString(t *testing.T) {
	input := "abcdefg"

	output := ReverseString(input)
	if output != "gfedcba" {
		t.Errorf("ReverseString(%s) returned %s, expected %s", input, output, "gfedcba")
	}
}

func TestListTotients(t *testing.T) {
	// φ: 0,1,1,2,2,4,2,6,4,6
	want := []int64{0, 1, 1, 2, 2, 4, 2, 6, 4, 6}
	got := ListTotients(int64(9))
	if len(got) != len(want) {
		t.Fatalf("ListTotients(9) length = %d, want %d", len(got), len(want))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("ListTotients(9)[%d] = %d, want %d", i, got[i], w)
		}
	}
}

func TestPhi(t *testing.T) {
	// φ: 0,1,1,2,2,4,2,6,4,6
	cases := []struct{ n, want int64 }{
		{0, 0}, {1, 1}, {2, 1}, {3, 2}, {4, 2},
		{5, 4}, {6, 2}, {7, 6}, {8, 4}, {9, 6},
		{-1, 0}, {-5, 0},
	}
	for _, tc := range cases {
		if got := Phi(tc.n); got != tc.want {
			t.Errorf("Phi(%d) = %d, want %d", tc.n, got, tc.want)
		}
	}
}

func TestPhiMatchesTotient(t *testing.T) {
	for n := int64(1); n <= 100; n++ {
		if got, want := Phi(n), Totient(n); got != want {
			t.Errorf("Phi(%d) = %d, Totient(%d) = %d; want equal", n, got, n, want)
		}
	}
}

func TestPhiMatchesListTotients(t *testing.T) {
	list := ListTotients(int64(50))
	for n := int64(0); n <= 50; n++ {
		if got, want := Phi(n), list[n]; got != want {
			t.Errorf("Phi(%d) = %d, ListTotients(50)[%d] = %d; want equal", n, got, n, want)
		}
	}
}

func TestAbs(t *testing.T) {
	cases := []struct{ n, want int64 }{
		{0, 0}, {5, 5}, {-5, 5}, {-100, 100},
	}
	for _, tc := range cases {
		if got := Abs(tc.n); got != tc.want {
			t.Errorf("Abs(%d) = %d, want %d", tc.n, got, tc.want)
		}
	}
}

func TestProduct(t *testing.T) {
	cases := []struct {
		nums []int64
		want int64
	}{
		{[]int64{}, 1},
		{[]int64{5}, 5},
		{[]int64{2, 3, 4}, 24},
		{[]int64{1, 2, 3, 4, 5}, 120},
	}
	for _, tc := range cases {
		if got := Product(tc.nums); got != tc.want {
			t.Errorf("Product(%v) = %d, want %d", tc.nums, got, tc.want)
		}
	}
}

func TestSumOfSquares(t *testing.T) {
	cases := []struct{ n, want int64 }{
		{1, 1}, {2, 5}, {3, 14}, {4, 30}, {10, 385},
	}
	for _, tc := range cases {
		if got := SumOfSquares(tc.n); got != tc.want {
			t.Errorf("SumOfSquares(%d) = %d, want %d", tc.n, got, tc.want)
		}
	}
}

func TestSquareOfSum(t *testing.T) {
	cases := []struct{ n, want int64 }{
		{1, 1}, {2, 9}, {3, 36}, {4, 100}, {10, 3025},
	}
	for _, tc := range cases {
		if got := SquareOfSum(tc.n); got != tc.want {
			t.Errorf("SquareOfSum(%d) = %d, want %d", tc.n, got, tc.want)
		}
	}
}

func TestListMobius(t *testing.T) {
	// μ(0..15), from OEIS A008683.
	want := []int64{0, 1, -1, -1, 0, -1, 1, -1, 0, 0, 1, -1, 0, -1, 1, 1}
	got := ListMobius(int64(15))
	if len(got) != len(want) {
		t.Fatalf("len(ListMobius(15)) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("ListMobius(15)[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestListMobiusZero(t *testing.T) {
	got := ListMobius(int64(0))
	if len(got) != 1 || got[0] != 0 {
		t.Errorf("ListMobius(0) = %v, want [0]", got)
	}
}

func TestRemoveDuplicatesFunc(t *testing.T) {
	eq := func(a, b int) bool { return a == b }
	got := RemoveDuplicatesFunc([]int{1, 2, 1, 3, 2, 4}, eq)
	want := []int{1, 2, 3, 4}
	if len(got) != len(want) {
		t.Fatalf("RemoveDuplicatesFunc length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("RemoveDuplicatesFunc[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}
