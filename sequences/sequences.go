// Package sequences provides integer sequences that have no package of
// their own: the Collatz sequence (n -> n/2 for even n, n -> 3n+1 for odd
// n, until it reaches 1), the Farey sequence, happy numbers, the length of
// a chain before it repeats, and linear recurrences modulo m.
package sequences

import (
	"github.com/t1sk3/eulerlib/v4/internal/mathx"
	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/constraints"
	"github.com/t1sk3/go-utils/types/pair"
)

// Collatz returns the Collatz sequence starting at n: n, ..., 1. It is
// empty for n <= 0.
//
// 3n+1 can exceed E long before n is near E's maximum (for int32, starting
// at 113383 already climbs above 2^31), so the stream panics if the next
// value would overflow E. CollatzLength does not have this limit.
func Collatz[E constraints.Signed](n E) stream.Stream[E] {
	return stream.FromSeq(func(yield func(E) bool) {
		if n <= 0 {
			return
		}
		max := mathx.MaxOf[E]()
		for m := n; ; {
			if !yield(m) || m == 1 {
				return
			}
			if m%2 == 0 {
				m /= 2
			} else {
				if m > (max-1)/3 {
					panic("sequences: Collatz value overflows the element type")
				}
				m = 3*m + 1
			}
		}
	})
}

// CollatzLength returns the number of terms of the Collatz sequence
// starting at n, counting n and the final 1: CollatzLength(1) is 1 and
// CollatzLength(6) is 9. It returns 0 for n <= 0. It works in uint64
// internally, so it is not limited by the size of E, and equals
// Collatz(n).Count() whenever that does not overflow.
func CollatzLength[E constraints.Signed](n E) int {
	if n <= 0 {
		return 0
	}
	return collatzLength(uint64(n), nil)
}

// collatzLength counts the terms starting at m, using lengths (if any) as a
// cache for values below len(lengths).
func collatzLength(m uint64, lengths []int) int {
	steps := 0
	for m != 1 {
		if m < uint64(len(lengths)) && lengths[m] != 0 {
			return steps + lengths[m]
		}
		if m%2 == 0 {
			m /= 2
		} else {
			if m > (^uint64(0)-1)/3 {
				panic("sequences: Collatz value overflows uint64")
			}
			m = 3*m + 1
		}
		steps++
	}
	return steps + 1
}

// CollatzLengths returns a table l of length limit (or 0 for limit <= 0)
// with l[i] = CollatzLength(i) for every i in [0, limit). Each length
// reuses the ones already computed, which makes it much faster than
// calling CollatzLength for every i (Project Euler #14).
func CollatzLengths[E constraints.Signed](limit E) []int {
	if limit <= 0 {
		return []int{}
	}
	l := make([]int, uint64(limit))
	for i := 1; i < len(l); i++ {
		l[i] = collatzLength(uint64(i), l)
	}
	return l
}

// LongestCollatz returns the starting number in [1, limit) with the
// longest Collatz sequence, paired with that sequence's length (Project
// Euler #14). Ties go to the smallest starting number. It returns the
// zero Pair for limit <= 1.
func LongestCollatz[E constraints.Signed](limit E) pair.Pair[E, int] {
	var best pair.Pair[E, int]
	for i, l := range CollatzLengths(limit) {
		if l > best.Second {
			best = pair.Of(E(i), l)
		}
	}
	return best
}
