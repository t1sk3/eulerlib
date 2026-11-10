package primes_test

import (
	"fmt"

	"github.com/t1sk3/eulerlib/v4/primes"
	"github.com/t1sk3/go-utils/stream"
)

func Example() {
	fmt.Println(primes.IsPrime(int64(1_000_000_007)))
	fmt.Println(primes.All[int]().Limit(10).ToSlice())
	fmt.Println(primes.NextPrime(100), primes.PrevPrime(100).Get())
	fmt.Println(primes.Between(90, 110).ToSlice())
	// Output:
	// true
	// [2 3 5 7 11 13 17 19 23 29]
	// 101 97
	// [97 101 103 107 109]
}

// Project Euler #7: the 10,001st prime.
func Example_euler7() {
	fmt.Println(primes.All[int]().Skip(10_000).FindFirst().Get())
	// Output: 104743
}

// Project Euler #10: the sum of all primes below two million.
func Example_euler10() {
	fmt.Println(stream.Sum(primes.Between(2, 2_000_000)))
	fmt.Println(primes.Count(2, 2_000_000))
	// Output:
	// 142913828922
	// 148933
}

func ExampleSieve() {
	sv := primes.NewSieve(50)
	fmt.Println(sv.Count(), sv.Has(47), sv.Has(49))
	fmt.Println(stream.Range(40, 60).Filter(sv.IsPrime).ToSlice())
	// Output:
	// 15 true false
	// [41 43 47 53 59]
}
