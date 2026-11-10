// Package geometry provides 2D coordinate helpers on top of go-utils'
// vector.Vec2: colinearity, quadrants, and an axial hexagonal grid
// coordinate, Hex.
package geometry

import (
	"fmt"
	"math"

	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/constraints"
	"github.com/t1sk3/go-utils/types/vector"
)

// signedReal is the element types with a sign that the quadrant and hex
// helpers accept: signed integers and floats.
type signedReal interface {
	constraints.Signed | constraints.Float
}

// AreColinear reports whether p1, p2 and p3 lie on one line, i.e. whether
// p2-p1 and p3-p1 are parallel (their perp-dot product is 0). For float
// vectors, compare (p2-p1).Cross(p3-p1) to 0 with a tolerance instead.
func AreColinear[E constraints.Signed](p1, p2, p3 vector.Vec2[E]) bool {
	return p2.Sub(p1).Cross(p3.Sub(p1)) == 0
}

// Quadrant returns the quadrant 1–4 of p, counter-clockwise from the
// positive x-axis: I is [0°, 90°), II [90°, 180°), III [180°, 270°) and IV
// [270°, 360°). Each quadrant includes its starting axis and excludes its
// ending one, so the positive x-axis is in I, the positive y-axis in II,
// the negative x-axis in III and the negative y-axis in IV. The origin is
// in IV by convention.
func Quadrant[E signedReal](p vector.Vec2[E]) int {
	switch {
	case p.X > 0 && p.Y >= 0:
		return 1
	case p.X <= 0 && p.Y > 0:
		return 2
	case p.X < 0 && p.Y <= 0:
		return 3
	default:
		return 4
	}
}

// SameQuadrant reports whether p1 and p2 are in the same quadrant, using
// Quadrant's boundary rules.
func SameQuadrant[E signedReal](p1, p2 vector.Vec2[E]) bool {
	return Quadrant(p1) == Quadrant(p2)
}

// Hex is a cell of a hexagonal grid in axial coordinates (Q, R). The third
// cube coordinate is S = -Q-R. Hex is a comparable value type, so it can be
// used as a map key, and marshals to JSON as {"q": ..., "r": ...}.
type Hex[E signedReal] struct {
	Q E `json:"q"`
	R E `json:"r"`
}

// HexOf returns the hex (q, r).
func HexOf[E signedReal](q, r E) Hex[E] { return Hex[E]{q, r} }

// S returns the third cube coordinate, -Q-R, so that Q+R+S = 0.
func (h Hex[E]) S() E { return -h.Q - h.R }

// Add returns h + o.
func (h Hex[E]) Add(o Hex[E]) Hex[E] { return Hex[E]{h.Q + o.Q, h.R + o.R} }

// Sub returns h - o.
func (h Hex[E]) Sub(o Hex[E]) Hex[E] { return Hex[E]{h.Q - o.Q, h.R - o.R} }

// Scale returns (kQ, kR).
func (h Hex[E]) Scale(k E) Hex[E] { return Hex[E]{k * h.Q, k * h.R} }

// Distance returns the number of steps between h and o on the grid:
// (|dq| + |dr| + |ds|) / 2.
func (h Hex[E]) Distance(o Hex[E]) E {
	d := h.Sub(o)
	return (abs(d.Q) + abs(d.R) + abs(d.S())) / 2
}

// hexDirections are the six unit steps, counter-clockwise from +Q.
var hexDirections = [6][2]int8{{1, 0}, {1, -1}, {0, -1}, {-1, 0}, {-1, 1}, {0, 1}}

// Neighbors returns the six hexes adjacent to h, counter-clockwise starting
// from h + (1, 0).
func (h Hex[E]) Neighbors() [6]Hex[E] {
	var out [6]Hex[E]
	for i, d := range hexDirections {
		out[i] = h.Add(Hex[E]{E(d[0]), E(d[1])})
	}
	return out
}

// ToVec2 returns (Q, R) as a vector, for vector arithmetic. It is not a
// cartesian position; use HexToCartesian for that.
func (h Hex[E]) ToVec2() vector.Vec2[E] { return vector.Of2(h.Q, h.R) }

// String formats h as "Hex(q, r)".
func (h Hex[E]) String() string { return fmt.Sprintf("Hex(%v, %v)", h.Q, h.R) }

// HexAreColinear reports whether three hexes lie on one line of the grid.
// No conversion to cartesian coordinates is needed: axial-to-cartesian is
// an invertible linear map A, and cross(Au, Av) = det(A)·cross(u, v), so
// the axial differences are parallel exactly when the cartesian ones are.
func HexAreColinear[E constraints.Signed](p1, p2, p3 Hex[E]) bool {
	return AreColinear(p1.ToVec2(), p2.ToVec2(), p3.ToVec2())
}

// HexSameQuadrant reports whether p1 and p2 are in the same quadrant of the
// (Q, R) axial plane, with Quadrant's boundary rules. The axial axes are 60°
// apart, so these are not the cartesian quadrants of HexToCartesian's
// output; they partition the grid in their own right.
func HexSameQuadrant[E constraints.Signed](p1, p2 Hex[E]) bool {
	return SameQuadrant(p1.ToVec2(), p2.ToVec2())
}

// CartesianToHex converts a cartesian point to fractional axial
// coordinates, for unit-size (circumradius 1) pointy-top hexagons centered
// on the origin. Round the result with RoundHex to get the containing cell.
func CartesianToHex[E constraints.Float](p vector.Vec2[E]) Hex[E] {
	sqrt3 := E(math.Sqrt(3))
	return Hex[E]{sqrt3/3*p.X - p.Y/3, 2 * p.Y / 3}
}

// HexToCartesian converts axial coordinates to the cartesian center of the
// cell, the inverse of CartesianToHex.
func HexToCartesian[E constraints.Float](h Hex[E]) vector.Vec2[E] {
	sqrt3 := E(math.Sqrt(3))
	return vector.Of2(sqrt3*h.Q+sqrt3/2*h.R, 3*h.R/2)
}

// RoundHex rounds fractional axial coordinates to the nearest cell: it
// rounds each cube coordinate, then recomputes the one with the largest
// rounding error so that Q+R+S stays 0.
func RoundHex[E constraints.Float](h Hex[E]) Hex[E] {
	q, r, s := h.Q, h.R, h.S()
	rq, rr, rs := E(math.Round(float64(q))), E(math.Round(float64(r))), E(math.Round(float64(s)))
	dq, dr, ds := abs(rq-q), abs(rr-r), abs(rs-s)
	switch {
	case dq > dr && dq > ds:
		rq = -rr - rs
	case dr > ds:
		rr = -rq - rs
	}
	return Hex[E]{rq, rr}
}

// HexRing returns the hexes at exactly radius steps from center, walking
// counter-clockwise; 6·radius of them, or just center for radius 0, and
// none for radius < 0.
func HexRing[E constraints.Signed](center Hex[E], radius int) stream.Stream[Hex[E]] {
	return stream.FromSeq(func(yield func(Hex[E]) bool) {
		if radius < 0 {
			return
		}
		if radius == 0 {
			yield(center)
			return
		}
		d4 := hexDirections[4]
		h := center.Add(Hex[E]{E(d4[0]), E(d4[1])}.Scale(E(radius)))
		for _, d := range hexDirections {
			for range radius {
				if !yield(h) {
					return
				}
				h = h.Add(Hex[E]{E(d[0]), E(d[1])})
			}
		}
	})
}

func abs[E signedReal](x E) E {
	if x < 0 {
		return -x
	}
	return x
}
