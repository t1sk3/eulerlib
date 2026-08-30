package vector_n

import "testing"

func TestNewVector2(t *testing.T) {
	v := NewVector2(3, 4)
	if v.X() != 3 {
		t.Errorf("X() = %d, want 3", v.X())
	}
	if v.Y() != 4 {
		t.Errorf("Y() = %d, want 4", v.Y())
	}
	if v.Len() != 2 {
		t.Errorf("Len() = %d, want 2", v.Len())
	}
}

func TestAsVector2(t *testing.T) {
	base := NewVector([]int{5, 6})
	v := AsVector2(base)
	if v.X() != 5 || v.Y() != 6 {
		t.Fatalf("AsVector2 = (%d, %d), want (5, 6)", v.X(), v.Y())
	}

	// shares backing storage with the original Vector
	v.SetX(99)
	if base.At(0) != 99 {
		t.Fatalf("AsVector2 does not share storage: base.At(0) = %d, want 99", base.At(0))
	}
}

func TestAsVector2WrongLengthPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("AsVector2 did not panic for a vector of length 3")
		}
	}()
	base := NewVector([]int{1, 2, 3})
	AsVector2(base)
}

func TestVector2SetXSetY(t *testing.T) {
	v := NewVector2(1, 2)
	v.SetX(10)
	v.SetY(20)
	if v.X() != 10 || v.Y() != 20 {
		t.Fatalf("after SetX/SetY = (%d, %d), want (10, 20)", v.X(), v.Y())
	}
}

func TestVector2Cross(t *testing.T) {
	a := NewVector2(1, 2)
	b := NewVector2(3, 4)
	got := a.Cross(*b)
	want := 1*4 - 2*3
	if got != want {
		t.Fatalf("Cross() = %d, want %d", got, want)
	}
}

func TestVector2Perp(t *testing.T) {
	v := NewVector2(1, 2)
	got := v.Perp()
	if got.X() != -2 || got.Y() != 1 {
		t.Fatalf("Perp() = (%d, %d), want (-2, 1)", got.X(), got.Y())
	}
}

func TestVector2Add(t *testing.T) {
	a := NewVector2(1, 2)
	b := NewVector2(3, 4)
	got := a.Add(*b)
	if got.X() != 4 || got.Y() != 6 {
		t.Fatalf("Add() = (%d, %d), want (4, 6)", got.X(), got.Y())
	}
}

func TestVector2Sub(t *testing.T) {
	a := NewVector2(4, 6)
	b := NewVector2(1, 2)
	got := a.Sub(*b)
	if got.X() != 3 || got.Y() != 4 {
		t.Fatalf("Sub() = (%d, %d), want (3, 4)", got.X(), got.Y())
	}
}

func TestVector2Mul(t *testing.T) {
	a := NewVector2(2, 3)
	b := NewVector2(4, 5)
	got := a.Mul(*b)
	if got.X() != 8 || got.Y() != 15 {
		t.Fatalf("Mul() = (%d, %d), want (8, 15)", got.X(), got.Y())
	}
}

func TestVector2Div(t *testing.T) {
	a := NewVector2(10, 20)
	b := NewVector2(2, 5)
	got := a.Div(*b)
	if got.X() != 5 || got.Y() != 4 {
		t.Fatalf("Div() = (%d, %d), want (5, 4)", got.X(), got.Y())
	}
}
