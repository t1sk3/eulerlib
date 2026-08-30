package vector_n

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

func TestVector3Add(t *testing.T) {
	a := NewVector3(1, 2, 3)
	b := NewVector3(4, 5, 6)
	got := a.Add(*b)
	if got.X() != 5 || got.Y() != 7 || got.Z() != 9 {
		t.Fatalf("Add() = (%d, %d, %d), want (5, 7, 9)", got.X(), got.Y(), got.Z())
	}
}
