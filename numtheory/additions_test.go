package numtheory

import (
	"math"
	"math/big"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/t1sk3/eulerlib/v4/primes"
)

func TestMultiplicativeOrder(t *testing.T) {
	for m := 1; m <= 300; m++ {
		for a := -3; a <= 2*m; a++ {
			want, ok := 0, false
			if Gcd(a, m) == 1 {
				x := ((a % m) + m) % m
				for k, p := 1, x%m; ; k, p = k+1, p*x%m {
					if p == 1%m {
						want, ok = k, true
						break
					}
				}
			}
			got, gotOK := MultiplicativeOrder(a, m).Value()
			if gotOK != ok || got != want {
				t.Fatalf("MultiplicativeOrder(%d, %d) = %d %v, want %d %v", a, m, got, gotOK, want, ok)
			}
		}
	}
	if MultiplicativeOrder(3, 0).IsPresent() || MultiplicativeOrder(3, -7).IsPresent() {
		t.Error("order modulo m <= 0 should be empty")
	}
	// 2^61-1 is prime and 2 has order 61 modulo it.
	if got := MultiplicativeOrder(uint64(2), 1<<61-1).Get(); got != 61 {
		t.Errorf("order of 2 mod 2^61-1 = %d", got)
	}
	if got := MultiplicativeOrder(uint8(3), 250).Get(); got != 100 {
		t.Errorf("order of 3 mod 250 [uint8] = %d, want 100", got)
	}
}

func TestRecurringCycleLength(t *testing.T) {
	brute := func(d int) int {
		seen := map[int]int{}
		for r, i := 1%d, 0; ; r, i = r*10%d, i+1 {
			if r == 0 {
				return 0
			}
			if j, ok := seen[r]; ok {
				return i - j
			}
			seen[r] = i
		}
	}
	best, bestD := 0, 0
	for d := 1; d < 1000; d++ {
		got := RecurringCycleLength(d)
		if got != brute(d) || RecurringCycleLength(-d) != got {
			t.Fatalf("RecurringCycleLength(%d) = %d, want %d", d, got, brute(d))
		}
		if got > best {
			best, bestD = got, d
		}
	}
	if bestD != 983 {
		t.Errorf("Euler #26: longest cycle below 1000 at d = %d, want 983", bestD)
	}
	if RecurringCycleLength(int8(-128)) != 0 || RecurringCycleLength(uint8(7)) != 6 {
		t.Error("RecurringCycleLength on narrow types")
	}
	defer func() {
		if recover() == nil {
			t.Error("RecurringCycleLength(0) should panic")
		}
	}()
	RecurringCycleLength(0)
}

func TestJacobiModSqrtPrimitiveRoot(t *testing.T) {
	for _, p := range primes.UpTo(400) {
		squares := map[int]bool{}
		for x := range p {
			squares[x*x%p] = true
		}
		for a := -p; a < 2*p; a++ {
			r := ((a % p) + p) % p
			root, ok := ModSqrt(a, p).Value()
			if ok != squares[r] {
				t.Fatalf("ModSqrt(%d, %d) present = %v, want %v", a, p, ok, squares[r])
			}
			if ok && (root*root%p != r || root > p-root && root != 0) {
				t.Fatalf("ModSqrt(%d, %d) = %d", a, p, root)
			}
			if p > 2 {
				want := -1
				if r == 0 {
					want = 0
				} else if squares[r] {
					want = 1
				}
				if got := Jacobi(a, p); got != want {
					t.Fatalf("Jacobi(%d, %d) = %d, want %d", a, p, got, want)
				}
			}
		}
		g := PrimitiveRoot(p)
		if ord := MultiplicativeOrder(g, p).Get(); p > 2 && ord != p-1 {
			t.Fatalf("PrimitiveRoot(%d) = %d has order %d", p, g, ord)
		}
		for h := 2; h < g; h++ {
			if MultiplicativeOrder(h, p).Get() == p-1 {
				t.Fatalf("PrimitiveRoot(%d) = %d, but %d is smaller", p, g, h)
			}
		}
	}
	// Jacobi is multiplicative in the modulus.
	for n := 1; n < 200; n += 2 {
		for a := range 60 {
			want := 1
			for _, p := range PrimeFactors(n) {
				want *= Jacobi(a, p)
			}
			if got := Jacobi(a, n); got != want {
				t.Fatalf("Jacobi(%d, %d) = %d, want %d", a, n, got, want)
			}
		}
	}
	// Large primes, including one with a big power of two in p-1.
	for _, p := range []uint64{18446744073709551557, 18446744069414584321, 1_000_000_007} {
		for range 200 {
			x := rand.Uint64N(p)
			a := MulMod(x, x, p)
			root := ModSqrt(a, p).Get()
			if MulMod(root, root, p) != a {
				t.Fatalf("ModSqrt(%d, %d) = %d", a, p, root)
			}
		}
	}
	for name, f := range map[string]func(){
		"ModSqrt composite":  func() { ModSqrt(2, 15) },
		"PrimitiveRoot comp": func() { PrimitiveRoot(15) },
		"Jacobi even":        func() { Jacobi(3, 8) },
		"Jacobi negative":    func() { Jacobi(3, -7) },
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s should panic", name)
				}
			}()
			f()
		}()
	}
}

