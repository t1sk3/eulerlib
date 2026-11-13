package digits

import (
	"math"
	"math/big"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/t1sk3/go-utils/stream"
)

func TestOf(t *testing.T) {
	cases := map[int][]int{0: {0}, 7: {7}, 10: {1, 0}, 1230: {1, 2, 3, 0}, -45: {4, 5}}
	for n, want := range cases {
		if got := Of(n); !slices.Equal(got, want) {
			t.Errorf("Of(%d) = %v, want %v", n, got, want)
		}
	}
	if got := Of(int64(math.MinInt64)); !slices.Equal(got, []int64{9, 2, 2, 3, 3, 7, 2, 0, 3, 6, 8, 5, 4, 7, 7, 5, 8, 0, 8}) {
		t.Errorf("Of(MinInt64) = %v", got)
	}
	if got := Of(uint64(math.MaxUint64)); len(got) != 20 {
		t.Errorf("Of(MaxUint64) has %d digits, want 20", len(got))
	}
}

func TestOfBase(t *testing.T) {
	if got := OfBase(10, 2); !slices.Equal(got, []int{1, 0, 1, 0}) {
		t.Errorf("OfBase(10, 2) = %v", got)
	}
	if got := OfBase(255, 16); !slices.Equal(got, []int{15, 15}) {
		t.Errorf("OfBase(255, 16) = %v", got)
	}
	mustPanic(t, "OfBase base 1", func() { OfBase(5, 1) })
	mustPanic(t, "OfBase base too big for E", func() { OfBase(int8(5), 200) })
}

func TestOfBigInt(t *testing.T) {
	cases := map[string][]int{"0": {0}, "7": {7}, "1230": {1, 2, 3, 0}, "-45": {4, 5}}
	for n, want := range cases {
		bi, _ := new(big.Int).SetString(n, 10)
		if got := OfBigInt(bi); !slices.Equal(got, want) {
			t.Errorf("OfBigInt(%s) = %v, want %v", n, got, want)
		}
	}
	big100fac, _ := new(big.Int).SetString(
		"93326215443944152681699238856266700490715968264381621468592963895217599993229915608941463976156518286253697920827223758251185210916864000000000000000000000000",
		10)
	if got := OfBigInt(big100fac); len(got) != 158 {
		t.Errorf("OfBigInt(100!) has %d digits, want 158", len(got))
	}
}

func TestSumBigInt(t *testing.T) {
	for n := -1000; n <= 100_000; n++ {
		bi := big.NewInt(int64(n))
		if got, want := SumBigInt(bi), SumString(strconv.Itoa(n)); got != want {
			t.Fatalf("SumBigInt(%d) = %d, want %d", n, got, want)
		}
	}
	big100fac, _ := new(big.Int).SetString(
		"93326215443944152681699238856266700490715968264381621468592963895217599993229915608941463976156518286253697920827223758251185210916864000000000000000000000000",
		10)
	if got, want := SumBigInt(big100fac), int64(648); got != want {
		t.Errorf("SumBigInt(100!) = %d, want %d", got, want)
	}
}

func TestFromDigitsRoundTrip(t *testing.T) {
	for n := 0; n <= 100_000; n++ {
		if got := FromDigits(Of(n)); got != n {
			t.Fatalf("FromDigits(Of(%d)) = %d", n, got)
		}
	}
	if FromDigits([]int{}) != 0 {
		t.Error("FromDigits(empty) != 0")
	}
}

func TestSumAndCount(t *testing.T) {
	for n := -1000; n <= 100_000; n++ {
		s := strconv.Itoa(n)
		if got, want := Sum(n), int(SumString(s)); got != want {
			t.Fatalf("Sum(%d) = %d, want %d", n, got, want)
		}
		if got, want := Count(n), len(strings.TrimPrefix(s, "-")); got != want {
			t.Fatalf("Count(%d) = %d, want %d", n, got, want)
		}
		if n >= 0 {
			if got, want := SumDigitString(s), SumString(s); got != want {
				t.Fatalf("SumDigitString(%s) = %d, want %d", s, got, want)
			}
		}
	}
	if got := Count(uint64(math.MaxUint64)); got != 20 {
		t.Errorf("Count(MaxUint64) = %d", got)
	}
	if got := SumString("2^15 = 32768"); got != 2+1+5+3+2+7+6+8 {
		t.Errorf("SumString = %d", got)
	}
}

func TestReverse(t *testing.T) {
	cases := map[int]int{0: 0, 1230: 321, 12345: 54321, -12: -21, 7: 7, 1000: 1}
	for n, want := range cases {
		if got := Reverse(n); got != want {
			t.Errorf("Reverse(%d) = %d, want %d", n, got, want)
		}
	}
}

