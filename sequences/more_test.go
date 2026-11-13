package sequences

import (
	"math/big"
	"testing"

	"github.com/t1sk3/eulerlib/v4/internal/testutil"
	"github.com/t1sk3/go-utils/types/pair"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func TestFarey(t *testing.T) {
	for n := -1; n <= 40; n++ {
		var want []pair.Pair[int, int]
		if n >= 1 {
			want = append(want, pair.Of(0, 1))
			// Insert every reduced fraction in order: a/b < c/d iff ad < bc.
			for b := 1; b <= n; b++ {
				for a := 1; a <= b; a++ {
					if gcd(a, b) == 1 {
						want = append(want, pair.Of(a, b))
					}
				}
			}
			for i := 1; i < len(want); i++ {
				for j := i; j > 0 && want[j].First*want[j-1].Second < want[j-1].First*want[j].Second; j-- {
					want[j], want[j-1] = want[j-1], want[j]
				}
			}
		}
		got := Farey(n).ToSlice()
		if len(got) != len(want) {
			t.Fatalf("Farey(%d) has %d terms, want %d", n, len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Fatalf("Farey(%d)[%d] = %v, want %v", n, i, got[i], want[i])
			}
		}
	}
	// Euler #72 (for a small bound): the length is 1 + Σ φ(k).
	if got := Farey(1000).Count(); got != 304193 {
		t.Errorf("len(Farey(1000)) = %d, want 304193", got)
	}
	// Narrow types: the intermediate values reach 2n.
	a, b := Farey(int8(127)).ToSlice(), Farey(127).ToSlice()
	if len(a) != len(b) {
		t.Fatalf("Farey[int8](127) has %d terms, want %d", len(a), len(b))
	}
	for i := range a {
		if int(a[i].First) != b[i].First || int(a[i].Second) != b[i].Second {
			t.Fatalf("Farey[int8](127)[%d] = %v, want %v", i, a[i], b[i])
		}
	}
	testutil.Reusable(t, Farey(12), -1)
	testutil.EarlyExit(t, Farey(12), 5)
}

func TestIsHappy(t *testing.T) {
	happy := []int{1, 7, 10, 13, 19, 23, 28, 31, 32, 44, 49, 68, 70, 79, 82, 86, 91, 94, 97, 100}
	i := 0
	for n := -3; n <= 100; n++ {
		want := i < len(happy) && happy[i] == n
		if want {
			i++
		}
		if IsHappy(n) != want {
			t.Errorf("IsHappy(%d) = %v", n, !want)
		}
	}
	// Euler #92 at a smaller scale: starting numbers below 10^5 that
	// arrive at 89, i.e. are not happy.
	unhappy := 0
	for n := 1; n < 100_000; n++ {
		if !IsHappy(n) {
			unhappy++
		}
	}
	if unhappy != 85623 {
		t.Errorf("unhappy numbers below 10^5: %d, want 85623", unhappy)
	}
	// The largest uint64 is handled without overflow (it is not happy).
	if IsHappy(uint64(1<<64 - 1)) {
		t.Error("IsHappy(2^64-1)")
	}
}

func TestChainLength(t *testing.T) {
	fact := [10]int{1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880}
	next := func(n int) int {
		s := 0
		for ; n > 0; n /= 10 {
			s += fact[n%10]
		}
		return s
	}
	for n, want := range map[int]int{145: 1, 169: 3, 871: 2, 69: 5, 78: 4, 540: 2} {
		if got := ChainLength(n, next); got != want {
			t.Errorf("ChainLength(%d) = %d, want %d", n, got, want)
		}
	}
	if got := ChainLength("a", func(s string) string {
		if len(s) == 4 {
			return "aa"
		}
		return s + "a"
	}); got != 4 {
		t.Errorf("ChainLength over strings = %d", got)
	}
}

func TestLinearRecurrenceMod(t *testing.T) {
	// Fibonacci.
	fib := []*big.Int{big.NewInt(0), big.NewInt(1)}
	for len(fib) < 400 {
		fib = append(fib, new(big.Int).Add(fib[len(fib)-1], fib[len(fib)-2]))
	}
	for _, m := range []int64{1, 2, 10, 1_000_000_007, 1<<63 - 1} {
		for n, f := range fib {
			want := new(big.Int).Mod(f, big.NewInt(m)).Int64()
			if got := LinearRecurrenceMod([]int64{1, 1}, []int64{0, 1}, uint64(n), m); got != want {
				t.Fatalf("F(%d) mod %d = %d, want %d", n, m, got, want)
			}
		}
	}
	// The last nine digits of F(10^18), via the Pisano period 1.5·10^9.
	const m = 1_000_000_000
	a := LinearRecurrenceMod([]int{1, 1}, []int{0, 1}, 1_000_000_000_000_000_000, m)
	b := LinearRecurrenceMod([]int{1, 1}, []int{0, 1}, 1_000_000_000_000_000_000%1_500_000_000, m)
	if a != b {
		t.Errorf("F(10^18) mod 10^9 = %d, want %d", a, b)
	}
	// An order-3 recurrence with negative coefficients and terms:
	// a(i) = 2a(i-1) - 3a(i-2) + 5a(i-3).
	coeffs, initial := []int{2, -3, 5}, []int{-1, 4, 2}
	seq := []*big.Int{big.NewInt(-1), big.NewInt(4), big.NewInt(2)}
	for len(seq) < 200 {
		l := len(seq)
		v := new(big.Int).Mul(big.NewInt(2), seq[l-1])
		v.Sub(v, new(big.Int).Mul(big.NewInt(3), seq[l-2]))
		v.Add(v, new(big.Int).Mul(big.NewInt(5), seq[l-3]))
		seq = append(seq, v)
	}
	for n, v := range seq {
		want := int(new(big.Int).Mod(v, big.NewInt(1009)).Int64())
		if got := LinearRecurrenceMod(coeffs, initial, uint64(n), 1009); got != want {
			t.Fatalf("a(%d) mod 1009 = %d, want %d", n, got, want)
		}
	}
	// Order 1 is a power: a(i) = 3a(i-1), a(0) = 1.
	if got := LinearRecurrenceMod([]uint8{3}, []uint8{1}, 100, 251); got != 20 { // 3^100 mod 251
		t.Errorf("3^100 mod 251 = %d", got)
	}
	// A modulus above 2^63.
	const big64 = uint64(18446744073709551557)
	want := new(big.Int).Mod(fib[399], new(big.Int).SetUint64(big64)).Uint64()
	if got := LinearRecurrenceMod([]uint64{1, 1}, []uint64{0, 1}, 399, big64); got != want {
		t.Errorf("F(399) mod %d = %d, want %d", big64, got, want)
	}
	for name, f := range map[string]func(){
		"m = 0":           func() { LinearRecurrenceMod([]int{1}, []int{1}, 5, 0) },
		"no coefficients": func() { LinearRecurrenceMod([]int{}, []int{}, 5, 7) },
		"length mismatch": func() { LinearRecurrenceMod([]int{1, 1}, []int{1}, 5, 7) },
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
