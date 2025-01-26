package eulerlib

import (
	"strconv"
	"testing"
)

func TestFactorial(t *testing.T) {
	testNums := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 20}
	want := []int64{1, 1, 2, 6, 24, 120, 720, 5040, 40320, 2432902008176640000}
	for i, num := range testNums {
		got := Factorial(num)
		if got != want[i] {
			t.Errorf("Factorial(%d) == %d, want %d", num, got, want[i])
		}
	}
}

func TestDigitSum(t *testing.T) {
	testNums := []int64{0, 5, 22, 562, 1234567890, 1234567890123456789}
	want := []int64{0, 5, 4, 13, 45, 90, 180}
	for i, num := range testNums {
		got := DigitSum(num)
		if got != want[i] {
			t.Errorf("DigitSum(%d) == %d, want %d", num, got, want[i])
		}
	}
}

func TestDigitSumString(t *testing.T) {
	testNums := []string{"0", "5", "22", "562", "1234567890", "1234567890123456789"}
	want := []int64{0, 5, 4, 13, 45, 90, 180}
	for i, num := range testNums {
		got := DigitSumString(num)
		if got != want[i] {
			t.Errorf("DigitSumString(%s) == %d, want %d", num, got, want[i])
		}
	}
}

func TestGcd(t *testing.T) {
	testNums := [][]int64{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {2, 4, 6, 8, 10, 12, 14, 16, 18, 20}, {12, 15, 18, 21, 24, 30}}
	want := []int64{1, 2, 3}
	for i, num := range testNums {
		got := Gcd(num...)
		if got != want[i] {
			t.Errorf("Gcd(%d) == %d, want %d", num, got, want[i])
		}
	}
}

func TestFactorize(t *testing.T) {
	testNums := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9553}
	want := []map[int64]int64{
		{1: 1},
		{2: 1},
		{3: 1},
		{2: 2},
		{5: 1},
		{2: 1, 3: 1},
		{7: 1},
		{2: 3},
		{3: 2},
		{2: 1, 5: 1},
		{41: 1, 233: 1},
	}
	for i, num := range testNums {
		got := Factorize(num)
		for k, v := range got {
			if v != want[i][k] {
				t.Errorf("Factorize(%d) == %d, want %d", num, got, want[i])
			}
		}
	}
}

func TestPowMod(t *testing.T) {
	testNums := [][]int64{{2, 3, 5}, {2, 10, 100}, {2, 100, 1000}, {2, 1000, 10000}, {2, 10000, 100000}, {2, 100000, 1000000}}
	want := []int64{3, 24, 376, 9376, 9376, 109376}
	for i, num := range testNums {
		got := PowMod(num[0], num[1], num[2])
		if got != want[i] {
			t.Errorf("PowMod(%d) == %d, want %d", num, got, want[i])
		}
	}
}

func TestBinomial(t *testing.T) {
	testNums := [][]int64{{5, 3}, {10, 5}, {20, 10}, {30, 15}, {40, 20}, {50, 25}}
	want := []int64{10, 252, 184756, 155117520, 137846528820, 126410606437752}
	for i, num := range testNums {
		got := Binomial(num[0], num[1])
		if got.Int64() != want[i] {
			t.Errorf("Binomial(%d) == %d, want %d", num, got, want[i])
		}
	}
}

func TestIsPalindrome(t *testing.T) {
	testNums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 11, 22, 33, 44, 55, 66, 77, 101, 111, 121, 131, 141, 151, 161, 171, 181, 191}
	want := []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true,
		true, true, true, true, true, true, true, true, true, true, true}
	for i, num := range testNums {
		got := IsPalindrome(strconv.Itoa(num))
		if got != want[i] {
			t.Errorf("IsPalindrome(%d) == %t, want %t", num, got, want[i])
		}
	}
}

func TestMaxInSlice(t *testing.T) {
	testNums := [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {2, 4, 6, 8, 10, 12, 14, 16, 18, 20}, {12, 15, 18, 21, 24, 30}}
	want := []int{10, 20, 30}
	for i, num := range testNums {
		got := MaxInSlice(num)
		if got != want[i] {
			t.Errorf("MaxInSlice(%d) == %d, want %d", num, got, want[i])
		}
	}
}