func TestIrootPerfectPower(t *testing.T) {
	check := func(n uint64, k int) {
		r := Iroot(n, k)
		lo := new(big.Int).Exp(new(big.Int).SetUint64(r), big.NewInt(int64(k)), nil)
		hi := new(big.Int).SetUint64(r)
		hi.Exp(hi.Add(hi, big.NewInt(1)), big.NewInt(int64(k)), nil)
		bn := new(big.Int).SetUint64(n)
		if lo.Cmp(bn) > 0 || hi.Cmp(bn) <= 0 {
			t.Fatalf("Iroot(%d, %d) = %d", n, k, r)
		}
	}
	for k := 1; k <= 70; k++ {
		for _, n := range []uint64{0, 1, 2, 3, 7, 8, 9, 26, 27, 28, 1 << 32, math.MaxUint64, math.MaxUint64 - 1, 1 << 63} {
			check(n, k)
		}
		for range 300 {
			check(rand.Uint64()>>rand.IntN(64), k)
		}
		// Exact powers and their neighbours.
		for b := uint64(2); b < 50; b++ {
			if p := powOrZero(b, k); p != 0 {
				check(p, k)
				check(p-1, k)
				check(p+1, k)
			}
		}
	}
	if Iroot(int8(127), 2) != 11 || Iroot(uint8(255), 3) != 6 {
		t.Error("Iroot on narrow types")
	}
	for n := range 20000 {
		wantB, wantK := 0, 0
		for k := 14; k >= 2 && wantK == 0; k-- {
			for b := 2; ; b++ {
				p := Pow(b, k)
				if p == n {
					wantB, wantK = b, k
				}
				if p >= n {
					break
				}
			}
		}
		b, k, ok := PerfectPower(n)
		if b != wantB || k != wantK || ok != (wantK != 0) {
			t.Fatalf("PerfectPower(%d) = %d %d %v, want %d %d", n, b, k, ok, wantB, wantK)
		}
	}
	if b, k, ok := PerfectPower(uint64(1) << 63); !ok || b != 2 || k != 63 {
		t.Errorf("PerfectPower(2^63) = %d %d %v", b, k, ok)
	}
	if b, k, ok := PerfectPower(uint64(4294967295) * 4294967295); !ok || b != 4294967295 || k != 2 {
		t.Errorf("PerfectPower((2^32-1)^2) = %d %d %v", b, k, ok)
	}
	if _, _, ok := PerfectPower(-8); ok {
		t.Error("PerfectPower(-8) should be false")
	}
}

