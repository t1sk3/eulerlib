package num_theory

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/t1sk3/eulerlib/v2/etc"
	"github.com/t1sk3/eulerlib/v2/prime_numbers"
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
	want := []int64{0, 5, 4, 13, 45, 90}
	for i, num := range testNums {
		got := DigitSum(num)
		if got != want[i] {
			t.Errorf("DigitSum(%d) == %d, want %d", num, got, want[i])
		}
	}
}

func TestDigitSumString(t *testing.T) {
	testNums := []string{"0", "5", "22", "562", "1234567890", "1234567890123456789"}
	want := []int64{0, 5, 4, 13, 45, 90}
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
		{}, // 1 has no prime factors
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
		if len(got) != len(want[i]) {
			t.Errorf("Factorize(%d) has %d factors, want %d", num, len(got), len(want[i]))
			continue
		}
		for k, v := range got {
			if v != want[i][k] {
				t.Errorf("Factorize(%d) == %v, want %v", num, got, want[i])
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
		got := etc.IsPalindrome(strconv.Itoa(num))
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
		got := etc.Max(num[0], num[1])
		if got != want[i] {
			t.Errorf("Max(%d) == %d, want %d", num, got, want[i])
		}
	}
}

func TestJoinSlice(t *testing.T) {
	testNums := [][]int{{1, 2, 3, 4, 5}, {2, 4, 6, 8, 10}, {12, 15, 18, 21, 24}}
	want := []string{"12345", "246810", "1215182124"}
	for i, num := range testNums {
		got := etc.JoinSlice(num)
		if got != want[i] {
			t.Errorf("JoinSlice(%d) == %s, want %s", num, got, want[i])
		}
	}
}

func TestMakeIntSlice(t *testing.T) {
	testNums := []int{12345, 246810, 1215182124}
	want := [][]int{{1, 2, 3, 4, 5}, {2, 4, 6, 8, 1, 0}, {1, 2, 1, 5, 1, 8, 2, 1, 2, 4}}
	for i, num := range testNums {
		got := etc.MakeIntSlice(num)
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
		got := etc.IsPandigital(num)
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
		got := Reduce(num, sumfunc)
		if got != want[i] {
			t.Errorf("Reduce(%d) == %d, want %d", num, got, want[i])
		}
	}
}

func TestTotient(t *testing.T) {
	testNums := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := []int64{1, 1, 2, 2, 4, 2, 6, 4, 6}
	for i, num := range testNums {
		got := etc.Totient(num)
		if got != want[i] {
			t.Errorf("Totient(%d) == %d, want %d", num, got, want[i])
		}
	}
}

func TestCombinations(t *testing.T) {
	testNums := [][]int64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9, 10}}
	want := [][][]int64{
		// bitmask order: 1..7 -> [1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]
		{{1}, {2}, {1, 2}, {3}, {1, 3}, {2, 3}, {1, 2, 3}},
		// bitmask order for 3 elements
		{{4}, {5}, {4, 5}, {6}, {4, 6}, {5, 6}, {4, 5, 6}},
		// bitmask order for 4 elements: 1..15
		{{7}, {8}, {7, 8}, {9}, {7, 9}, {8, 9}, {7, 8, 9}, {10}, {7, 10}, {8, 10}, {7, 8, 10}, {9, 10}, {7, 9, 10}, {8, 9, 10}, {7, 8, 9, 10}},
	}
	for i, num := range testNums {
		// pass 0 to indicate "any length" (function filters by ones count when n>0)
		got := Combinations(num, 0)
		if len(got) != len(want[i]) {
			t.Errorf("Combinations(%d) == %d, want %d", num, got, want[i])
			continue
		}
		for j := range got {
			if len(got[j]) != len(want[i][j]) {
				t.Errorf("Combinations(%d)[%d] == %d, want %d", num, j, got[j], want[i][j])
				continue
			}
			for k := range got[j] {
				if got[j][k] != want[i][j][k] {
					t.Errorf("Combinations(%d)[%d][%d] == %d, want %d", num, j, k, got[j][k], want[i][j][k])
				}
			}
		}
	}
}

func TestCountDivisors(t *testing.T) {
	cases := []struct{ n, want int64 }{
		{1, 1}, {2, 2}, {3, 2}, {4, 3}, {6, 4}, {12, 6}, {36, 9}, {100, 9}, {28, 6},
	}
	for _, tc := range cases {
		if got := CountDivisors(tc.n); got != tc.want {
			t.Errorf("CountDivisors(%d) = %d, want %d", tc.n, got, tc.want)
		}
	}
}

func TestSumDivisors(t *testing.T) {
	cases := []struct{ n, want int64 }{
		{1, 0}, {2, 1}, {6, 6}, {12, 16}, {28, 28}, {220, 284}, {496, 496},
	}
	for _, tc := range cases {
		if got := SumDivisors(tc.n); got != tc.want {
			t.Errorf("SumDivisors(%d) = %d, want %d", tc.n, got, tc.want)
		}
	}
}

