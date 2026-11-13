package memo

import "testing"

func TestMemoize3(t *testing.T) {
	calls := 0
	var paths func(x, y, z int) int
	paths = Memoize3(func(x, y, z int) int {
		calls++
		if x == 0 && y == 0 && z == 0 {
			return 1
		}
		n := 0
		if x > 0 {
			n += paths(x-1, y, z)
		}
		if y > 0 {
			n += paths(x, y-1, z)
		}
		if z > 0 {
			n += paths(x, y, z-1)
		}
		return n
	})
	// Lattice paths in a 6×6×6 box: 18! / (6!)³.
	if got := paths(6, 6, 6); got != 17153136 {
		t.Errorf("paths(6, 6, 6) = %d", got)
	}
	if calls != 7*7*7 {
		t.Errorf("fn ran %d times, want %d", calls, 7*7*7)
	}
	// Arguments of different types, and argument order matters.
	f := Memoize3(func(a string, b int, c bool) string {
		if c {
			return a
		}
		return a[:b]
	})
	if f("hello", 2, false) != "he" || f("hello", 2, true) != "hello" || f("hello", 2, false) != "he" {
		t.Error("Memoize3 with mixed key types")
	}
}

func TestMemoizeWithClear(t *testing.T) {
	calls, offset := 0, 0
	f, clear := MemoizeWithClear(func(n int) int {
		calls++
		return n + offset
	})
	if f(1) != 1 || f(1) != 1 || f(2) != 2 || calls != 2 {
		t.Fatalf("before clear: %d calls", calls)
	}
	offset = 10
	if f(1) != 1 {
		t.Error("a cached value changed without a clear")
	}
	clear()
	if f(1) != 11 || f(1) != 11 || calls != 3 {
		t.Errorf("after clear: f(1) = %d, %d calls", f(1), calls)
	}
	clear()
	clear() // clearing an empty cache is fine
	if f(2) != 12 {
		t.Error("after a second clear")
	}
}
