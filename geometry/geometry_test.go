package geometry

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/t1sk3/eulerlib/v4/internal/testutil"
	"github.com/t1sk3/go-utils/types/vector"
)

// The cases below are ported unchanged from eulerlib v3's
// vector/coordinate_test.go, so behaviour is identical.

func v(x, y int) vector.Vec2[int] { return vector.Of2(x, y) }

func TestAreColinear(t *testing.T) {
	tests := []struct {
		name       string
		p1, p2, p3 vector.Vec2[int]
		want       bool
	}{
		{"colinear through origin", v(0, 0), v(1, 1), v(2, 2), true},
		{"colinear, none at origin", v(0, 4), v(2, 2), v(4, 0), true},
		{"not colinear", v(0, 0), v(1, 0), v(0, 1), false},
		{"duplicate point is degenerate colinear", v(1, 1), v(1, 1), v(5, 9), true},
	}
	for _, tt := range tests {
		if got := AreColinear(tt.p1, tt.p2, tt.p3); got != tt.want {
			t.Errorf("%s: AreColinear = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestHex(t *testing.T) {
	h := HexOf(2, -3)
	if h.Q != 2 || h.R != -3 || h.S() != 1 || h.Q+h.R+h.S() != 0 {
		t.Errorf("HexOf(2, -3) = %v, S = %d", h, h.S())
	}
	if h != (Hex[int]{Q: 2, R: -3}) {
		t.Error("Hex should be comparable")
	}
	if h.Add(HexOf(1, 1)) != HexOf(3, -2) || h.Sub(HexOf(1, 1)) != HexOf(1, -4) || h.Scale(2) != HexOf(4, -6) {
		t.Error("Add/Sub/Scale")
	}
	if h.String() != "Hex(2, -3)" || h.ToVec2() != v(2, -3) {
		t.Error("String/ToVec2")
	}
	m := map[Hex[int]]bool{h: true}
	if !m[HexOf(2, -3)] {
		t.Error("Hex as map key")
	}
	b, err := json.Marshal(h)
	if err != nil || string(b) != `{"q":2,"r":-3}` {
		t.Errorf("json = %s, %v", b, err)
	}
	var back Hex[int]
	if err := json.Unmarshal(b, &back); err != nil || back != h {
		t.Errorf("json round trip = %v, %v", back, err)
	}
}

func TestDistanceNeighbors(t *testing.T) {
	o := HexOf(0, 0)
	for _, n := range o.Neighbors() {
		if o.Distance(n) != 1 {
			t.Errorf("neighbor %v at distance %d", n, o.Distance(n))
		}
	}
	if HexOf(1, -3).Distance(HexOf(-2, 2)) != 5 {
		t.Error("Distance")
	}
	// Distance equals the BFS step count.
	dist := map[Hex[int]]int{o: 0}
	frontier := []Hex[int]{o}
	for d := 1; d <= 6; d++ {
		var next []Hex[int]
		for _, h := range frontier {
			for _, n := range h.Neighbors() {
				if _, ok := dist[n]; !ok {
					dist[n] = d
					next = append(next, n)
				}
			}
		}
		frontier = next
	}
	for h, d := range dist {
		if o.Distance(h) != d {
			t.Fatalf("Distance(%v) = %d, BFS %d", h, o.Distance(h), d)
		}
	}
}

func TestHexRing(t *testing.T) {
	c := HexOf(2, -1)
	for r := 0; r <= 5; r++ {
		ring := HexRing(c, r).ToSlice()
		want := max(6*r, 1)
		if len(ring) != want {
			t.Fatalf("HexRing radius %d has %d hexes, want %d", r, len(ring), want)
		}
		seen := map[Hex[int]]bool{}
		for i, h := range ring {
			if c.Distance(h) != r || seen[h] {
				t.Fatalf("radius %d: bad hex %v", r, h)
			}
			seen[h] = true
			if r > 0 && ring[(i+1)%len(ring)].Distance(h) != 1 {
				t.Fatalf("radius %d: ring not contiguous at %d", r, i)
			}
		}
	}
	if HexRing(c, -1).Count() != 0 {
		t.Error("negative radius")
	}
	testutil.Reusable(t, HexRing(c, 3), -1)
	testutil.EarlyExit(t, HexRing(c, 3), 4)
}

func TestHexAreColinear(t *testing.T) {
	tests := []struct {
		name       string
		p1, p2, p3 Hex[int]
		want       bool
	}{
		{"colinear through origin", HexOf(0, 0), HexOf(1, 2), HexOf(2, 4), true},
		{"colinear, none at origin", HexOf(1, -1), HexOf(2, -2), HexOf(4, -4), true},
		{"not colinear", HexOf(0, 0), HexOf(1, 0), HexOf(0, 1), false},
		{"duplicate point is degenerate colinear", HexOf(2, 3), HexOf(2, 3), HexOf(-1, 5), true},
	}
	for _, tt := range tests {
		if got := HexAreColinear(tt.p1, tt.p2, tt.p3); got != tt.want {
			t.Errorf("%s: HexAreColinear = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestQuadrant(t *testing.T) {
	tests := []struct {
		name string
		p    vector.Vec2[int]
		want int
	}{
		{"interior of I", v(1, 2), 1},
		{"interior of II", v(-1, 2), 2},
		{"interior of III", v(-1, -2), 3},
		{"interior of IV", v(1, -2), 4},
		{"positive x-axis is I's lower boundary", v(3, 0), 1},
		{"positive y-axis is II's lower boundary", v(0, 3), 2},
		{"negative x-axis is III's lower boundary", v(-3, 0), 3},
		{"negative y-axis is IV's lower boundary", v(0, -3), 4},
		{"origin is IV by convention", v(0, 0), 4},
	}
	for _, tt := range tests {
		if got := Quadrant(tt.p); got != tt.want {
			t.Errorf("%s: Quadrant = %d, want %d", tt.name, got, tt.want)
		}
	}
	if Quadrant(vector.Of2(-0.5, 0.25)) != 2 {
		t.Error("Quadrant on floats")
	}
}

func TestSameQuadrant(t *testing.T) {
	tests := []struct {
		name   string
		p1, p2 vector.Vec2[int]
		want   bool
	}{
		{"both quadrant I", v(1, 2), v(5, 9), true},
		{"both quadrant III", v(-1, -2), v(-5, -9), true},
		{"different quadrants", v(1, 2), v(-1, 2), false},
		{"p1 on positive y-axis", v(0, 2), v(1, 2), false},
		{"p2 on positive x-axis", v(1, 2), v(1, 0), true},
		{"both at origin", v(0, 0), v(0, 0), true},
	}
	for _, tt := range tests {
		if got := SameQuadrant(tt.p1, tt.p2); got != tt.want {
			t.Errorf("%s: SameQuadrant = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestHexSameQuadrant(t *testing.T) {
	tests := []struct {
		name   string
		p1, p2 Hex[int]
		want   bool
	}{
		{"both quadrant I", HexOf(1, 2), HexOf(5, 9), true},
		{"both quadrant III", HexOf(-1, -2), HexOf(-5, -9), true},
		{"different quadrants", HexOf(1, 2), HexOf(-1, 2), false},
		{"p1 on axis", HexOf(0, 2), HexOf(1, 2), false},
	}
	for _, tt := range tests {
		if got := HexSameQuadrant(tt.p1, tt.p2); got != tt.want {
			t.Errorf("%s: HexSameQuadrant = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCartesianHexRoundTrip(t *testing.T) {
	for _, h := range []Hex[float64]{HexOf(0.0, 0.0), HexOf(1.0, 0.0), HexOf(0.0, 1.0), HexOf(3.0, -2.0), HexOf(-4.0, 5.0)} {
		got := CartesianToHex(HexToCartesian(h))
		if math.Abs(got.Q-h.Q) > 1e-9 || math.Abs(got.R-h.R) > 1e-9 {
			t.Errorf("round trip for %v: got %v", h, got)
		}
	}
	if p := HexToCartesian(HexOf(0.0, 0.0)); p.X != 0 || p.Y != 0 {
		t.Errorf("HexToCartesian(0, 0) = %v", p)
	}
}

func TestRoundHex(t *testing.T) {
	tests := []struct {
		name         string
		in           Hex[float64]
		wantQ, wantR float64
	}{
		{"already integral", HexOf(2.0, -1.0), 2, -1},
		{"rounds toward nearest cell", HexOf(1.2, 0.41), 1, 1},
		{"fixes up largest cube error", HexOf(0.6, 0.6), 1, 0},
	}
	for _, tt := range tests {
		got := RoundHex(tt.in)
		if got.Q != tt.wantQ || got.R != tt.wantR || got.Q+got.R+got.S() != 0 {
			t.Errorf("%s: RoundHex = %v, want (%v, %v)", tt.name, got, tt.wantQ, tt.wantR)
		}
	}
	// Every point inside a cell rounds to that cell.
	for _, h := range []Hex[float64]{HexOf(0.0, 0.0), HexOf(2.0, -1.0), HexOf(-3.0, 4.0)} {
		c := HexToCartesian(h)
		for a := 0.0; a < 2*math.Pi; a += 0.3 {
			p := vector.Of2(c.X+0.8*math.Cos(a), c.Y+0.8*math.Sin(a)) // inside the inscribed circle (r = √3/2)
			if got := RoundHex(CartesianToHex(p)); got != h {
				t.Fatalf("point %v near %v rounded to %v", p, h, got)
			}
		}
	}
}