func TestIsPalindrome(t *testing.T) {
	// Cross-check against the string definition.
	for n := 0; n <= 100_000; n++ {
		if got, want := IsPalindrome(n), IsPalindromeString(strconv.Itoa(n)); got != want {
			t.Fatalf("IsPalindrome(%d) = %v, want %v", n, got, want)
		}
		for _, b := range []int{2, 3, 16} {
			if got, want := IsPalindromeInBase(n, b), IsPalindromeString(strconv.FormatInt(int64(n), b)); got != want {
				t.Fatalf("IsPalindromeInBase(%d, %d) = %v, want %v", n, b, got, want)
			}
		}
	}
	if !IsPalindromeInBase(585, 2) || !IsPalindrome(585) {
		t.Error("585 is a palindrome in base 10 and base 2")
	}
	if IsPalindrome(uint64(math.MaxUint64)) {
		t.Error("IsPalindrome(MaxUint64) should be false")
	}
	if !IsPalindrome(uint64(18446744066044764481)) {
		t.Error("IsPalindrome(18446744066044764481)")
	}
	if !IsPalindromeString("") || !IsPalindromeString("été") || IsPalindromeString("ab") {
		t.Error("IsPalindromeString")
	}
}

func TestIsPandigital(t *testing.T) {
	cases := map[int64]bool{
		123456789: true, 918273645: true, 192384576: true,
		1234567890: false, // has a 0
		112345678:  false, // repeat
		12345678:   false, // only 8 digits
		1234567899: false, // 10 digits
		-123456789: true,  // |n|
		0:          false,
	}
	for n, want := range cases {
		if got := IsPandigital(n); got != want {
			t.Errorf("IsPandigital(%d) = %v, want %v", n, got, want)
		}
	}
	if !IsPandigitalN(2143, 4) || IsPandigitalN(2143, 5) || IsPandigitalN(1123, 4) || !IsPandigitalN(1, 1) {
		t.Error("IsPandigitalN")
	}
	mustPanic(t, "IsPandigitalN k=10", func() { IsPandigitalN(1, 10) })
}

func TestIsPandigitalInBase(t *testing.T) {
	if !IsPandigitalInBase(1023456789, 10) || !IsPandigitalInBase(10234567899, 10) || IsPandigitalInBase(123456789, 10) {
		t.Error("base 10")
	}
	if !IsPandigitalInBase(2, 2) || !IsPandigitalInBase(0b1100, 2) || IsPandigitalInBase(0b111, 2) {
		t.Error("base 2")
	}
	// 0x10FEDCBA98765432 contains every hex digit.
	if !IsPandigitalInBase(uint64(0x10FEDCBA98765432), 16) {
		t.Error("base 16")
	}
	if !IsPandigitalString("1023456789", 10) || IsPandigitalString("123456789", 10) || IsPandigitalString("10234567a9", 10) {
		t.Error("IsPandigitalString base 10")
	}
	if !IsPandigitalString("FEDCBA9876543210", 16) || IsPandigitalString("fedcba9876543210", 16) {
		t.Error("IsPandigitalString base 16 (upper-case alphabet)")
	}
	// Agrees with the numeric version.
	for n := int64(0); n < 200_000; n += 7 {
		for _, b := range []int{2, 3, 4, 5} {
			if IsPandigitalInBase(n, b) != IsPandigitalString(ToBase(n, b), b) {
				t.Fatalf("IsPandigitalInBase(%d, %d) disagrees with IsPandigitalString", n, b)
			}
		}
	}
}

func TestConcat(t *testing.T) {
	if got := Concat(1, 23, 4); got != 1234 {
		t.Errorf("Concat(1, 23, 4) = %d", got)
	}
	if got := Concat(192, 384, 576); got != 192384576 {
		t.Errorf("Concat(192, 384, 576) = %d", got)
	}
	if got := Concat(5, 0, 1); got != 501 {
		t.Errorf("Concat(5, 0, 1) = %d", got)
	}
	if got := Concat[int](); got != 0 {
		t.Errorf("Concat() = %d", got)
	}
	if got := ConcatString(1, 23, 0, 4); got != "12304" {
		t.Errorf("ConcatString = %q", got)
	}
	mustPanic(t, "Concat negative", func() { Concat(1, -2) })
}

func TestToBase(t *testing.T) {
	for _, n := range []int64{0, 1, -1, 35, 36, 1234567, -987654321, math.MaxInt64, math.MinInt64} {
		for b := 2; b <= 36; b++ {
			if got, want := ToBase(n, b), strings.ToUpper(strconv.FormatInt(n, b)); got != want {
				t.Fatalf("ToBase(%d, %d) = %q, want %q", n, b, got, want)
			}
		}
	}
	if got := ToBase(uint64(math.MaxUint64), 16); got != "FFFFFFFFFFFFFFFF" {
		t.Errorf("ToBase(MaxUint64, 16) = %q", got)
	}
	if got := ToBase(61, 62); got != "z" {
		t.Errorf("ToBase(61, 62) = %q", got)
	}
	if got := ToBase(62, 62); got != "10" {
		t.Errorf("ToBase(62, 62) = %q", got)
	}
	mustPanic(t, "ToBase 63", func() { ToBase(1, 63) })
	mustPanic(t, "ToBase 1", func() { ToBase(1, 1) })
}

func TestSumOf(t *testing.T) {
	if got := stream.RangeClosed(1, 10).Collect(SumOf[int]()); got != 46 {
		t.Errorf("SumOf over 1..10 = %d, want 46", got)
	}
	if got := stream.Empty[int]().Collect(SumOf[int]()); got != 0 {
		t.Errorf("SumOf over empty = %d", got)
	}
}

func mustPanic(t *testing.T, name string, f func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s: expected a panic", name)
		}
	}()
	f()
}
