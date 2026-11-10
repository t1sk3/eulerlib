package sequences_test

import (
	"fmt"

	"github.com/t1sk3/eulerlib/v4/sequences"
)

func Example() {
	fmt.Println(sequences.Collatz(6).ToSlice(), sequences.CollatzLength(27))
	// Output: [6 3 10 5 16 8 4 2 1] 112
}

// Project Euler #14: the starting number under one million with the
// longest Collatz chain.
func ExampleLongestCollatz() {
	fmt.Println(sequences.LongestCollatz(1_000_000))
	// Output: (837799, 525)
}
