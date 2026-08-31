package vector

import "testing"

func TestNewVector3(t *testing.T) {
	v := NewVector3(1, 2, 3)
	if v.X() != 1 || v.Y() != 2 || v.Z() != 3 {
		t.Fatalf("NewVector3 = (%d, %d, %d), want (1, 2, 3)", v.X(), v.Y(), v.Z())
	}
	if v.Len() != 3 {
		t.Errorf("Len() = %d, want 3", v.Len())
	}
}

func TestAsVector3(t *testing.T) {
	base := NewVector([]int{1, 2, 3})
	v := AsVector3(base)
	if v.X() != 1 || v.Y() != 2 || v.Z() != 3 {
		t.Fatalf("AsVector3 = (%d, %d, %d), want (1, 2, 3)", v.X(), v.Y(), v.Z())
	}
}

func TestAsVector3WrongLengthPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("AsVector3 did not panic for a vector of length 2")
		}
	}()
	base := NewVector([]int{1, 2})
	AsVector3(base)
}

func TestVector3Cross(t *testing.T) {
	// standard basis vectors: i x j = k
	i := NewVector3(1, 0, 0)
	j := NewVector3(0, 1, 0)
	got := i.Cross(*j)
	if got.X() != 0 || got.Y() != 0 || got.Z() != 1 {
		t.Fatalf("Cross(i, j) = (%d, %d, %d), want (0, 0, 1)", got.X(), got.Y(), got.Z())
	}
}

func TestVector3CrossSelfIsZero(t *testing.T) {
	v := NewVector3(1, 2, 3)
	got := v.Cross(*v)
	if got.X() != 0 || got.Y() != 0 || got.Z() != 0 {
		t.Fatalf("Cross(v, v) = (%d, %d, %d), want (0, 0, 0)", got.X(), got.Y(), got.Z())
	}
}

func TestVector3SetXSetYSetZ(t *testing.T) {
	v := NewVector3(1, 2, 3)
	v.SetX(10)
	v.SetY(20)
	v.SetZ(30)
	if v.X() != 10 || v.Y() != 20 || v.Z() != 30 {
		t.Fatalf("after SetX/SetY/SetZ = (%d, %d, %d), want (10, 20, 30)", v.X(), v.Y(), v.Z())
	}
}

func TestVector3Add(t *testing.T) {
	a := NewVector3(1, 2, 3)
	b := NewVector3(4, 5, 6)
	got := a.Add(*b)
	if got.X() != 5 || got.Y() != 7 || got.Z() != 9 {
		t.Fatalf("Add() = (%d, %d, %d), want (5, 7, 9)", got.X(), got.Y(), got.Z())
	}
}

func TestVector3Sub(t *testing.T) {
	a := NewVector3(5, 7, 9)
	b := NewVector3(1, 2, 3)
	got := a.Sub(*b)
	if got.X() != 4 || got.Y() != 5 || got.Z() != 6 {
		t.Fatalf("Sub() = (%d, %d, %d), want (4, 5, 6)", got.X(), got.Y(), got.Z())
	}
}

func TestVector3Mul(t *testing.T) {
	a := NewVector3(2, 3, 4)
	b := NewVector3(5, 6, 7)
	got := a.Mul(*b)
	if got.X() != 10 || got.Y() != 18 || got.Z() != 28 {
		t.Fatalf("Mul() = (%d, %d, %d), want (10, 18, 28)", got.X(), got.Y(), got.Z())
	}
}

func TestVector3Div(t *testing.T) {
	a := NewVector3(10, 20, 30)
	b := NewVector3(2, 5, 3)
	got := a.Div(*b)
	if got.X() != 5 || got.Y() != 4 || got.Z() != 10 {
		t.Fatalf("Div() = (%d, %d, %d), want (5, 4, 10)", got.X(), got.Y(), got.Z())
	}
}
