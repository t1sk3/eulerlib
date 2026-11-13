// Package combinatorics counts and enumerates arrangements: factorials,
// binomial coefficients, integer partitions, and lazy streams of
// permutations, combinations, subsets and Cartesian products.
//
// Every enumeration copies its input once per iteration, never modifies the
// caller's slice, and yields a fresh slice per element, so callers may keep
// (or modify) the slices they receive. That costs one allocation per
// element; collecting a stream of n! permutations still holds all of them
// in memory, so prefer ranging over it or bounding it with Limit.
package combinatorics

import (
	"cmp"
	"math"
	"math/big"
	"slices"

	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/constraints"
)

// Factorial returns n!. It overflows E silently for large n (21! does not
// fit in an int64; use FactorialBig). It panics if n < 0.
func Factorial[E constraints.Integer](n E) E {
	if n < 0 {
		panic("combinatorics: Factorial of a negative number")
	}
	r := E(1)
	for i := E(2); i <= n; i++ {
		r *= i
	}
	return r
}

// FactorialBig returns n! as a *big.Int. It panics if n < 0.
func FactorialBig(n int64) *big.Int {
	if n < 0 {
		panic("combinatorics: FactorialBig of a negative number")
	}
	return new(big.Int).MulRange(1, n)
}

// Binomial returns the binomial coefficient C(n, k), "n choose k": the
// number of k-element subsets of an n-element set. It is 0 when k < 0 or
// k > n. It panics if n < 0.
func Binomial[E constraints.Integer](n, k E) *big.Int {
	if n < 0 {
		panic("combinatorics: Binomial with negative n")
	}
	if k < 0 || k > n {
		return new(big.Int)
	}
	if k > n-k {
		k = n - k // C(n, k) = C(n, n-k), and n-k is the smaller one
	}
	if uint64(n) <= math.MaxInt64 {
		return new(big.Int).Binomial(int64(n), int64(k))
	}
	// n is a uint64 beyond big.Int.Binomial's int64 arguments:
	// C(n, k) = Π (n-k+i)/i for i = 1..k, exact at every step.
	r := big.NewInt(1)
	f := new(big.Int)
	for i := uint64(1); i <= uint64(k); i++ {
		r.Mul(r, f.SetUint64(uint64(n)-uint64(k)+i))
		r.Quo(r, f.SetUint64(i))
	}
	return r
}

// Partition returns p(n), the number of ways to write n as a sum of
// positive integers, ignoring order. p(0) is 1 and p(n) is 0 for n < 0.
func Partition[E constraints.Integer](n E) *big.Int {
	if n < 0 {
		return new(big.Int)
	}
	t := PartitionTable(int(n))
	return t[len(t)-1]
}

// PartitionTable returns p(0), p(1), ..., p(n), computed with Euler's
// pentagonal number recurrence in O(n√n) (Project Euler #76 and #78 need
// the whole table). It panics if n < 0.
func PartitionTable(n int) []*big.Int {
	if n < 0 {
		panic("combinatorics: PartitionTable of a negative number")
	}
	p := make([]*big.Int, n+1)
	p[0] = big.NewInt(1)
	for i := 1; i <= n; i++ {
		sum := new(big.Int)
		for k := 1; ; k++ {
			g1 := k * (3*k - 1) / 2
			if g1 > i {
				break
			}
			g2 := k * (3*k + 1) / 2
			add := k%2 == 1
			for _, g := range [2]int{g1, g2} {
				if g > i {
					continue
				}
				if add {
					sum.Add(sum, p[i-g])
				} else {
					sum.Sub(sum, p[i-g])
				}
			}
		}
		p[i] = sum
	}
	return p
}

// PermutationCount returns the number of distinct orderings of items,
// treating equal elements as indistinguishable: n! / (c1! c2! ...) where
// the ci count each distinct value. PermutationCount of "AAB" is 3. It
// panics if the result does not fit in an int.
func PermutationCount[E comparable](items []E) int {
	r := multinomial(stream.FromSlice(items).CountBy(stream.Identity[E]))
	if !r.IsInt64() || int64(int(r.Int64())) != r.Int64() {
		panic("combinatorics: PermutationCount overflows int")
	}
	return int(r.Int64())
}

// multinomial returns (Σc)! / Π(c!) for the counts in m.
func multinomial[K comparable](m map[K]int) *big.Int {
	r, total := big.NewInt(1), 0
	for _, c := range m {
		total += c
		r.Mul(r, new(big.Int).Binomial(int64(total), int64(c)))
	}
	return r
}

// DigitFactorialSum returns the sum of the factorials of the decimal
// digits of |n|: DigitFactorialSum(145) = 1! + 4! + 5! = 145 (Project
// Euler #34). In eulerlib v3 this was num_theory.FactorialDigitSum.
func DigitFactorialSum[E constraints.Integer](n E) E {
	fact := [10]uint64{1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880}
	if n == 0 {
		return 1
	}
	var s E
	for n != 0 {
		d := n % 10
		if d < 0 {
			d = -d
		}
		s += E(fact[d])
		n /= 10
	}
	return s
}

