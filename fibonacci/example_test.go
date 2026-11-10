package fibonacci_test

import (
	"fmt"

	"github.com/t1sk3/eulerlib/v4/fibonacci"
	"github.com/t1sk3/go-utils/stream"
)

func Example() {
	fmt.Println(fibonacci.First[int](10))
	fmt.Println(fibonacci.Nth[int64](90))
	fmt.Println(fibonacci.NthBig(100))
	// Output:
	// [0 1 1 2 3 5 8 13 21 34]
	// 2880067194370816120
	// 354224848179261915075
}

// Project Euler #2: the sum of the even Fibonacci numbers up to four million.
func ExampleAll() {
	fmt.Println(stream.Sum(fibonacci.All[int]().TakeWhile(stream.AtMost(4_000_000)).Filter(stream.IsEven[int])))
	// Output: 4613732
}

// Project Euler #25: the index of the first Fibonacci number with 1000 digits.
func ExampleIndexOfFirstWithDigits() {
	fmt.Println(fibonacci.IndexOfFirstWithDigits(1000))
	// Output: 4782
}
