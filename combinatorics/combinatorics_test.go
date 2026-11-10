package combinatorics

import (
	"fmt"
	"slices"
	"testing"

	"github.com/t1sk3/eulerlib/v4/internal/testutil"
	"github.com/t1sk3/go-utils/stream"
)

func TestFactorial(t *testing.T) {
	want := 1
	for n := range 21 {
		if n > 0 {
			want *= n
		}
		if got := Factorial(n); got != want {
			t.Fatalf("Factorial(%d) = %d, want %d", n, got, want)
		}
		if got := FactorialBig(int64(n)); got.Int64() != int64(want) {
			t.Fatalf("FactorialBig(%d) = %s", n, got)
		}
	}
	if got := FactorialBig(25).String(); got != "15511210043330985984000000" {
		t.Errorf("FactorialBig(25) = %s", got)
	}
	mustPanic(t, "Factorial(-1)", func() { Factorial(-1) })
}

func TestBinomial(t *testing.T) {
	// Pascal's triangle.
	row := []int64{1}
	for n := range 40 {
		for k := range row {
			if got := Binomial(n, k); got.Int64() != row[k] {
				t.Fatalf("Binomial(%d, %d) = %s, want %d", n, k, got, row[k])
			}
		}
		next := make([]int64, len(row)+1)
		next[0], next[len(row)] = 1, 1
		for k := 1; k < len(row); k++ {
			next[k] = row[k-1] + row[k]
		}
		row = next
	}
	if Binomial(5, -1).Sign() != 0 || Binomial(5, 6).Sign() != 0 {
		t.Error("Binomial out of range should be 0")
	}
	if got := Binomial(40, 20).String(); got != "137846528820" {
		t.Errorf("Binomial(40, 20) = %s (Euler #15)", got)
	}
}

func TestPartition(t *testing.T) {
	// p(0..20) from OEIS A000041.
	want := []int64{1, 1, 2, 3, 5, 7, 11, 15, 22, 30, 42, 56, 77, 101, 135, 176, 231, 297, 385, 490, 627}
	table := PartitionTable(20)
	for n, w := range want {
		if table[n].Int64() != w || Partition(n).Int64() != w {
			t.Fatalf("p(%d) = %s, want %d", n, table[n], w)
		}
	}
	if got := Partition(100).String(); got != "190569292" {
		t.Errorf("p(100) = %s", got)
	}
	if Partition(-1).Sign() != 0 {
		t.Error("p(-1) should be 0")
	}
	// Brute-force count of partitions for small n.
	var count func(n, max int) int
	count = func(n, max int) int {
		if n == 0 {
			return 1
		}
		c := 0
		for k := min(n, max); k >= 1; k-- {
			c += count(n-k, k)
		}
		return c
	}
	for n := range 40 {
		if int64(count(n, n)) != Partition(n).Int64() {
			t.Fatalf("p(%d) disagrees with brute force", n)
		}
	}
}

func TestPermutationCount(t *testing.T) {
	if PermutationCount([]rune("AAB")) != 3 || PermutationCount([]int{1, 2, 3, 4}) != 24 || PermutationCount([]int{}) != 1 {
		t.Error("PermutationCount")
	}
	if got := PermutationCount([]rune("MISSISSIPPI")); got != 34650 {
		t.Errorf("PermutationCount(MISSISSIPPI) = %d", got)
	}
	if got := PermutationCount([]int{1, 1, 2, 2, 3}); got != LexPermutations([]int{1, 1, 2, 2, 3}).Count() {
		t.Errorf("PermutationCount disagrees with LexPermutations: %d", got)
	}
}

func TestDigitFactorialSum(t *testing.T) {
	if DigitFactorialSum(145) != 145 || DigitFactorialSum(40585) != 40585 || DigitFactorialSum(0) != 1 || DigitFactorialSum(-145) != 145 {
		t.Error("DigitFactorialSum")
	}
}

func key[E any](s []E) string { return fmt.Sprint(s) }

func TestPermutations(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	orig := slices.Clone(items)
	all := Permutations(items).ToSlice()
	if len(all) != 120 {
		t.Fatalf("Permutations of 5 gave %d", len(all))
	}
	seen := map[string]bool{}
	for _, p := range all {
		seen[key(p)] = true
		s := slices.Clone(p)
		slices.Sort(s)
		if !slices.Equal(s, orig) {
			t.Fatalf("%v is not a permutation of %v", p, orig)
		}
	}
	if len(seen) != 120 {
		t.Errorf("Permutations gave %d distinct results, want 120", len(seen))
	}
	if !slices.Equal(items, orig) {
		t.Errorf("Permutations modified its input: %v", items)
	}
	if !slices.Equal(all[0], orig) {
		t.Errorf("first permutation = %v, want the input", all[0])
	}
	// The caller may modify yielded slices.
	for p := range Permutations(items).All() {
		p[0] = -1
	}
	if got := Permutations([]int{}).ToSlice(); len(got) != 1 || len(got[0]) != 0 {
		t.Errorf("Permutations(empty) = %v", got)
	}
	testutil.Reusable(t, Permutations(items), -1)
	testutil.EarlyExit(t, Permutations(items), 3)
}

