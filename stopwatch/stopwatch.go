// Package stopwatch provides a tiny timing helper for benchmarking
// brute-force solutions: how long did that loop actually take?
package stopwatch

import (
	"fmt"
	"io"
	"os"
	"time"
)

// config holds the settings a Stopwatch call can be tuned with via Option.
type config struct {
	writer io.Writer
}

// Option configures a stopwatch.
type Option func(*config)

// WithWriter sets where the elapsed-time line is printed (default os.Stderr).
func WithWriter(w io.Writer) Option {
	return func(c *config) { c.writer = w }
}

// Start begins timing and returns a Stop function. Calling Stop prints
// "<label>: <elapsed>" to the configured writer and returns the elapsed
// duration. The common pattern is to defer it, timing everything up to the
// end of the enclosing function:
//
//	defer stopwatch.Start("part 2")()
//
// Or capture the returned duration directly:
//
//	stop := stopwatch.Start("brute force")
//	... // do the work
//	elapsed := stop()
func Start(label string, opts ...Option) func() time.Duration {
	cfg := config{writer: os.Stderr}
	for _, opt := range opts {
		opt(&cfg)
	}
	begin := time.Now()
	return func() time.Duration {
		elapsed := time.Since(begin)
		fmt.Fprintf(cfg.writer, "%s: %s\n", label, elapsed)
		return elapsed
	}
}

// Time runs fn, prints "<label>: <elapsed>" the same way Start's returned
// function does, and returns the elapsed duration.
func Time(label string, fn func(), opts ...Option) time.Duration {
	stop := Start(label, opts...)
	fn()
	return stop()
}
