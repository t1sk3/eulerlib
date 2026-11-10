// Package progress renders a live progress bar while ranging over a slice,
// a count, an iter.Seq or a go-utils Stream. Rendering is delegated to
// github.com/schollz/progressbar/v3.
package progress

import (
	"fmt"
	"io"
	"iter"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/t1sk3/go-utils/stream"
)

// Option configures the underlying progress bar. It's an alias for
// progressbar.Option, so any option from that package works here too —
// see https://pkg.go.dev/github.com/schollz/progressbar/v3#Option.
type Option = progressbar.Option

// WithWidth sets the number of characters used to render the bar itself.
func WithWidth(width int) Option {
	return progressbar.OptionSetWidth(width)
}

// WithWriter sets where the bar is rendered (default os.Stderr).
func WithWriter(w io.Writer) Option {
	return func(pb *progressbar.ProgressBar) {
		progressbar.OptionSetWriter(w)(pb)
		progressbar.OptionOnCompletion(func() { fmt.Fprintln(w) })(pb)
	}
}

// WithLabel sets a description printed before the bar.
func WithLabel(label string) Option {
	return progressbar.OptionSetDescription(label)
}

// newBar builds a progress bar with sane, good-looking defaults (mirroring
// progressbar.Default), letting opts override any of them.
func newBar(total int, opts ...Option) *progressbar.ProgressBar {
	defaults := []Option{
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(10),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionThrottle(65 * time.Millisecond),
		progressbar.OptionOnCompletion(func() {}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
	}
	return progressbar.NewOptions(total, append(defaults, opts...)...)
}

// Slice ranges over s exactly like range does, drawing a live progress
// bar (default: to stderr) as iteration proceeds:
//
//	for i, e := range progress.Slice(s) {
//	    ...
//	}
//
// The bar is created when iteration starts, so ranging over the result
// twice draws two bars. Breaking out of the loop early still finalizes the
// bar at whatever count was reached.
func Slice[T any](s []T, opts ...Option) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		bar := newBar(len(s), opts...)
		defer bar.Finish()
		for i, e := range s {
			bar.Add(1)
			if !yield(i, e) {
				return
			}
		}
	}
}

// N ranges over [0, n) like range does over an integer, drawing a live
// progress bar as iteration proceeds:
//
//	for i := range progress.N(n) {
//	    ...
//	}
func N(n int, opts ...Option) iter.Seq[int] {
	return func(yield func(int) bool) {
		bar := newBar(n, opts...)
		defer bar.Finish()
		for i := range n {
			bar.Add(1)
			if !yield(i) {
				return
			}
		}
	}
}

// Spinner ranges over an unbounded count 0, 1, 2, ..., drawing a live
// spinner (no percentage or ETA, since the total is unknown). Useful for
// open-ended searches; break out of the loop once you are done:
//
//	for i := range progress.Spinner() {
//	    if found(i) {
//	        break
//	    }
//	}
func Spinner(opts ...Option) iter.Seq[int] {
	return func(yield func(int) bool) {
		bar := newBar(-1, opts...)
		defer bar.Finish()
		for i := 0; ; i++ {
			bar.Add(1)
			if !yield(i) {
				return
			}
		}
	}
}

// Seq ranges over seq, drawing a progress bar that expects total elements.
// With total < 0 it draws a spinner instead, for sequences of unknown
// length.
func Seq[T any](seq iter.Seq[T], total int, opts ...Option) iter.Seq[T] {
	return func(yield func(T) bool) {
		bar := newBar(total, opts...)
		defer bar.Finish()
		for e := range seq {
			bar.Add(1)
			if !yield(e) {
				return
			}
		}
	}
}

// Stream returns s with a progress bar that advances as elements are
// pulled through it; see Seq. The bar is drawn each time the returned
// stream is consumed:
//
//	sum := stream.Sum(progress.Stream(primes.Between(2, 2_000_000), 148_933))
func Stream[T any](s stream.Stream[T], total int, opts ...Option) stream.Stream[T] {
	return stream.FromSeq(Seq(s.All(), total, opts...))
}
