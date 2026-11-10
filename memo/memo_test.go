package memo

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestMemoizeCachesResult(t *testing.T) {
	calls := 0
	square := Memoize(func(n int) int {
		calls++
		return n * n
	})

	if got := square(4); got != 16 {
		t.Errorf("square(4) = %d, want 16", got)
	}
	if got := square(4); got != 16 {
		t.Errorf("square(4) = %d, want 16", got)
	}
	if got := square(5); got != 25 {
		t.Errorf("square(5) = %d, want 25", got)
	}
	if calls != 2 {
		t.Errorf("fn called %d times, want 2 (one per distinct argument)", calls)
	}
}

func TestMemoizeRecursive(t *testing.T) {
	calls := 0
	var fib func(n int) int64
	fib = Memoize(func(n int) int64 {
		calls++
		if n < 2 {
			return int64(n)
		}
		return fib(n-1) + fib(n-2)
	})

	if got := fib(30); got != 832040 {
		t.Errorf("fib(30) = %d, want 832040", got)
	}
	if calls != 31 {
		t.Errorf("fn called %d times, want 31 (one per distinct n)", calls)
	}
}

func TestMemoize2CachesResult(t *testing.T) {
	calls := 0
	add := Memoize2(func(a, b int) int {
		calls++
		return a + b
	})

	if got := add(2, 3); got != 5 {
		t.Errorf("add(2, 3) = %d, want 5", got)
	}
	if got := add(2, 3); got != 5 {
		t.Errorf("add(2, 3) = %d, want 5", got)
	}
	if got := add(3, 2); got != 5 {
		t.Errorf("add(3, 2) = %d, want 5", got)
	}
	if calls != 2 {
		t.Errorf("fn called %d times, want 2 (one per distinct (a,b) pair)", calls)
	}
}

func TestMemoizeConcurrentRecursive(t *testing.T) {
	var calls atomic.Int64
	var fib func(n int) int64
	fib = MemoizeConcurrent(func(n int) int64 {
		calls.Add(1)
		if n < 2 {
			return int64(n)
		}
		return fib(n-1) + fib(n-2)
	})
	var wg sync.WaitGroup
	results := make([]int64, 8)
	for g := range 8 {
		wg.Go(func() { results[g] = fib(90) })
	}
	wg.Wait()
	for _, r := range results {
		if r != 2880067194370816120 {
			t.Fatalf("fib(90) = %d", r)
		}
	}
	// Each of the 91 keys is computed at least once, and at most once per
	// goroutine.
	if c := calls.Load(); c < 91 || c > 8*91 {
		t.Errorf("fn called %d times", c)
	}
	before := calls.Load()
	fib(90)
	if calls.Load() != before {
		t.Error("cached call recomputed")
	}
}
