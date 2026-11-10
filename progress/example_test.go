package progress_test

import (
	"fmt"
	"io"

	"github.com/t1sk3/eulerlib/v4/primes"
	"github.com/t1sk3/eulerlib/v4/progress"
	"github.com/t1sk3/go-utils/stream"
)

func ExampleStream() {
	total := primes.Count(2, 2_000_000) // known up front, so the bar shows a percentage
	sum := stream.Sum(progress.Stream(primes.Between(2, 2_000_000), total, progress.WithWriter(io.Discard)))
	fmt.Println(total, sum)
	// Output: 148933 142913828922
}

func ExampleN() {
	sum := 0
	for i := range progress.N(1000, progress.WithWriter(io.Discard)) {
		sum += i
	}
	fmt.Println(sum)
	// Output: 499500
}
