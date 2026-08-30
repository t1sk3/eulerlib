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
