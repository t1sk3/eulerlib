package combinatorics

import (
	"math/big"
	"slices"
	"testing"

	"github.com/t1sk3/eulerlib/v4/internal/testutil"
)

func bigs(vs ...int64) []string {
	out := make([]string, len(vs))
	for i, v := range vs {
		out[i] = big.NewInt(v).String()
	}
	return out
}

func TestCountingNumbers(t *testing.T) {
	seq := func(f func(int) *big.Int, n int) []string {
		out := make([]string, n)
		for i := range out {
			out[i] = f(i).String()
		}
		return out
	}
	if got := seq(Catalan, 10); !slices.Equal(got, bigs(1, 1, 2, 5, 14, 42, 132, 429, 1430, 4862)) {
		t.Errorf("Catalan: %v", got)
	}
	if got := seq(Derangements, 10); !slices.Equal(got, bigs(1, 0, 1, 2, 9, 44, 265, 1854, 14833, 133496)) {
		t.Errorf("Derangements: %v", got)
	}
	if got := seq(Bell, 11); !slices.Equal(got, bigs(1, 1, 2, 5, 15, 52, 203, 877, 4140, 21147, 115975)) {
		t.Errorf("Bell: %v", got)
	}
	if got := Catalan(100).String(); got != "896519947090131496687170070074100632420837521538745909320" {
		t.Errorf("Catalan(100) = %s", got)
	}
	for n := 0; n <= 25; n++ {
		// Σ_k S(n, k) = B(n), Σ_k c(n, k) = n!, and derangements count the
		// permutations without fixed points.
		s1, s2 := new(big.Int), new(big.Int)
		for k := 0; k <= n+1; k++ {
			s1.Add(s1, Stirling1(n, k))
			s2.Add(s2, Stirling2(n, k))
		}
		if s2.Cmp(Bell(n)) != 0 || s1.Cmp(FactorialBig(int64(n))) != 0 {
			t.Fatalf("Stirling row sums for n = %d: %s, %s", n, s1, s2)
		}
	}
	if Stirling2(5, 2).Int64() != 15 || Stirling2(10, 3).Int64() != 9330 || Stirling1(5, 2).Int64() != 50 || Stirling1(4, 1).Int64() != 6 {
		t.Error("Stirling values")
	}
	if Stirling2(0, 0).Int64() != 1 || Stirling2(3, 0).Sign() != 0 || Stirling2(0, 3).Sign() != 0 || Stirling1(3, 5).Sign() != 0 {
		t.Error("Stirling edge cases")
	}
	for n := 0; n <= 8; n++ {
		items := make([]int, n)
		for i := range items {
			items[i] = i
		}
		count := Permutations(items).Filter(func(p []int) bool {
			for i, v := range p {
				if i == v {
					return false
				}
			}
			return true
		}).Count()
		if int64(count) != Derangements(n).Int64() {
			t.Errorf("Derangements(%d) = %s, counted %d", n, Derangements(n), count)
		}
	}
	mustPanic(t, "Catalan(-1)", func() { Catalan(-1) })
	mustPanic(t, "Derangements(-1)", func() { Derangements(-1) })
	mustPanic(t, "Bell(-1)", func() { Bell(-1) })
	mustPanic(t, "Stirling2(-1, 0)", func() { Stirling2(-1, 0) })
	mustPanic(t, "Stirling1(1, -1)", func() { Stirling1(1, -1) })
}

func TestWaysToSum(t *testing.T) {
	if got := WaysToSum(200, []int{1, 2, 5, 10, 20, 50, 100, 200}).Int64(); got != 73682 {
		t.Errorf("Euler #31 = %d", got)
	}
	// With every part allowed it is the partition function.
	table := PartitionTable(60)
	all := make([]int, 60)
	for i := range all {
		all[i] = i + 1
	}
	for n := 0; n <= 60; n++ {
		if got := WaysToSum(n, all); got.Cmp(table[n]) != 0 {
			t.Fatalf("WaysToSum(%d, 1..60) = %s, want %s", n, got, table[n])
		}
	}
	// Euler #76: 100 as a sum of at least two positive integers.
	parts := make([]int, 99)
	for i := range parts {
		parts[i] = i + 1
	}
	if got := WaysToSum(100, parts).Int64(); got != 190569291 {
		t.Errorf("Euler #76 = %d", got)
	}
	// Euler #77: first value writable as a sum of primes in over 5000 ways.
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}
	for n := 2; ; n++ {
		if WaysToSum(n, primes).Int64() > 5000 {
			if n != 71 {
				t.Errorf("Euler #77 = %d", n)
			}
			break
		}
	}
	in := []int{5, 1, 5, 2}
	if WaysToSum(6, in).Int64() != 5 || !slices.Equal(in, []int{5, 1, 5, 2}) { // 5+1, 2+2+2, 2+2+1+1, 2+1·4, 1·6
		t.Errorf("WaysToSum(6, [5 1 5 2]) = %s, input now %v", WaysToSum(6, in), in)
	}
	if WaysToSum(0, nil).Int64() != 1 || WaysToSum(5, nil).Sign() != 0 || WaysToSum(-1, []int{1}).Sign() != 0 {
		t.Error("WaysToSum edge cases")
	}
	mustPanic(t, "WaysToSum with part 0", func() { WaysToSum(5, []int{0, 1}) })
}

