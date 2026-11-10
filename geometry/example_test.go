package geometry_test

import (
	"fmt"

	"github.com/t1sk3/eulerlib/v4/geometry"
	"github.com/t1sk3/go-utils/types/vector"
)

func Example() {
	fmt.Println(geometry.AreColinear(vector.Of2(0, 0), vector.Of2(1, 1), vector.Of2(3, 3)))
	fmt.Println(geometry.Quadrant(vector.Of2(-2, 5)))

	h := geometry.HexOf(1, -3)
	fmt.Println(h, h.S(), h.Distance(geometry.HexOf(-2, 2)))
	fmt.Println(geometry.HexRing(geometry.HexOf(0, 0), 1).Count())
	// Output:
	// true
	// 2
	// Hex(1, -3) 2 5
	// 6
}