func TestSumOfTwoSquares(t *testing.T) {
	for n := -2; n < 5000; n++ {
		wantA, wantB, want := 0, 0, false
		for a := 0; a*a*2 <= n && !want; a++ {
			if b := Isqrt(n - a*a); b*b == n-a*a {
				wantA, wantB, want = a, b, true
			}
		}
		got, ok := SumOfTwoSquares(n).Value()
		if ok != want || got.First != wantA || got.Second != wantB {
			t.Fatalf("SumOfTwoSquares(%d) = %v %v, want (%d, %d) %v", n, got, ok, wantA, wantB, want)
		}
	}
	for n := 0; n <= 255; n++ {
		a, okA := SumOfTwoSquares(uint8(n)).Value()
		b, okB := SumOfTwoSquares(n).Value()
		if okA != okB || int(a.First) != b.First || int(a.Second) != b.Second {
			t.Fatalf("SumOfTwoSquares[uint8](%d) = %v %v", n, a, okA)
		}
	}
	// A prime ≡ 1 (mod 4) is a sum of two squares; one ≡ 3 (mod 4) is not.
	if p, ok := SumOfTwoSquares(uint64(1_000_000_009 * 4)).Value(); !ok || p.First*p.First+p.Second*p.Second != 4_000_000_036 {
		t.Errorf("SumOfTwoSquares(4·(1e9+9)) = %v %v", p, ok)
	}
	if SumOfTwoSquares(1_000_000_007).IsPresent() {
		t.Error("1e9+7 ≡ 3 (mod 4) is not a sum of two squares")
	}
}

func TestDivisorCountSieve(t *testing.T) {
	d := DivisorCountSieve(5000)
	if len(d) != 5001 || d[0] != 0 {
		t.Fatalf("DivisorCountSieve(5000): len %d, d[0] = %d", len(d), d[0])
	}
	for i := 1; i <= 5000; i++ {
		if d[i] != CountDivisors(i) {
			t.Fatalf("d[%d] = %d, want %d", i, d[i], CountDivisors(i))
		}
	}
	for i, v := range DivisorCountSieve(uint8(255)) {
		if int(v) != CountDivisors(i) {
			t.Fatalf("uint8 d[%d] = %d", i, v)
		}
	}
	if got := DivisorCountSieve(int8(127)); len(got) != 128 || got[120] != 16 {
		t.Errorf("DivisorCountSieve[int8](127)")
	}
	if got := DivisorCountSieve(0); !slices.Equal(got, []int{0}) {
		t.Errorf("DivisorCountSieve(0) = %v", got)
	}
}

func TestSqrtContinuedFraction(t *testing.T) {
	a0, period := SqrtContinuedFraction(23)
	if a0 != 4 || !slices.Equal(period, []int{1, 3, 1, 8}) {
		t.Errorf("√23 = [%d; %v]", a0, period)
	}
	if a0, period := SqrtContinuedFraction(49); a0 != 7 || len(period) != 0 {
		t.Errorf("√49 = [%d; %v]", a0, period)
	}
	odd := 0
	for n := 2; n <= 10000; n++ {
		a0, period := SqrtContinuedFraction(n)
		if a0 != Isqrt(n) {
			t.Fatalf("a0 of √%d = %d", n, a0)
		}
		if IsSquare(n) {
			continue
		}
		// The period ends in 2·a0 and is a palindrome before that.
		l := len(period)
		if period[l-1] != 2*a0 || !isPalindrome(period[:l-1]) {
			t.Fatalf("√%d = [%d; %v]", n, a0, period)
		}
		if l%2 == 1 {
			odd++
		}
	}
	if odd != 1322 {
		t.Errorf("Euler #64: %d odd periods, want 1322", odd)
	}
	if a0, period := SqrtContinuedFraction(uint8(255)); a0 != 15 || !slices.Equal(period, []uint8{1, 30}) {
		t.Errorf("√255 [uint8] = [%d; %v]", a0, period)
	}
	if a0, period := SqrtContinuedFraction(uint64(math.MaxUint64)); a0 != math.MaxUint32 || !slices.Equal(period, []uint64{1, 2 * math.MaxUint32}) {
		t.Errorf("√(2^64-1) = [%d; %v]", a0, period)
	}
}

