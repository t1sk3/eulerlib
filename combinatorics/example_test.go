package combinatorics_test

import (
	"fmt"

	"github.com/t1sk3/eulerlib/v4/combinatorics"
	"github.com/t1sk3/eulerlib/v4/digits"
)

func Example() {
	fmt.Println(combinatorics.Factorial(10), combinatorics.Binomial(40, 20))
	fmt.Println(combinatorics.LexPermutations([]int{1, 1, 2}).ToSlice())
	fmt.Println(combinatorics.Combinations([]string{"a", "b", "c"}, 2).ToSlice())
	fmt.Println(combinatorics.Partition(100))
	// Output:
	// 3628800 137846528820
	// [[1 1 2] [1 2 1] [2 1 1]]
	// [[a b] [a c] [b c]]
	// 190569292
}

// Project Euler #24: the millionth lexicographic permutation of 0..9.
func ExampleNthPermutation() {
	p := combinatorics.NthPermutation([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 999_999)
	fmt.Println(digits.FromDigits(p))
	// Output: 2783915460
}
