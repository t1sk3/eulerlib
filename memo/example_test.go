package memo_test

import (
	"fmt"

	"github.com/t1sk3/eulerlib/v4/memo"
)

// Project Euler #15: lattice paths through a 20×20 grid.
func ExampleMemoize2() {
	var paths func(r, c int) int64
	paths = memo.Memoize2(func(r, c int) int64 {
		if r == 0 || c == 0 {
			return 1
		}
		return paths(r-1, c) + paths(r, c-1)
	})
	fmt.Println(paths(20, 20))
	// Output: 137846528820
}

func ExampleMemoize() {
	var fib func(n int) int64
	fib = memo.Memoize(func(n int) int64 {
		if n < 2 {
			return int64(n)
		}
		return fib(n-1) + fib(n-2)
	})
	fmt.Println(fib(90))
	// Output: 2880067194370816120
}
