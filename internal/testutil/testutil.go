// Package testutil has the checks every lazy eulerlib sequence must pass:
// it can be consumed more than once with the same result, and breaking out
// of it early leaves no goroutine behind.
package testutil

import (
	"bytes"
	"reflect"
	"runtime"
	"runtime/pprof"
	"strings"
	"testing"
	"time"

	"github.com/t1sk3/go-utils/stream"
)

// Reusable fails t unless consuming the first n elements of s twice gives
// the same elements both times. For a finite stream, n < 0 consumes all of
// it.
func Reusable[E any](t testing.TB, s stream.Stream[E], n int) {
	t.Helper()
	take := func() []E {
		if n < 0 {
			return s.ToSlice()
		}
		return s.Limit(n).ToSlice()
	}
	first, second := take(), take()
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("stream is not reusable: first pass %v, second pass %v", first, second)
	}
}

// EarlyExit ranges over s, breaks out after k elements, and fails t if a
// goroutine is left behind: runtime.NumGoroutine must return to its
// starting value, and the goroutineleak profile must not list any eulerlib
// frame.
func EarlyExit[E any](t testing.TB, s stream.Stream[E], k int) {
	t.Helper()
	before := runtime.NumGoroutine()
	i := 0
	for range s.All() {
		i++
		if i >= k {
			break
		}
	}
	if i < k {
		t.Fatalf("stream ended after %d elements, wanted at least %d to test early exit", i, k)
	}
	NoLeaks(t, before)
}

// NoLeaks fails t if the number of goroutines does not settle back to
// before, or if the goroutineleak profile names eulerlib code.
func NoLeaks(t testing.TB, before int) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for runtime.NumGoroutine() > before && time.Now().Before(deadline) {
		runtime.Gosched()
		time.Sleep(time.Millisecond)
	}
	if now := runtime.NumGoroutine(); now > before {
		t.Fatalf("goroutine leak: %d goroutines before, %d after", before, now)
	}
	if p := pprof.Lookup("goroutineleak"); p != nil {
		runtime.GC()
		var buf bytes.Buffer
		if err := p.WriteTo(&buf, 1); err == nil && strings.Contains(buf.String(), "t1sk3/eulerlib") {
			t.Fatalf("goroutineleak profile lists eulerlib frames:\n%s", buf.String())
		}
	}
}