func TestIsAbundantDeficientPerfect(t *testing.T) {
	for _, n := range []int64{6, 28, 496} {
		if !IsPerfect(n) {
			t.Errorf("IsPerfect(%d) = false, want true", n)
		}
		if IsAbundant(n) {
			t.Errorf("IsAbundant(%d) = true, want false", n)
		}
		if IsDeficient(n) {
			t.Errorf("IsDeficient(%d) = true, want false", n)
		}
	}
	if !IsAbundant(int64(12)) {
		t.Error("IsAbundant(12) = false, want true")
	}
	if !IsDeficient(int64(8)) {
		t.Error("IsDeficient(8) = false, want true")
	}
}

func TestIsAmicable(t *testing.T) {
	for _, n := range []int64{220, 284, 1184, 1210} {
		if !IsAmicable(n) {
			t.Errorf("IsAmicable(%d) = false, want true", n)
		}
	}
	for _, n := range []int64{1, 6, 28, 100} {
		if IsAmicable(n) {
			t.Errorf("IsAmicable(%d) = true, want false", n)
		}
	}
}

func TestPowBigInt(t *testing.T) {
	cases := []struct {
		b, n int64
		want int64
	}{
		{2, 0, 1}, {2, 1, 2}, {2, 10, 1024}, {3, 4, 81}, {5, 3, 125},
	}
	for _, tc := range cases {
		b := big.NewInt(tc.b)
		got := PowBigInt(b, tc.n)
		if got.Int64() != tc.want {
			t.Errorf("PowBigInt(%d, %d) = %d, want %d", tc.b, tc.n, got.Int64(), tc.want)
		}
		if b.Int64() != tc.b {
			t.Errorf("PowBigInt mutated b: got %d, want %d", b.Int64(), tc.b)
		}
	}
}

func TestFactorizeBigInt(t *testing.T) {
	n := big.NewInt(12)
	pairs := FactorizeBigInt(n)
	// 12 = 2^2 * 3^1
	if len(pairs) != 2 {
		t.Fatalf("FactorizeBigInt(12) returned %d pairs, want 2", len(pairs))
	}
	got := map[int64]int64{}
	for _, pair := range pairs {
		got[pair[0].Int64()] = pair[1].Int64()
	}
	if got[2] != 2 || got[3] != 1 {
		t.Errorf("FactorizeBigInt(12) = %v, want {2:2, 3:1}", got)
	}
	if n.Int64() != 12 {
		t.Errorf("FactorizeBigInt mutated n: got %d, want 12", n.Int64())
	}
}

func TestPrimeFactorsBigInt(t *testing.T) {
	n := big.NewInt(360) // 2^3 * 3^2 * 5
	orig := new(big.Int).Set(n)
	pairs := PrimeFactorsBigInt(n)
	if n.Cmp(orig) != 0 {
		t.Error("PrimeFactorsBigInt mutated input")
	}
	got := map[int64]int64{}
	for _, pair := range pairs {
		got[pair[0]] = pair[1]
	}
	if got[2] != 3 || got[3] != 2 || got[5] != 1 {
		t.Errorf("PrimeFactorsBigInt(360) = %v, want {2:3, 3:2, 5:1}", got)
	}
}

func TestReduceWithInit(t *testing.T) {
	add := func(a, b int64) int64 { return a + b }
	if got := ReduceWithInit(int64(10), []int64{1, 2, 3}, add); got != 16 {
		t.Errorf("ReduceWithInit(10, [1,2,3], +) = %d, want 16", got)
	}
	if got := ReduceWithInit(int64(0), []int64{}, add); got != 0 {
		t.Errorf("ReduceWithInit(0, [], +) = %d, want 0", got)
	}
}

func TestDigitSumStringLarge(t *testing.T) {
	// 25 nines → digit sum = 225
	got := DigitSumString("9999999999999999999999999")
	if got != 225 {
		t.Errorf("DigitSumString(25×'9') = %d, want 225", got)
	}
}

func TestDigitCount(t *testing.T) {
	cases := []struct{ n, want int64 }{
		{0, 1}, {5, 1}, {9, 1}, {10, 2}, {99, 2}, {100, 3}, {1234567890, 10},
	}
	for _, tc := range cases {
		if got := DigitCount(tc.n); got != tc.want {
			t.Errorf("DigitCount(%d) = %d, want %d", tc.n, got, tc.want)
		}
	}
}

