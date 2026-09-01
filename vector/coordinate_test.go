package vector

import (
	"math"
	"testing"
)

func TestAreColinear(t *testing.T) {
	tests := []struct {
		name       string
		p1, p2, p3 *Vector2[int]
		want       bool
	}{
		{"colinear through origin", NewVector2(0, 0), NewVector2(1, 1), NewVector2(2, 2), true},
		{"colinear, none at origin", NewVector2(0, 4), NewVector2(2, 2), NewVector2(4, 0), true},
		{"not colinear", NewVector2(0, 0), NewVector2(1, 0), NewVector2(0, 1), false},
		{"duplicate point is degenerate colinear", NewVector2(1, 1), NewVector2(1, 1), NewVector2(5, 9), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AreColinear(*tt.p1, *tt.p2, *tt.p3); got != tt.want {
				t.Errorf("AreColinear(%v, %v, %v) = %v, want %v", tt.p1, tt.p2, tt.p3, got, tt.want)
			}
		})
	}
}

func TestNewHexCoordinate(t *testing.T) {
	h := NewHexCoordinate(2, -3)
	if h.Q() != 2 {
		t.Errorf("Q() = %d, want 2", h.Q())
	}
	if h.R() != -3 {
		t.Errorf("R() = %d, want -3", h.R())
	}
	if h.S() != 1 {
		t.Errorf("S() = %d, want 1", h.S())
	}
	if h.Q()+h.R()+h.S() != 0 {
		t.Errorf("q+r+s = %d, want 0", h.Q()+h.R()+h.S())
	}
}

func TestHexAreColinear(t *testing.T) {
	tests := []struct {
		name       string
		p1, p2, p3 *HexCoordinate[int]
		want       bool
	}{
		{"colinear through origin", NewHexCoordinate(0, 0), NewHexCoordinate(1, 2), NewHexCoordinate(2, 4), true},
		{"colinear, none at origin", NewHexCoordinate(1, -1), NewHexCoordinate(2, -2), NewHexCoordinate(4, -4), true},
		{"not colinear", NewHexCoordinate(0, 0), NewHexCoordinate(1, 0), NewHexCoordinate(0, 1), false},
		{"duplicate point is degenerate colinear", NewHexCoordinate(2, 3), NewHexCoordinate(2, 3), NewHexCoordinate(-1, 5), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HexAreColinear(*tt.p1, *tt.p2, *tt.p3); got != tt.want {
				t.Errorf("HexAreColinear(%v, %v, %v) = %v, want %v", tt.p1, tt.p2, tt.p3, got, tt.want)
			}
		})
	}
}

func TestQuadrant(t *testing.T) {
	tests := []struct {
		name string
		p    *Vector2[int]
		want int
	}{
		{"quadrant I", NewVector2(1, 1), 1},
		{"quadrant II", NewVector2(-1, 1), 2},
		{"quadrant III", NewVector2(-1, -1), 3},
		{"quadrant IV", NewVector2(1, -1), 4},
		{"on x-axis", NewVector2(3, 0), 0},
		{"on y-axis", NewVector2(0, -2), 0},
		{"origin", NewVector2(0, 0), 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Quadrant(*tt.p); got != tt.want {
				t.Errorf("Quadrant(%v) = %d, want %d", tt.p, got, tt.want)
			}
		})
	}
}

func TestSameQuadrant(t *testing.T) {
	tests := []struct {
		name   string
		p1, p2 *Vector2[int]
		want   bool
	}{
		{"both quadrant I", NewVector2(1, 2), NewVector2(5, 9), true},
		{"both quadrant III", NewVector2(-1, -2), NewVector2(-5, -9), true},
		{"different quadrants", NewVector2(1, 2), NewVector2(-1, 2), false},
		{"p1 on axis", NewVector2(0, 2), NewVector2(1, 2), false},
		{"p2 on axis", NewVector2(1, 2), NewVector2(1, 0), false},
		{"both at origin", NewVector2(0, 0), NewVector2(0, 0), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SameQuadrant(*tt.p1, *tt.p2); got != tt.want {
				t.Errorf("SameQuadrant(%v, %v) = %v, want %v", tt.p1, tt.p2, got, tt.want)
			}
		})
	}
}

func TestHexSameQuadrant(t *testing.T) {
	tests := []struct {
		name   string
		p1, p2 *HexCoordinate[int]
		want   bool
	}{
		{"both quadrant I", NewHexCoordinate(1, 2), NewHexCoordinate(5, 9), true},
		{"both quadrant III", NewHexCoordinate(-1, -2), NewHexCoordinate(-5, -9), true},
		{"different quadrants", NewHexCoordinate(1, 2), NewHexCoordinate(-1, 2), false},
		{"p1 on axis", NewHexCoordinate(0, 2), NewHexCoordinate(1, 2), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HexSameQuadrant(*tt.p1, *tt.p2); got != tt.want {
				t.Errorf("HexSameQuadrant(%v, %v) = %v, want %v", tt.p1, tt.p2, got, tt.want)
			}
		})
	}
}

func TestCartesianHexRoundTrip(t *testing.T) {
	hexes := []*HexCoordinate[float64]{
		NewHexCoordinate(0.0, 0.0),
		NewHexCoordinate(1.0, 0.0),
		NewHexCoordinate(0.0, 1.0),
		NewHexCoordinate(3.0, -2.0),
		NewHexCoordinate(-4.0, 5.0),
	}

	for _, h := range hexes {
		p := HexToCartesian(*h)
		got := CartesianToHex(*p)
		if math.Abs(got.Q()-h.Q()) > 1e-9 || math.Abs(got.R()-h.R()) > 1e-9 {
			t.Errorf("round-trip for (%v, %v): got (%v, %v)", h.Q(), h.R(), got.Q(), got.R())
		}
	}
}

func TestHexToCartesianOrigin(t *testing.T) {
	h := NewHexCoordinate(0.0, 0.0)
	p := HexToCartesian(*h)
	if p.X() != 0 || p.Y() != 0 {
		t.Errorf("HexToCartesian(0,0) = (%v, %v), want (0, 0)", p.X(), p.Y())
	}
}

func TestRoundHex(t *testing.T) {
	tests := []struct {
		name  string
		in    *HexCoordinate[float64]
		wantQ float64
		wantR float64
	}{
		{"already integral", NewHexCoordinate(2.0, -1.0), 2, -1},
		{"rounds toward nearest cell", NewHexCoordinate(1.2, 0.41), 1, 1},
		{"fixes up largest cube error", NewHexCoordinate(0.6, 0.6), 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RoundHex(*tt.in)
			if got.Q() != tt.wantQ || got.R() != tt.wantR {
				t.Errorf("RoundHex(%v, %v) = (%v, %v), want (%v, %v)", tt.in.Q(), tt.in.R(), got.Q(), got.R(), tt.wantQ, tt.wantR)
			}
			if got.Q()+got.R()+got.S() != 0 {
				t.Errorf("rounded q+r+s = %v, want 0", got.Q()+got.R()+got.S())
			}
		})
	}
}
