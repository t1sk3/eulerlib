package pythagoras_test

import (
	"cmp"
	"fmt"

	"github.com/t1sk3/eulerlib/v4/pythagoras"
)

func Example() {
	fmt.Println(pythagoras.IsTriplet(3, 4, 5))
	byC := func(a, b pythagoras.Triplet[int]) int { return cmp.Compare(a.C, b.C) }
	fmt.Println(pythagoras.PrimitiveTriplets(30).Sorted(byC).ToSlice())
	// Output:
	// true
	// [{3 4 5} {5 12 13} {8 15 17} {7 24 25} {20 21 29}]
}

// Project Euler #9: the product abc of the triplet with a + b + c = 1000.
func ExampleTripletsWithPerimeterAtMost() {
	t := pythagoras.TripletsWithPerimeterAtMost(1000).
		Filter(func(t pythagoras.Triplet[int]) bool { return t.Perimeter() == 1000 }).
		FindFirst().Get()
	fmt.Println(t, t.A*t.B*t.C)
	// Output: {200 375 425} 31875000
}
