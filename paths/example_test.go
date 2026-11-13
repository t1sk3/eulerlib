package paths_test

import (
	"fmt"

	"github.com/t1sk3/eulerlib/v4/paths"
)

func ExampleTriangleMax() {
	fmt.Println(paths.TriangleMax([][]int{
		{3},
		{7, 4},
		{2, 4, 6},
		{8, 5, 9, 3},
	}))
	// Output: 23
}

func ExampleGridMinTwoWays() {
	grid := [][]int{
		{131, 673, 234},
		{201, 96, 342},
		{630, 803, 746},
	}
	fmt.Println(paths.GridMinTwoWays(grid), paths.GridMinFourWays(grid))
	// Output: 1516 1516
}