func TestBinomialMod(t *testing.T) {
	for _, p := range []int{2, 3, 5, 7, 13, 101} {
		bp := big.NewInt(int64(p))
		for n := 0; n <= 120; n++ {
			for k := -1; k <= n+1; k++ {
				want := new(big.Int).Mod(Binomial(n, k), bp).Int64()
				if got := BinomialMod(n, k, p); int64(got) != want {
					t.Fatalf("BinomialMod(%d, %d, %d) = %d, want %d", n, k, p, got, want)
				}
			}
		}
	}
	const p = 1_000_000_007
	for _, c := range [][2]uint64{{1_000_000, 500}, {1_000_000_006, 3}, {1 << 62, 7}, {1<<64 - 1, 5}} {
		var want *big.Int
		if c[0] <= 1<<63-1 {
			want = new(big.Int).Binomial(int64(c[0]), int64(c[1]))
		} else { // C(n, 5) = n(n-1)(n-2)(n-3)(n-4)/120
			want = big.NewInt(1)
			for i := uint64(0); i < c[1]; i++ {
				want.Mul(want, new(big.Int).SetUint64(c[0]-i))
			}
			want.Quo(want, FactorialBig(int64(c[1])))
		}
		want.Mod(want, big.NewInt(p))
		if got := BinomialMod(c[0], c[1], p); got != want.Uint64() {
			t.Errorf("BinomialMod(%d, %d, p) = %d, want %s", c[0], c[1], got, want)
		}
	}
	// Lucas with two base-p digits each: n = 2p + 999999986, k = p + 1, so
	// C(n, k) ≡ C(2, 1)·C(999999986, 1).
	if got := BinomialMod(3_000_000_000, 1_000_000_008, p); got != 999_999_965 {
		t.Errorf("BinomialMod(3e9, 1e9+8, p) = %d, want 999999965", got)
	}
	if BinomialMod(uint8(200), 100, 251) != uint8(new(big.Int).Mod(Binomial(200, 100), big.NewInt(251)).Int64()) {
		t.Error("BinomialMod[uint8]")
	}
	mustPanic(t, "BinomialMod negative n", func() { BinomialMod(-1, 0, 7) })
	mustPanic(t, "BinomialMod p = 1", func() { BinomialMod(5, 2, 1) })
}

func TestFactorialsMod(t *testing.T) {
	const p = 1_000_000_007
	fact, inv := FactorialsMod(2000, int64(p))
	if len(fact) != 2001 || len(inv) != 2001 || fact[0] != 1 || inv[0] != 1 {
		t.Fatal("FactorialsMod table shape")
	}
	for i := 0; i <= 2000; i++ {
		want := new(big.Int).Mod(FactorialBig(int64(i)), big.NewInt(p)).Int64()
		if fact[i] != want || fact[i]*inv[i]%p != 1 {
			t.Fatalf("fact[%d] = %d, inv = %d, want %d", i, fact[i], inv[i], want)
		}
	}
	for _, c := range [][2]int{{2000, 1000}, {1999, 3}, {50, 50}, {7, 0}} {
		got := fact[c[0]] * inv[c[1]] % p * inv[c[0]-c[1]] % p
		if got != BinomialMod(int64(c[0]), int64(c[1]), p) {
			t.Errorf("C(%d, %d) from the tables = %d", c[0], c[1], got)
		}
	}
	f8, i8 := FactorialsMod(12, uint8(13))
	for i := range f8 {
		if int(f8[i])*int(i8[i])%13 != 1 {
			t.Errorf("uint8 tables at %d: %d, %d", i, f8[i], i8[i])
		}
	}
	mustPanic(t, "FactorialsMod p <= n", func() { FactorialsMod(13, 13) })
	mustPanic(t, "FactorialsMod n < 0", func() { FactorialsMod(-1, 13) })
}

