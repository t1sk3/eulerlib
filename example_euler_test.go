// Project Euler solutions written with the v4 API. Every Example's answer is
// also checked against a brute-force reference in TestEulerAgainstBruteForce,
// so the expected outputs are computed, not copied from somewhere.
package eulerlib_test

import (
	"fmt"
	"maps"
	"slices"
	"testing"

	"github.com/t1sk3/eulerlib/v4/combinatorics"
	"github.com/t1sk3/eulerlib/v4/digits"
	"github.com/t1sk3/eulerlib/v4/figurate"
	"github.com/t1sk3/eulerlib/v4/numtheory"
	"github.com/t1sk3/eulerlib/v4/primes"
	"github.com/t1sk3/eulerlib/v4/sequences"
	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/pair"
)

// Solutions, shared by the examples and the brute-force test.

func euler1() int {
	return stream.Sum(stream.Range(1, 1000).Filter(numtheory.DivisibleByAny(3, 5)))
}

func euler5() int {
	return stream.RangeClosed(1, 20).Collect(numtheory.LcmOf[int]())
}

func euler7() int {
	return primes.All[int]().Skip(10_000).FindFirst().Get()
}

func euler10() int {
	return stream.Sum(primes.Between(2, 2_000_000))
}

func euler14() pair.Pair[int, int] {
	return sequences.LongestCollatz(1_000_000)
}

func euler21() int {
	return numtheory.AmicablePairs(10_000).Fold(0, func(acc int, p pair.Pair[int, int]) int {
		return acc + p.First + p.Second
	})
}

func euler24() int {
	return digits.FromDigits(combinatorics.NthPermutation([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 999_999))
}

func euler45() int64 {
	// Every hexagonal number is triangular; H(1) = 1 and H(143) = 40755 are
	// the first two that are also pentagonal.
	return figurate.Hexagonals[int64]().Filter(figurate.IsPentagonal[int64]).Skip(2).FindFirst().Get()
}

func Example_euler1() {
	fmt.Println(euler1())
	// Output: 233168
}

func Example_euler5() {
	fmt.Println(euler5())
	// Output: 232792560
}

func Example_euler7() {
	fmt.Println(euler7())
	// Output: 104743
}

func Example_euler10() {
	fmt.Println(euler10())
	// Output: 142913828922
}

func Example_euler14() {
	fmt.Println(euler14())
	// Output: (837799, 525)
}

func Example_euler21() {
	fmt.Println(euler21())
	// Output: 31626
}

func Example_euler24() {
	fmt.Println(euler24())
	// Output: 2783915460
}

func Example_euler45() {
	fmt.Println(euler45())
	// Output: 1533776805
}

// The "Working with streams" snippets from the README.
func Example_streams() {
	// Infinite prime stream, bounded with TakeWhile.
	fmt.Println(primes.All[int]().TakeWhile(stream.LessThan(30)).ToSlice())

	// Factors is an iter.Seq2: range over it, or collect it into a map.
	for p, e := range numtheory.Factors(360) {
		fmt.Print(p, "^", e, ";")
	}
	fmt.Println()
	fmt.Println(maps.Collect(numtheory.Factors(360)))

	// Collectors and predicate builders plug into stream pipelines.
	fmt.Println(stream.RangeClosed(1, 10).Collect(numtheory.LcmOf[int]()))
	fmt.Println(stream.Range(1, 20).Filter(numtheory.DivisibleBy(7)).ToSlice())

	// Lazy permutations: nothing is generated until it is consumed.
	for p := range combinatorics.LexPermutations([]int{1, 2, 3}).Limit(2).All() {
		fmt.Print(p)
	}
	fmt.Println()
	// Output:
	// [2 3 5 7 11 13 17 19 23 29]
	// 2^3;3^2;5^1;
	// map[2:3 3:2 5:1]
	// 2520
	// [7 14]
	// [1 2 3][1 3 2]
}

func TestEulerAgainstBruteForce(t *testing.T) {
	isPrime := func(n int) bool {
		if n < 2 {
			return false
		}
		for d := 2; d*d <= n; d++ {
			if n%d == 0 {
				return false
			}
		}
		return true
	}
	properSum := func(n int) int {
		s := 0
		for d := 1; d < n; d++ {
			if n%d == 0 {
				s += d
			}
		}
		return s
	}

	t.Run("1", func(t *testing.T) {
		want := 0
		for n := 1; n < 1000; n++ {
			if n%3 == 0 || n%5 == 0 {
				want += n
			}
		}
		if got := euler1(); got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
	t.Run("5", func(t *testing.T) {
		want := 1
		for !func() bool {
			for d := 1; d <= 20; d++ {
				if want%d != 0 {
					return false
				}
			}
			return true
		}() {
			want++
		}
		if got := euler5(); got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
	t.Run("7", func(t *testing.T) {
		n, count := 1, 0
		for count < 10_001 {
			n++
			if isPrime(n) {
				count++
			}
		}
		if got := euler7(); got != n {
			t.Errorf("got %d, want %d", got, n)
		}
	})
	t.Run("10", func(t *testing.T) {
		want := 0
		for n := 2; n < 2_000_000; n++ {
			if isPrime(n) {
				want += n
			}
		}
		if got := euler10(); got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
	t.Run("14", func(t *testing.T) {
		bestN, bestL := 0, 0
		for n := 1; n < 1_000_000; n++ {
			l := 1
			for m := n; m != 1; l++ {
				if m%2 == 0 {
					m /= 2
				} else {
					m = 3*m + 1
				}
			}
			if l > bestL {
				bestN, bestL = n, l
			}
		}
		if got := euler14(); got.First != bestN || got.Second != bestL {
			t.Errorf("got %v, want (%d, %d)", got, bestN, bestL)
		}
	})
	t.Run("21", func(t *testing.T) {
		want := 0
		for a := 2; a < 10_000; a++ {
			if b := properSum(a); b != a && properSum(b) == a {
				want += a
			}
		}
		if got := euler21(); got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
	t.Run("24", func(t *testing.T) {
		// Generate permutations of the sorted digits in lexicographic order
		// by recursion, stopping at the millionth.
		count, want := 0, ""
		var rec func(prefix, rest []int) bool
		rec = func(prefix, rest []int) bool {
			if len(rest) == 0 {
				count++
				if count == 1_000_000 {
					want = fmt.Sprint(prefix)
					return true
				}
				return false
			}
			for i := range rest {
				if rec(append(slices.Clone(prefix), rest[i]), slices.Concat(rest[:i:i], rest[i+1:])) {
					return true
				}
			}
			return false
		}
		rec(nil, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
		if got := fmt.Sprint(digits.Of(euler24())); got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
	t.Run("45", func(t *testing.T) {
		tri := map[int64]bool{}
		pent := map[int64]bool{}
		for k := int64(1); k < 100_000; k++ {
			tri[k*(k+1)/2] = true
			pent[k*(3*k-1)/2] = true
		}
		var found []int64
		for k := int64(1); len(found) < 3; k++ {
			if h := k * (2*k - 1); tri[h] && pent[h] {
				found = append(found, h)
			}
		}
		if got := euler45(); got != found[2] {
			t.Errorf("got %d, want %d", got, found[2])
		}
	})
}
