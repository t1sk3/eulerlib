package digits

import (
	"math"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestRotations(t *testing.T) {
	for _, c := range []struct {
		n    int
		want []int
	}{
		{197, []int{197, 971, 719}},
		{109, []int{109, 91, 910}},
		{7, []int{7}},
		{0, []int{0}},
		{-12, []int{12, 21}},
		{1000, []int{1000, 1, 10, 100}},
	} {
		if got := Rotations(c.n); !slices.Equal(got, c.want) {
			t.Errorf("Rotations(%d) = %v, want %v", c.n, got, c.want)
		}
	}
	for n := 1; n < 5000; n++ {
		s := strconv.Itoa(n)
		for i, r := range Rotations(n) {
			want, _ := strconv.Atoi(s[i:] + s[:i])
			if r != want {
				t.Fatalf("Rotations(%d)[%d] = %d, want %d", n, i, r, want)
			}
		}
	}
	if got := Rotations(uint64(math.MaxUint64)); len(got) != 20 || got[0] != math.MaxUint64 {
		t.Errorf("Rotations(MaxUint64) = %v", got)
	}
}

func TestSignatureIsPermutation(t *testing.T) {
	if Signature(1223) != [10]uint8{0, 1, 2, 1} || Signature(0) != [10]uint8{1} || Signature(-505) != [10]uint8{1, 0, 0, 0, 0, 2} {
		t.Error("Signature")
	}
	if got := Signature(int64(math.MinInt64)); got != [10]uint8{2, 0, 3, 3, 1, 2, 1, 3, 3, 1} { // 9223372036854775808
		t.Errorf("Signature(MinInt64) = %v", got)
	}
	if !IsPermutation(1487, 4817) || IsPermutation(1487, 4818) || IsPermutation(100, 10) || !IsPermutation(-21, 12) {
		t.Error("IsPermutation")
	}
	// Euler #52: the smallest x with x, 2x, ..., 6x made of the same digits.
	for x := 1; ; x++ {
		ok := true
		for k := 2; k <= 6 && ok; k++ {
			ok = IsPermutation(x, k*x)
		}
		if ok {
			if x != 142857 {
				t.Errorf("Euler #52 = %d", x)
			}
			break
		}
	}
}

func TestPowerSum(t *testing.T) {
	if PowerSum(1634, 4) != 1634 || PowerSum(-19, 2) != 82 || PowerSum(0, 3) != 0 || PowerSum(0, 0) != 1 || PowerSum(905, 0) != 3 {
		t.Error("PowerSum")
	}
	for n := range 3000 {
		if PowerSum(n, 1) != Sum(n) {
			t.Fatalf("PowerSum(%d, 1) != Sum", n)
		}
	}
	// Euler #30: numbers equal to the sum of the fifth powers of their digits.
	sum := 0
	for n := 10; n < 6*59049; n++ {
		if PowerSum(n, 5) == n {
			sum += n
		}
	}
	if sum != 443839 {
		t.Errorf("Euler #30 = %d", sum)
	}
	defer func() {
		if recover() == nil {
			t.Error("PowerSum with k < 0 should panic")
		}
	}()
	PowerSum(5, -1)
}

func TestLookAndSay(t *testing.T) {
	s := "1"
	for _, want := range []string{"11", "21", "1211", "111221", "312211", "13112221"} {
		if s = LookAndSay(s); s != want {
			t.Fatalf("LookAndSay = %q, want %q", s, want)
		}
	}
	if LookAndSay("") != "" || LookAndSay("aaaaaaaaaaaa") != "12a" {
		t.Error("LookAndSay edge cases")
	}
}

func TestChampernowneDigit(t *testing.T) {
	var b strings.Builder
	for i := 1; b.Len() < 20000; i++ {
		b.WriteString(strconv.Itoa(i))
	}
	s := b.String()
	for n := 1; n <= 20000; n++ {
		if got := ChampernowneDigit(n); got != int(s[n-1]-'0') {
			t.Fatalf("ChampernowneDigit(%d) = %d, want %c", n, got, s[n-1])
		}
	}
	prod := 1
	for n := 1; n <= 1_000_000; n *= 10 {
		prod *= ChampernowneDigit(n)
	}
	if prod != 210 {
		t.Errorf("Euler #40 = %d", prod)
	}
	if d := ChampernowneDigit(math.MaxInt); d < 0 || d > 9 {
		t.Errorf("ChampernowneDigit(MaxInt) = %d", d)
	}
	defer func() {
		if recover() == nil {
			t.Error("ChampernowneDigit(0) should panic")
		}
	}()
	ChampernowneDigit(0)
}
