package digits_test

import (
	"fmt"

	"github.com/t1sk3/eulerlib/v4/digits"
	"github.com/t1sk3/go-utils/stream"
)

func Example() {
	fmt.Println(digits.Of(1230))
	fmt.Println(digits.Sum(1230), digits.Count(1230), digits.Reverse(1230))
	fmt.Println(digits.IsPalindromeInBase(585, 2), digits.ToBase(585, 2))
	fmt.Println(digits.Concat(192, 384, 576), digits.IsPandigital(192384576))
	// Output:
	// [1 2 3 0]
	// 6 4 321
	// true 1001001001
	// 192384576 true
}

// Project Euler #36: the sum of all numbers below one million that are
// palindromic in base 10 and base 2.
func Example_euler36() {
	sum := stream.Sum(stream.Range(1, 1_000_000).
		Filter(digits.IsPalindrome[int]).
		Filter(func(n int) bool { return digits.IsPalindromeInBase(n, 2) }))
	fmt.Println(sum)
	// Output: 872187
}
