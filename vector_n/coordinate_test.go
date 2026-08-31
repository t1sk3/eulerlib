package vector_n

import (
	"math"
	"testing"
)

func TestAreColinear(t *testing.T) {
	tests := []struct {
		name   string
		p1, p2 *Vector2[int]
		want   bool
	}{
		{"parallel same direction", NewVector2(1, 2), NewVector2(2, 4), true},
		{"parallel opposite direction", NewVector2(1, 2), NewVector2(-2, -4), true},
		{"zero vector", NewVector2(0, 0), NewVector2(3, 5), true},
		{"not colinear", NewVector2(1, 2), NewVector2(3, 4), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AreColinear(*tt.p1, *tt.p2); got != tt.want {
				t.Errorf("AreColinear(%v, %v) = %v, want %v", tt.p1, tt.p2, got, tt.want)
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
		name   string
		p1, p2 *HexCoordinate[int]
		want   bool
	}{
		{"parallel same direction", NewHexCoordinate(1, 2), NewHexCoordinate(2, 4), true},
		{"parallel opposite direction", NewHexCoordinate(1, -1), NewHexCoordinate(-3, 3), true},
		{"zero vector", NewHexCoordinate(0, 0), NewHexCoordinate(4, -1), true},
		{"not colinear", NewHexCoordinate(1, 0), NewHexCoordinate(0, 1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HexAreColinear(*tt.p1, *tt.p2); got != tt.want {
				t.Errorf("HexAreColinear(%v, %v) = %v, want %v", tt.p1, tt.p2, got, tt.want)
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
