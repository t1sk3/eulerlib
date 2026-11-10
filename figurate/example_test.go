package figurate_test

import (
	"fmt"

	"github.com/t1sk3/eulerlib/v4/figurate"
)

func Example() {
	fmt.Println(figurate.Triangulars[int]().Limit(5).ToSlice())
	fmt.Println(figurate.NthPentagonal(10), figurate.IsPentagonal(145), figurate.IsHexagonal(145))
	fmt.Println(figurate.NthPolygonal(8, 3))
	// Output:
	// [1 3 6 10 15]
	// 145 true false
	// 21
}

// Project Euler #45: the next triangular number after 40755 that is also
// pentagonal and hexagonal (every hexagonal number is triangular).
func Example_euler45() {
	fmt.Println(figurate.Hexagonals[int64]().Filter(figurate.IsPentagonal[int64]).Skip(2).FindFirst().Get())
	// Output: 1533776805
}