// Permutations returns every ordering of items, n! of them, using Heap's
// algorithm. Equal elements are not merged (see LexPermutations for that),
// and the order is Heap's, not lexicographic. The first permutation is
// items itself. items is not modified.
func Permutations[E any](items []E) stream.Stream[[]E] {
	return stream.FromSeq(func(yield func([]E) bool) {
		a := slices.Clone(items)
		if !yield(slices.Clone(a)) {
			return
		}
		c := make([]int, len(a))
		for i := 1; i < len(a); {
			if c[i] < i {
				if i%2 == 0 {
					a[0], a[i] = a[i], a[0]
				} else {
					a[c[i]], a[i] = a[i], a[c[i]]
				}
				if !yield(slices.Clone(a)) {
					return
				}
				c[i]++
				i = 1
			} else {
				c[i] = 0
				i++
			}
		}
	})
}

// LexPermutations returns the distinct orderings of items in lexicographic
// order, starting from the sorted order. Equal elements are not repeated:
// LexPermutations([1 1 2]) yields [1 1 2], [1 2 1], [2 1 1]. items is not
// modified.
func LexPermutations[E cmp.Ordered](items []E) stream.Stream[[]E] {
	return stream.FromSeq(func(yield func([]E) bool) {
		a := slices.Clone(items)
		slices.Sort(a)
		for {
			if !yield(slices.Clone(a)) || !nextPermutation(a) {
				return
			}
		}
	})
}

// nextPermutation rearranges a into the next lexicographic permutation and
// reports whether there was one.
func nextPermutation[E cmp.Ordered](a []E) bool {
	i := len(a) - 2
	for i >= 0 && a[i] >= a[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := len(a) - 1
	for a[j] <= a[i] {
		j--
	}
	a[i], a[j] = a[j], a[i]
	slices.Reverse(a[i+1:])
	return true
}

// NthPermutation returns the k-th (0-based) distinct permutation of items
// in lexicographic order, the same element LexPermutations(items) would
// yield at position k, without generating the ones before it (Project
// Euler #24). It handles repeated elements. It panics if k is negative or
// not less than the number of distinct permutations.
func NthPermutation[E cmp.Ordered](items []E, k int) []E {
	if k < 0 {
		panic("combinatorics: NthPermutation with negative k")
	}
	rest := slices.Clone(items)
	slices.Sort(rest)
	counts := map[E]int{}
	for _, v := range rest {
		counts[v]++
	}
	if big.NewInt(int64(k)).Cmp(multinomial(counts)) >= 0 {
		panic("combinatorics: NthPermutation k out of range")
	}
	kk := big.NewInt(int64(k))
	out := make([]E, 0, len(rest))
	for len(rest) > 0 {
		for i := 0; i < len(rest); i++ {
			if i > 0 && rest[i] == rest[i-1] {
				continue
			}
			// How many distinct permutations start with rest[i]?
			counts[rest[i]]--
			n := multinomial(counts)
			if kk.Cmp(n) < 0 {
				out = append(out, rest[i])
				rest = slices.Delete(rest, i, i+1)
				break
			}
			counts[rest[i]]++
			kk.Sub(kk, n)
		}
	}
	return out
}

// Combinations returns every k-element combination of items, as slices in
// the order the elements appear in items, in lexicographic order of their
// indexes. There are C(len(items), k) of them; positions are combined, so
// equal elements are not merged. It yields one empty slice for k == 0 and
// nothing for k < 0 or k > len(items). Unlike eulerlib v3's bitmask
// version, it works for any number of items.
func Combinations[E any](items []E, k int) stream.Stream[[]E] {
	return stream.FromSeq(func(yield func([]E) bool) {
		combinations(slices.Clone(items), k, yield)
	})
}

func combinations[E any](a []E, k int, yield func([]E) bool) bool {
	n := len(a)
	if k < 0 || k > n {
		return true
	}
	idx := make([]int, k)
	for i := range idx {
		idx[i] = i
	}
	for {
		c := make([]E, k)
		for i, j := range idx {
			c[i] = a[j]
		}
		if !yield(c) {
			return false
		}
		// Advance the rightmost index that can still move.
		i := k - 1
		for i >= 0 && idx[i] == n-k+i {
			i--
		}
		if i < 0 {
			return true
		}
		idx[i]++
		for j := i + 1; j < k; j++ {
			idx[j] = idx[j-1] + 1
		}
	}
}

// Subsets returns all 2^len(items) subsets of items, ordered by size and
// then like Combinations: the empty subset first and items itself last.
func Subsets[E any](items []E) stream.Stream[[]E] {
	return stream.FromSeq(func(yield func([]E) bool) {
		a := slices.Clone(items)
		for k := 0; k <= len(a); k++ {
			if !combinations(a, k, yield) {
				return
			}
		}
	})
}

// CartesianProduct returns every slice whose i-th element is taken from
// sets[i], with the last set varying fastest (like nested loops). It
// yields one empty slice when sets is empty, and nothing if any set is
// empty.
func CartesianProduct[E any](sets ...[]E) stream.Stream[[]E] {
	return stream.FromSeq(func(yield func([]E) bool) {
		ss := make([][]E, len(sets))
		for i, s := range sets {
			if len(s) == 0 {
				return
			}
			ss[i] = slices.Clone(s)
		}
		idx := make([]int, len(ss))
		for {
			c := make([]E, len(ss))
			for i, j := range idx {
				c[i] = ss[i][j]
			}
			if !yield(c) {
				return
			}
			i := len(ss) - 1
			for i >= 0 && idx[i] == len(ss[i])-1 {
				idx[i] = 0
				i--
			}
			if i < 0 {
				return
			}
			idx[i]++
		}
	})
}
