package plot

import (
	"iter"

	"github.com/t1sk3/go-utils/types/constraints"
	"github.com/t1sk3/go-utils/types/pair"
	"github.com/t1sk3/go-utils/types/vector"
)

// PointFromVec2 returns v as a labeled 2D point.
func PointFromVec2[E constraints.Number](v vector.Vec2[E], label string) Point2D[E] {
	return Point2D[E]{X: v.X, Y: v.Y, Label: label}
}

// PointFromVec3 returns v as a labeled 3D point.
func PointFromVec3[E constraints.Number](v vector.Vec3[E], label string) Point3D[E] {
	return Point3D[E]{X: v.X, Y: v.Y, Z: v.Z, Label: label}
}

// PointFromPair returns the pair (x, y) as a labeled 2D point.
func PointFromPair[E constraints.Number](p pair.Pair[E, E], label string) Point2D[E] {
	return Point2D[E]{X: p.First, Y: p.Second, Label: label}
}

// Points2D collects seq into unlabeled 2D points. Pass s.All() to plot a
// stream.
func Points2D[E constraints.Number](seq iter.Seq[vector.Vec2[E]]) []Point2D[E] {
	var out []Point2D[E]
	for v := range seq {
		out = append(out, PointFromVec2(v, ""))
	}
	return out
}

// Points3D collects seq into unlabeled 3D points. Pass s.All() to plot a
// stream.
func Points3D[E constraints.Number](seq iter.Seq[vector.Vec3[E]]) []Point3D[E] {
	var out []Point3D[E]
	for v := range seq {
		out = append(out, PointFromVec3(v, ""))
	}
	return out
}
