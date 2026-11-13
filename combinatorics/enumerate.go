package combinatorics

import (
	"slices"

	"github.com/t1sk3/go-utils/stream"
)

// KPermutations returns every ordered selection of k distinct positions of
// items, n!/(n-k)! of them, in lexicographic order of their indexes:
// KPermutations([a b c], 2) yields [a b], [a c], [b a], [b c], [c a],
// [c b]. Positions are selected, so equal elements are not merged. It
// yields one empty slice for k == 0 and nothing for k < 0 or
// k > len(items). items is not modified.
func KPermutations[E any](items []E, k int) stream.Stream[[]E] {
	return stream.FromSeq(func(yield func([]E) bool) {
		a := slices.Clone(items)
		n := len(a)
		if k < 0 || k > n {
			return
		}
		idx := make([]int, 0, k)
		used := make([]bool, n)
		// rec extends idx to k indexes and reports whether to keep going.
		var rec func() bool
		rec = func() bool {
			if len(idx) == k {
				c := make([]E, k)
				for i, j := range idx {
					c[i] = a[j]
				}
				return yield(c)
			}
			for j := range n {
				if used[j] {
					continue
				}
				used[j] = true
				idx = append(idx, j)
				ok := rec()
				idx = idx[:len(idx)-1]
				used[j] = false
				if !ok {
					return false
				}
			}
			return true
		}
		rec()
	})
}

// CombinationsWithReplacement returns every k-element multiset of items:
// selections where a position may be picked more than once and order does
// not matter, C(n+k-1, k) of them, in lexicographic order of their
// indexes. CombinationsWithReplacement([a b], 2) yields [a a], [a b],
// [b b]. It yields one empty slice for k == 0, and nothing for k < 0 or
// when items is empty and k > 0. items is not modified.
func CombinationsWithReplacement[E any](items []E, k int) stream.Stream[[]E] {
	return stream.FromSeq(func(yield func([]E) bool) {
		a := slices.Clone(items)
		n := len(a)
		if k < 0 || (n == 0 && k > 0) {
			return
		}
		idx := make([]int, k) // non-decreasing
		for {
			c := make([]E, k)
			for i, j := range idx {
				c[i] = a[j]
			}
			if !yield(c) {
				return
			}
			// Advance the rightmost index that can still move, and reset
			// everything after it to the same value.
			i := k - 1
			for i >= 0 && idx[i] == n-1 {
				i--
			}
			if i < 0 {
				return
			}
			v := idx[i] + 1
			for j := i; j < k; j++ {
				idx[j] = v
			}
		}
	})
}

// Partitions returns every way to write n as a sum of positive integers,
// ignoring order, each as a slice of parts in descending order. They come
// in reverse lexicographic order: Partitions(4) yields [4], [3 1], [2 2],
// [2 1 1], [1 1 1 1]. There are Partition(n) of them. It yields one empty
// slice for n == 0 and nothing for n < 0.
func Partitions(n int) stream.Stream[[]int] {
	return stream.FromSeq(func(yield func([]int) bool) {
		if n < 0 {
			return
		}
		a := []int{}
		if n > 0 {
			a = append(a, n)
		}
		for {
			if !yield(slices.Clone(a)) {
				return
			}
			// Drop the trailing 1s, lower the last part above 1, and
			// refill what was removed with parts no larger than it.
			rest := 1
			for len(a) > 0 && a[len(a)-1] == 1 {
				a = a[:len(a)-1]
				rest++
			}
			if len(a) == 0 {
				return
			}
			a[len(a)-1]--
			v := a[len(a)-1]
			for rest > v {
				a = append(a, v)
				rest -= v
			}
			a = append(a, rest)
		}
	})
}
