package sequences

import (
	"github.com/t1sk3/eulerlib/v2/utils"
)

// Collatz returns the full Collatz sequence starting from n until it reaches 1.
// Returns an empty slice for n <= 0.
func Collatz[E utils.SignedInteger](n E) []E {
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
func CollatzLength[E utils.SignedInteger](n E) int {
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
