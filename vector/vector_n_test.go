package vector

import "testing"

func assertElements(t *testing.T, got []int, want []int) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length = %d, want %d (got %v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("element[%d] = %d, want %d (got %v, want %v)", i, got[i], want[i], got, want)
		}
	}
}

func TestNewVector(t *testing.T) {
	v := NewVector([]int{1, 2, 3})
	assertElements(t, v.Elements(), []int{1, 2, 3})
}

func TestNewVectorCopiesInput(t *testing.T) {
	src := []int{1, 2, 3}
	v := NewVector(src)
	src[0] = 99
	if v.At(0) != 1 {
		t.Fatalf("NewVector did not copy input slice: v.At(0) = %d, want 1", v.At(0))
	}
}

func TestVectorLen(t *testing.T) {
	v := NewVector([]int{1, 2, 3, 4})
	if got := v.Len(); got != 4 {
		t.Fatalf("Len() = %d, want 4", got)
	}
}

func TestVectorAt(t *testing.T) {
	v := NewVector([]int{10, 20, 30})
	for i, want := range []int{10, 20, 30} {
		if got := v.At(i); got != want {
			t.Errorf("At(%d) = %d, want %d", i, got, want)
		}
	}
}

func TestVectorSet(t *testing.T) {
	v := NewVector([]int{1, 2, 3})
	v.Set(1, 42)
	if got := v.At(1); got != 42 {
		t.Fatalf("after Set(1, 42), At(1) = %d, want 42", got)
	}
}

func TestVectorElementsIsCopy(t *testing.T) {
	v := NewVector([]int{1, 2, 3})
	e := v.Elements()
	e[0] = 99
	if v.At(0) != 1 {
		t.Fatalf("Elements() did not return a copy: v.At(0) = %d, want 1", v.At(0))
	}
}

func TestVectorSwap(t *testing.T) {
	v := NewVector([]int{1, 2, 3})
	v.Swap(0, 2)
	assertElements(t, v.Elements(), []int{3, 2, 1})
}

func TestVectorAdd(t *testing.T) {
	a := NewVector([]int{1, 2, 3})
	b := NewVector([]int{4, 5, 6})
	got := a.Add(*b)
	assertElements(t, got.Elements(), []int{5, 7, 9})
	// originals unchanged
	assertElements(t, a.Elements(), []int{1, 2, 3})
	assertElements(t, b.Elements(), []int{4, 5, 6})
}

func TestVectorSub(t *testing.T) {
	a := NewVector([]int{5, 7, 9})
	b := NewVector([]int{4, 5, 6})
	got := a.Sub(*b)
	assertElements(t, got.Elements(), []int{1, 2, 3})
}

func TestVectorDot(t *testing.T) {
	a := NewVector([]int{1, 2, 3})
	b := NewVector([]int{4, 5, 6})
	got := a.Dot(*b)
	want := 1*4 + 2*5 + 3*6
	if got != want {
		t.Fatalf("Dot() = %d, want %d", got, want)
	}
}

func TestVectorMul(t *testing.T) {
	a := NewVector([]int{1, 2, 3})
	b := NewVector([]int{4, 5, 6})
	got := a.Mul(*b)
	assertElements(t, got.Elements(), []int{4, 10, 18})
}

func TestVectorDiv(t *testing.T) {
	a := NewVector([]int{10, 20, 30})
	b := NewVector([]int{2, 5, 3})
	got := a.Div(*b)
	assertElements(t, got.Elements(), []int{5, 4, 10})
}

func TestVectorMismatchedLengthsPanic(t *testing.T) {
	a := NewVector([]int{1, 2, 3})
	b := NewVector([]int{1, 2})

	tests := []struct {
		name string
		fn   func()
	}{
		{"Add", func() { a.Add(*b) }},
		{"Sub", func() { a.Sub(*b) }},
		{"Dot", func() { a.Dot(*b) }},
		{"Mul", func() { a.Mul(*b) }},
		{"Div", func() { a.Div(*b) }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("%s did not panic on mismatched lengths", tc.name)
				}
			}()
			tc.fn()
		})
	}
}