func TestMax(t *testing.T) {
	testNums := [][]int{{1, 2}, {2, 4}, {12, 15}, {2, 1}, {4, 2}, {15, 12}}
	want := []int{2, 4, 15, 2, 4, 15}
	for i, num := range testNums {
		got := Max(num[0], num[1])
		if got != want[i] {
			t.Errorf("Max(%d) == %d, want %d", num, got, want[i])
		}
	}
}

func TestJoinSlice(t *testing.T) {
	testNums := [][]int{{1, 2, 3, 4, 5}, {2, 4, 6, 8, 10}, {12, 15, 18, 21, 24}}
	want := []string{"12345", "246810", "1215182124"}
	for i, num := range testNums {
		got := JoinSlice(num)
		if got != want[i] {
			t.Errorf("JoinSlice(%d) == %s, want %s", num, got, want[i])
		}
	}
}

func TestMakeIntSlice(t *testing.T) {
	testNums := []int{12345, 246810, 1215182124}
	want := [][]int{{1, 2, 3, 4, 5}, {2, 4, 6, 8, 1, 0}, {1, 2, 1, 5, 1, 8, 2, 1, 2, 4}}
	for i, num := range testNums {
		got := MakeIntSlice(num)
		for j, e := range got {
			if e != want[i][j] {
				t.Errorf("MakeIntSlice(%d) == %d, want %d", num, got, want[i])
			}
		}
	}
}

func TestIsPandigital(t *testing.T) {
	testNums := []int64{12345, 246810, 1215182124, 123456789, 1234567890, 12345678901}
	want := []bool{false, false, false, true, true, true}
	for i, num := range testNums {
		got := IsPandigital(num)
		if got != want[i] {
			t.Errorf("IsPandigital(%d) == %t, want %t", num, got, want[i])
		}
	}
}

func TestIsSquare(t *testing.T) {
	testNums := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 16, 25, 36, 49, 64, 81}
	want := []bool{true, true, false, false, true, false, false, false, false, true, true, true, true, true, true, true}
	for i, num := range testNums {
		got := IsSquare(num)
		if got != want[i] {
			t.Errorf("IsSquare(%d) == %t, want %t", num, got, want[i])
		}
	}
}

func TestLcd(t *testing.T) {
	testNums := [][]int64{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {2, 4, 6, 8, 10, 12, 14, 16, 18, 20}, {12, 15, 18, 21, 24, 30}}
	want := []int64{2520, 5040, 2520}
	for i, num := range testNums {
		got := Lcd(num...)
		if got != want[i] {
			t.Errorf("Lcd(%d) == %d, want %d", num, got, want[i])
		}
	}
}

func TestLcm(t *testing.T) {
	testNums := [][]int64{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {2, 4, 6, 8, 10, 12, 14, 16, 18, 20}, {12, 15, 18, 21, 24, 30}}
	want := []int64{2520, 5040, 2520}
	for i, num := range testNums {
		got := Lcm(num...)
		if got != want[i] {
			t.Errorf("Lcm(%d) == %d, want %d", num, got, want[i])
		}
	}
}

func TestReduce(t *testing.T) {
	testNums := [][]int64{{2, 4}, {3, 6}, {4, 8}, {5, 10, 20}, {6, 12}, {7, 14}, {8, 16, 6, 100}, {9, 18}, {10, 20}}
	sumfunc := func(a, b int64) int64 { return a + b }
	want := []int64{6, 9, 12, 35, 18, 21, 130, 27, 30}
	for i, num := range testNums {
		got := Reduce(sumfunc, num)
		if got != want[i] {
			t.Errorf("Reduce(%d) == %d, want %d", num, got, want[i])
		}
	}
}

func TestTotient(t *testing.T) {
	testNums := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := []int64{1, 1, 2, 2, 4, 2, 6, 4, 6}
	for i, num := range testNums {
		got := Totient(num)
		if got != want[i] {
			t.Errorf("Totient(%d) == %d, want %d", num, got, want[i])
		}
	}
}
