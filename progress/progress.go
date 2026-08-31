// Package progress renders a live progress bar while ranging over a slice
// or a count of iterations, using Go 1.23 range-over-func. Rendering is
// delegated to github.com/schollz/progressbar/v3.
package progress

import (
	"fmt"
	"io"
	"iter"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
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
	return progressbar.OptionSetWriter(w)
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
		progressbar.OptionOnCompletion(func() { fmt.Fprintln(os.Stderr) }),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
	}
	return progressbar.NewOptions(total, append(defaults, opts...)...)
}

// ProgressBar ranges over s exactly like range does, drawing a live progress
// bar (default: to stderr) as iteration proceeds:
//
//	for i, e := range progress.ProgressBar(s) {
//	    ...
//	}
//
// Breaking out of the loop early still finalizes the bar at whatever count
// was reached.
func ProgressBar[T any](s []T, opts ...Option) iter.Seq2[int, T] {
	bar := newBar(len(s), opts...)
	return func(yield func(int, T) bool) {
		defer bar.Finish()
		for i, e := range s {
			if !yield(i, e) {
				return
			}
			bar.Add(1)
		}
	}
}

// ProgressBarN ranges over [0, n) like range does over an integer, drawing a
// live progress bar as iteration proceeds:
//
//	for i := range progress.ProgressBarN(n) {
//	    ...
//	}
func ProgressBarN(n int, opts ...Option) iter.Seq[int] {
	bar := newBar(n, opts...)
	return func(yield func(int) bool) {
		defer bar.Finish()
		for i := 0; i < n; i++ {
			if !yield(i) {
				return
			}
			bar.Add(1)
		}
	}
}