func isPalindrome(a []int) bool {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		if a[i] != a[j] {
			return false
		}
	}
	return true
}

func TestConvergents(t *testing.T) {
	// √2: 1/1, 3/2, 7/5, 17/12, 41/29
	want := [][2]int64{{1, 1}, {3, 2}, {7, 5}, {17, 12}, {41, 29}}
	i := 0
	for h, k := range Convergents(SqrtCFTerms(2)) {
		if h.Int64() != want[i][0] || k.Int64() != want[i][1] {
			t.Fatalf("convergent %d of √2 = %s/%s", i, h, k)
		}
		h.SetInt64(0) // the caller owns the values
		if i++; i == len(want) {
			break
		}
	}
	// Euler #57: expansions of √2 whose numerator has more digits.
	count, i := 0, 0
	for h, k := range Convergents(SqrtCFTerms(2)) {
		if i > 0 && len(h.String()) > len(k.String()) {
			count++
		}
		if i++; i > 1000 {
			break
		}
	}
	if count != 153 {
		t.Errorf("Euler #57 = %d, want 153", count)
	}
	// Euler #65: digit sum of the numerator of the 100th convergent of e.
	i = 0
	for h := range Convergents(ECFTerms()) {
		if i++; i == 100 {
			sum := 0
			for _, c := range h.String() {
				sum += int(c - '0')
			}
			if sum != 272 {
				t.Errorf("Euler #65 = %d, want 272", sum)
			}
			break
		}
	}
	// A finite continued fraction: 415/93 = [4; 2, 6, 7].
	var last [2]string
	n := 0
	for h, k := range Convergents(slices.Values([]uint8{4, 2, 6, 7})) {
		last = [2]string{h.String(), k.String()}
		n++
	}
	if n != 4 || last != [2]string{"415", "93"} {
		t.Errorf("[4; 2, 6, 7] = %v after %d terms", last, n)
	}
	// A square's terms stop after a0.
	n = 0
	for range SqrtCFTerms(9) {
		n++
	}
	if n != 1 {
		t.Errorf("SqrtCFTerms(9) has %d terms", n)
	}
}

func TestPell(t *testing.T) {
	for d := -1; d <= 400; d++ {
		sol, ok := Pell(d).Value()
		if ok != (d > 0 && !IsSquare(d)) {
			t.Fatalf("Pell(%d) present = %v", d, ok)
		}
		if !ok {
			continue
		}
		x, y := sol.First, sol.Second
		lhs := new(big.Int).Mul(x, x)
		lhs.Sub(lhs, new(big.Int).Mul(new(big.Int).Mul(y, y), big.NewInt(int64(d))))
		if lhs.Cmp(big.NewInt(1)) != 0 || y.Sign() <= 0 {
			t.Fatalf("Pell(%d) = (%s, %s)", d, x, y)
		}
		// Minimal: no smaller y works (checked where that is cheap).
		if y.IsInt64() && y.Int64() < 3000 {
			for yy := int64(1); yy < y.Int64(); yy++ {
				if IsSquare(int64(d)*yy*yy + 1) {
					t.Fatalf("Pell(%d) = (%s, %s) is not minimal: y = %d works", d, x, y, yy)
				}
			}
		}
	}
	if got := Pell(61).Get(); got.First.String() != "1766319049" || got.Second.String() != "226153980" {
		t.Errorf("Pell(61) = %v", got)
	}
	// Euler #66: the d <= 1000 with the largest minimal x.
	best, bestD := new(big.Int), 0
	for d := 2; d <= 1000; d++ {
		if sol, ok := Pell(d).Value(); ok && sol.First.Cmp(best) > 0 {
			best, bestD = sol.First, d
		}
	}
	if bestD != 661 {
		t.Errorf("Euler #66 = %d, want 661", bestD)
	}
	if !Pell(uint8(255)).IsPresent() || Pell(uint8(225)).IsPresent() {
		t.Error("Pell on uint8")
	}
}