func TestKPermutations(t *testing.T) {
	got := KPermutations([]string{"a", "b", "c"}, 2).ToSlice()
	want := [][]string{{"a", "b"}, {"a", "c"}, {"b", "a"}, {"b", "c"}, {"c", "a"}, {"c", "b"}}
	if !slices.EqualFunc(got, want, slices.Equal[[]string]) {
		t.Errorf("KPermutations = %v", got)
	}
	items := []int{3, 1, 4, 1, 5}
	for k := -1; k <= 6; k++ {
		all := KPermutations(items, k).ToSlice()
		want := 0
		if k >= 0 && k <= 5 {
			want = 1
			for i := 0; i < k; i++ {
				want *= 5 - i
			}
		}
		if len(all) != want {
			t.Errorf("KPermutations(5 items, %d): %d, want %d", k, len(all), want)
		}
		for _, p := range all {
			if len(p) != k {
				t.Fatalf("KPermutations(%d) yielded %v", k, p)
			}
		}
	}
	// With k = n it is LexPermutations for distinct sorted items.
	a := KPermutations([]int{1, 2, 3, 4}, 4).ToSlice()
	b := LexPermutations([]int{1, 2, 3, 4}).ToSlice()
	if !slices.EqualFunc(a, b, slices.Equal[[]int]) {
		t.Error("KPermutations(n, n) is not in lexicographic order")
	}
	if !slices.Equal(items, []int{3, 1, 4, 1, 5}) {
		t.Error("KPermutations modified its input")
	}
	testutil.Reusable(t, KPermutations(items, 3), -1)
	testutil.EarlyExit(t, KPermutations(items, 3), 7)
}

func TestCombinationsWithReplacement(t *testing.T) {
	got := CombinationsWithReplacement([]string{"a", "b"}, 2).ToSlice()
	if !slices.EqualFunc(got, [][]string{{"a", "a"}, {"a", "b"}, {"b", "b"}}, slices.Equal[[]string]) {
		t.Errorf("CombinationsWithReplacement = %v", got)
	}
	for n := 0; n <= 5; n++ {
		items := make([]int, n)
		for i := range items {
			items[i] = i
		}
		for k := -1; k <= 5; k++ {
			all := CombinationsWithReplacement(items, k).ToSlice()
			want := int64(0)
			switch {
			case k == 0:
				want = 1
			case k > 0 && n > 0:
				want = Binomial(n+k-1, k).Int64()
			}
			if int64(len(all)) != want {
				t.Errorf("CombinationsWithReplacement(%d items, %d): %d, want %d", n, k, len(all), want)
			}
			for i, c := range all {
				if !slices.IsSorted(c) || (i > 0 && slices.Compare(all[i-1], c) >= 0) {
					t.Fatalf("CombinationsWithReplacement(%d, %d) out of order at %v", n, k, c)
				}
			}
		}
	}
	testutil.Reusable(t, CombinationsWithReplacement([]int{1, 2, 3}, 3), -1)
	testutil.EarlyExit(t, CombinationsWithReplacement([]int{1, 2, 3}, 3), 4)
}

func TestPartitions(t *testing.T) {
	got := Partitions(4).ToSlice()
	if !slices.EqualFunc(got, [][]int{{4}, {3, 1}, {2, 2}, {2, 1, 1}, {1, 1, 1, 1}}, slices.Equal[[]int]) {
		t.Errorf("Partitions(4) = %v", got)
	}
	table := PartitionTable(30)
	for n := 0; n <= 30; n++ {
		count := 0
		var prev []int
		for p := range Partitions(n).All() {
			sum := 0
			for i, v := range p {
				sum += v
				if v < 1 || (i > 0 && v > p[i-1]) {
					t.Fatalf("Partitions(%d) yielded %v", n, p)
				}
			}
			if sum != n || (prev != nil && slices.Compare(prev, p) <= 0) {
				t.Fatalf("Partitions(%d): %v after %v", n, p, prev)
			}
			prev = p
			count++
		}
		if int64(count) != table[n].Int64() {
			t.Errorf("Partitions(%d): %d, want %s", n, count, table[n])
		}
	}
	if Partitions(-1).Count() != 0 || len(Partitions(0).ToSlice()[0]) != 0 {
		t.Error("Partitions edge cases")
	}
	testutil.Reusable(t, Partitions(8), -1)
	testutil.EarlyExit(t, Partitions(8), 5)
}