func TestLexPermutations(t *testing.T) {
	if got := key(LexPermutations([]int{1, 1, 2}).ToSlice()); got != "[[1 1 2] [1 2 1] [2 1 1]]" {
		t.Errorf("LexPermutations([1 1 2]) = %s", got)
	}
	// Against sorting all permutations from Heap's algorithm.
	items := []int{3, 1, 4, 1, 5, 9}
	want := []string{}
	seen := map[string]bool{}
	heap := Permutations(items).ToSlice()
	slices.SortFunc(heap, slices.Compare)
	for _, p := range heap {
		if k := key(p); !seen[k] {
			seen[k] = true
			want = append(want, k)
		}
	}
	got := LexPermutations(items).Map(key[int]).ToSlice()
	if !slices.Equal(got, want) {
		t.Fatalf("LexPermutations disagrees with sorted Heap permutations: %d vs %d", len(got), len(want))
	}
	testutil.Reusable(t, LexPermutations(items), -1)
	testutil.EarlyExit(t, LexPermutations(items), 3)
}

func TestNthPermutation(t *testing.T) {
	for _, items := range [][]int{{0, 1, 2, 3, 4, 5}, {3, 1, 4, 1, 5, 9}, {2, 2, 2}, {}} {
		all := LexPermutations(items).ToSlice()
		for k, want := range all {
			if got := NthPermutation(items, k); !slices.Equal(got, want) {
				t.Fatalf("NthPermutation(%v, %d) = %v, want %v", items, k, got, want)
			}
			if got := LexPermutations(items).Skip(k).FindFirst().Get(); !slices.Equal(got, want) {
				t.Fatalf("Skip(%d) mismatch", k)
			}
		}
		mustPanic(t, "k too large", func() { NthPermutation(items, len(all)) })
	}
	digits := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	if got := key(NthPermutation(digits, 999_999)); got != "[2 7 8 3 9 1 5 4 6 0]" {
		t.Errorf("Euler #24 = %s", got)
	}
	mustPanic(t, "negative k", func() { NthPermutation(digits, -1) })
}

func TestCombinations(t *testing.T) {
	if got := key(Combinations([]string{"a", "b", "c", "d"}, 2).ToSlice()); got != "[[a b] [a c] [a d] [b c] [b d] [c d]]" {
		t.Errorf("Combinations(abcd, 2) = %s", got)
	}
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for k := -1; k <= len(items)+1; k++ {
		if got, want := int64(Combinations(items, k).Count()), Binomial(len(items), k).Int64(); got != want {
			t.Fatalf("Combinations(10, %d).Count() = %d, want %d", k, got, want)
		}
	}
	if got := Combinations(items, 0).ToSlice(); len(got) != 1 || len(got[0]) != 0 {
		t.Errorf("Combinations(k=0) = %v", got)
	}
	// 70 elements: the v3 bitmask version broke at 64.
	big := make([]int, 70)
	for i := range big {
		big[i] = i
	}
	if got := Combinations(big, 2).Count(); got != 70*69/2 {
		t.Errorf("Combinations(70, 2).Count() = %d", got)
	}
	last := Combinations(big, 2).FindLast().Get()
	if !slices.Equal(last, []int{68, 69}) {
		t.Errorf("last combination = %v", last)
	}
	testutil.Reusable(t, Combinations(items, 3), -1)
	testutil.EarlyExit(t, Combinations(items, 3), 3)
}

func TestSubsets(t *testing.T) {
	if got := key(Subsets([]int{1, 2, 3}).ToSlice()); got != "[[] [1] [2] [3] [1 2] [1 3] [2 3] [1 2 3]]" {
		t.Errorf("Subsets([1 2 3]) = %s", got)
	}
	if got := Subsets(make([]int, 12)).Count(); got != 4096 {
		t.Errorf("Subsets of 12 = %d", got)
	}
	// Works past 64 elements (only the first few are taken).
	if got := Subsets(make([]int, 100)).Limit(5).Count(); got != 5 {
		t.Errorf("Subsets of 100 Limit(5) = %d", got)
	}
	testutil.Reusable(t, Subsets([]int{1, 2, 3, 4}), -1)
	testutil.EarlyExit(t, Subsets([]int{1, 2, 3, 4}), 3)
}

func TestCartesianProduct(t *testing.T) {
	if got := key(CartesianProduct([]int{1, 2}, []int{3, 4, 5}).ToSlice()); got != "[[1 3] [1 4] [1 5] [2 3] [2 4] [2 5]]" {
		t.Errorf("CartesianProduct = %s", got)
	}
	if got := CartesianProduct[int]().ToSlice(); len(got) != 1 || len(got[0]) != 0 {
		t.Errorf("CartesianProduct() = %v", got)
	}
	if got := CartesianProduct([]int{1}, []int{}).Count(); got != 0 {
		t.Errorf("CartesianProduct with an empty set = %d", got)
	}
	d := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	if got := CartesianProduct(d, d, d).Count(); got != 1000 {
		t.Errorf("10^3 = %d", got)
	}
	testutil.Reusable(t, CartesianProduct(d, d), -1)
	testutil.EarlyExit(t, CartesianProduct(d, d), 3)
}

func TestInputNotModified(t *testing.T) {
	items := []int{5, 3, 1}
	for _, s := range []stream.Stream[[]int]{Permutations(items), LexPermutations(items), Combinations(items, 2), Subsets(items), CartesianProduct(items, items)} {
		for p := range s.All() {
			for i := range p {
				p[i] = 0
			}
		}
		if !slices.Equal(items, []int{5, 3, 1}) {
			t.Fatalf("input modified: %v", items)
		}
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

func BenchmarkPermutations10(b *testing.B) {
	items := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for b.Loop() {
		n := 0
		for range Permutations(items).All() {
			n++
		}
	}
}
