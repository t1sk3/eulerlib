package numtheory_test

import (
	"fmt"
	"maps"
	"slices"

	"github.com/t1sk3/eulerlib/v4/numtheory"
	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/pair"
)

func Example() {
	fmt.Println(numtheory.Divisors(28), numtheory.SumDivisors(28))
	fmt.Println(numtheory.Factorize(360))
	for p, e := range numtheory.Factors(600851475143) { // Project Euler #3
		fmt.Printf("%d^%d;", p, e)
	}
	fmt.Println()
	fmt.Println(numtheory.PowMod(int64(2), 1_000_000, 1_000_000_007))
	fmt.Println(numtheory.ModInverse(int64(3), 11).Get())
	// Output:
	// [1 2 4 7 14 28] 28
	// map[2:3 3:2 5:1]
	// 71^1;839^1;1471^1;6857^1;
	// 235042059
	// 4
}

// Project Euler #1: the sum of all multiples of 3 or 5 below 1000.
func ExampleDivisibleByAny() {
	fmt.Println(stream.Sum(stream.Range(1, 1000).Filter(numtheory.DivisibleByAny(3, 5))))
	// Output: 233168
}

// Project Euler #5: the smallest number divisible by every number 1..20.
func ExampleLcmOf() {
	fmt.Println(stream.RangeClosed(1, 20).Collect(numtheory.LcmOf[int]()))
	// Output: 232792560
}

// Project Euler #21: the sum of all amicable numbers under 10000.
func ExampleAmicablePairs() {
	sum := numtheory.AmicablePairs(10_000).Fold(0, func(acc int, p pairOfInts) int {
		return acc + p.First + p.Second
	})
	fmt.Println(sum)
	// Output: 31626
}

func ExampleFactors() {
	m := maps.Collect(numtheory.Factors(360))
	keys := slices.Sorted(maps.Keys(m))
	fmt.Println(keys, m[2])
	// Output: [2 3 5] 3
}

type pairOfInts = pair.Pair[int, int]
