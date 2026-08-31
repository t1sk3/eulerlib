package vector_n

import (
	"math"

	"github.com/t1sk3/eulerlib/v2/utils"
)

// AreColinear reports whether p1 and p2 are colinear, i.e. p1, p2, and the
// origin all lie on a common line — equivalently, one vector is a scalar
// multiple of the other (either may be the zero vector). It is defined as
// the perp-dot product (Vector2.Cross) being zero.
func AreColinear[E utils.SignedNumber](p1, p2 Vector2[E]) bool {
	return p1.Cross(p2) == 0
}

// HexCoordinate is a point in axial hexagonal coordinates (q, r), the
// standard two-axis system used to address cells of a hexagonal grid. It
// embeds Vector2 to reuse vector arithmetic (Add, Sub, Cross, ...); use
// Q()/R() rather than the inherited X()/Y() for hex-grid clarity.
type HexCoordinate[E utils.SignedNumber] struct {
	Vector2[E]
}

// NewHexCoordinate creates an axial hex coordinate from q and r.
func NewHexCoordinate[E utils.SignedNumber](q, r E) *HexCoordinate[E] {
	return &HexCoordinate[E]{*NewVector2(q, r)}
}

// Q returns the q axial coordinate.
func (h *HexCoordinate[E]) Q() E { return h.X() }

// R returns the r axial coordinate.
func (h *HexCoordinate[E]) R() E { return h.Y() }

// S returns the third cube coordinate, s = -q-r. Cube coordinates satisfy
// q+r+s = 0, which is useful for rounding and distance calculations.
func (h *HexCoordinate[E]) S() E { return -h.Q() - h.R() }

// HexAreColinear reports whether p1 and p2 are colinear on the hex grid,
// i.e. p1, p2, and the grid origin lie on a common line.
//
// No conversion to cartesian coordinates is needed: axial-to-cartesian
// (see HexToCartesian) is a fixed invertible linear map A, and for any
// linear map, cross(A*u, A*v) = det(A) * cross(u, v). Since det(A) != 0,
// the perp-dot product of the axial coordinates is zero exactly when the
// corresponding cartesian points are colinear.
func HexAreColinear[E utils.SignedNumber](p1, p2 HexCoordinate[E]) bool {
	return AreColinear(p1.Vector2, p2.Vector2)
}

// CartesianToHex converts a cartesian point to fractional axial hex
// coordinates, assuming unit-size (circumradius 1), pointy-top hexagons
// centered on the origin. The result is generally fractional — round it
// with RoundHex to snap it to the containing hex cell.
func CartesianToHex[E utils.Float](p Vector2[E]) *HexCoordinate[E] {
	sqrt3 := E(math.Sqrt(3))
	q := sqrt3/3*p.X() - p.Y()/3
	r := 2 * p.Y() / 3
	return NewHexCoordinate(q, r)
}

// HexToCartesian converts axial hex coordinates to a cartesian point, the
// inverse of CartesianToHex (unit-size, pointy-top hexagons centered on the
// origin).
func HexToCartesian[E utils.Float](h HexCoordinate[E]) *Vector2[E] {
	sqrt3 := E(math.Sqrt(3))
	x := sqrt3*h.Q() + sqrt3/2*h.R()
	y := 3 * h.R() / 2
	return NewVector2(x, y)
}

// RoundHex rounds fractional axial hex coordinates (as produced by
// CartesianToHex) to the nearest hex cell. It rounds each cube coordinate
// independently, then fixes up the component with the largest rounding
// error so that q+r+s stays 0.
func RoundHex[E utils.Float](h HexCoordinate[E]) *HexCoordinate[E] {
	q, r, s := h.Q(), h.R(), h.S()

	rq := E(math.Round(float64(q)))
	rr := E(math.Round(float64(r)))
	rs := E(math.Round(float64(s)))

	dq, dr, ds := absFloat(rq-q), absFloat(rr-r), absFloat(rs-s)

	switch {
	case dq > dr && dq > ds:
		rq = -rr - rs
	case dr > ds:
		rr = -rq - rs
	default:
		rs = -rq - rr
	}

	return NewHexCoordinate(rq, rr)
}

// absFloat returns the absolute value of x.
func absFloat[E utils.Float](x E) E {
	if x < 0 {
		return -x
	}
	return x
}