func TestIsCoprime(t *testing.T) {
	cases := []struct {
		a, b int64
		want bool
	}{
		{1, 1, true}, {2, 3, true}, {8, 9, true}, {6, 9, false}, {17, 34, false}, {35, 64, true},
	}
	for _, tc := range cases {
		if got := IsCoprime(tc.a, tc.b); got != tc.want {
			t.Errorf("IsCoprime(%d, %d) = %t, want %t", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestExtGcd(t *testing.T) {
	cases := []struct{ a, b int64 }{
		{240, 46}, {46, 240}, {17, 5}, {0, 5}, {5, 0}, {-12, 18}, {12, -18},
	}
	for _, tc := range cases {
		g, x, y := ExtGcd(tc.a, tc.b)
		wantG := Gcd(etc.Abs(tc.a), etc.Abs(tc.b))
		if tc.a == 0 && tc.b == 0 {
			continue
		}
		if wantG == 0 {
			wantG = etc.Abs(tc.a) + etc.Abs(tc.b)
		}
		if g != wantG {
			t.Errorf("ExtGcd(%d, %d) gcd = %d, want %d", tc.a, tc.b, g, wantG)
		}
		if tc.a*x+tc.b*y != g {
			t.Errorf("ExtGcd(%d, %d) = (%d, %d, %d), but a*x+b*y = %d, want %d", tc.a, tc.b, g, x, y, tc.a*x+tc.b*y, g)
		}
	}
}

func TestModInverse(t *testing.T) {
	cases := []struct {
		a, m, want int64
		ok         bool
	}{
		{3, 11, 4, true},   // 3*4 = 12 ≡ 1 (mod 11)
		{10, 17, 12, true}, // 10*12 = 120 ≡ 1 (mod 17)
		{6, 9, 0, false},   // gcd(6,9) = 3, no inverse
		{1, 1, 0, true},    // any x works mod 1; canonical is 0
	}
	for _, tc := range cases {
		got, ok := ModInverse(tc.a, tc.m)
		if ok != tc.ok {
			t.Errorf("ModInverse(%d, %d) ok = %t, want %t", tc.a, tc.m, ok, tc.ok)
			continue
		}
		if ok && got != tc.want {
			t.Errorf("ModInverse(%d, %d) = %d, want %d", tc.a, tc.m, got, tc.want)
		}
		if ok && (tc.a*got)%tc.m != 1%tc.m {
			t.Errorf("ModInverse(%d, %d) = %d, but %d*%d mod %d = %d, want 1", tc.a, tc.m, got, tc.a, got, tc.m, (tc.a*got)%tc.m)
		}
	}
}

func TestCRT(t *testing.T) {
	// x ≡ 2 (mod 3), x ≡ 3 (mod 5), x ≡ 2 (mod 7) -> x = 23 (mod 105)
	x, m, ok := CRT([]int64{2, 3, 2}, []int64{3, 5, 7})
	if !ok {
		t.Fatal("CRT([2,3,2], [3,5,7]) ok = false, want true")
	}
	if m != 105 {
		t.Errorf("CRT([2,3,2], [3,5,7]) modulus = %d, want 105", m)
	}
	if x != 23 {
		t.Errorf("CRT([2,3,2], [3,5,7]) = %d, want 23", x)
	}

	// Non-coprime, consistent moduli: x ≡ 5 (mod 6), x ≡ 3 (mod 4) -> x = 11 (mod 12)
	x, m, ok = CRT([]int64{5, 3}, []int64{6, 4})
	if !ok || m != 12 || x != 11 {
		t.Errorf("CRT([5,3], [6,4]) = (%d, %d, %t), want (11, 12, true)", x, m, ok)
	}

	// Contradictory system: no solution.
	if _, _, ok := CRT([]int64{1, 2}, []int64{4, 4}); ok {
		t.Error("CRT([1,2], [4,4]) ok = true, want false (contradictory system)")
	}

	// Mismatched lengths.
	if _, _, ok := CRT([]int64{1, 2}, []int64{4}); ok {
		t.Error("CRT with mismatched lengths ok = true, want false")
	}
}

func TestPartition(t *testing.T) {
	// p(0)..p(10), from OEIS A000041.
	want := []int64{1, 1, 2, 3, 5, 7, 11, 15, 22, 30, 42}
	for n, w := range want {
		if got := Partition(int64(n)); got.Int64() != w {
			t.Errorf("Partition(%d) = %d, want %d", n, got, w)
		}
	}
	if got := Partition(int64(100)); got.String() != "190569292" {
		t.Errorf("Partition(100) = %s, want 190569292", got)
	}
	if got := Partition(int64(-1)); got.Sign() != 0 {
		t.Errorf("Partition(-1) = %s, want 0", got)
	}
}

func TestFactorizeSPF(t *testing.T) {
	spf := prime_numbers.ListSmallestPrimeFactors(int64(1000))
	for _, n := range []int64{1, 2, 360, 997, 1000} {
		got := FactorizeSPF(n, spf)
		want := Factorize(n)
		if len(got) != len(want) {
			t.Fatalf("FactorizeSPF(%d) = %v, want %v", n, got, want)
		}
		for p, exp := range want {
			if got[p] != exp {
				t.Errorf("FactorizeSPF(%d)[%d] = %d, want %d", n, p, got[p], exp)
			}
		}
	}
}
