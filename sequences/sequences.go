// Package sequences provides functions for generating numeric sequences,
// currently the Collatz conjecture sequence.
package sequences

import (
	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/constraints"
	"github.com/t1sk3/go-utils/types/pair"
)

// Collatz returns the full Collatz sequence starting from n until it reaches 1.
// Returns an empty slice for n <= 0.
func Collatz[E constraints.Signed](n E) []E {
	if n <= 0 {
		return []E{}
	}
	seq := []E{n}
	for n != 1 {
		if n%2 == 0 {
			n /= 2
		} else {
			n = 3*n + 1
		}
		seq = append(seq, n)
	}
	return seq
}

// CollatzLength returns the number of steps in the Collatz sequence starting from n.
// Returns 0 for n <= 0.
func CollatzLength[E constraints.Signed](n E) int {
	if n <= 0 {
		return 0
	}
	length := 1
	for n != 1 {
		if n%2 == 0 {
			n /= 2
		} else {
			n = 3*n + 1
		}
		length++
	}
	return length
}

// LongestCollatz returns the starting number in [1, limit) that produces the
// longest Collatz chain, paired with that chain's length (Project Euler
// #14). If several starting numbers are tied for the longest chain, the
// smallest is returned. Returns the zero Pair for limit <= 1.
func LongestCollatz[E constraints.Signed](limit E) pair.Pair[E, int] {
	withLength := stream.Map(stream.Range(E(1), limit), func(n E) pair.Pair[E, int] {
		return pair.Of(n, CollatzLength(n))
	})
	return withLength.Max(func(a, b pair.Pair[E, int]) int { return a.Second - b.Second }).Get()
}
