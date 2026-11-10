package stopwatch_test

import (
	"fmt"
	"io"
	"time"

	"github.com/t1sk3/eulerlib/v4/stopwatch"
)

func ExampleStart() {
	stop := stopwatch.Start("brute force", stopwatch.WithWriter(io.Discard))
	time.Sleep(time.Millisecond)
	fmt.Println(stop() >= time.Millisecond)
	// Output: true
}

func ExampleTime() {
	elapsed := stopwatch.Time("sum", func() {
		s := 0
		for i := range 1000 {
			s += i
		}
		_ = s
	}, stopwatch.WithWriter(io.Discard))
	fmt.Println(elapsed > 0)
	// Output: true
}
