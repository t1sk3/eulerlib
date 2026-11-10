package plot

import (
	"slices"
	"testing"

	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/pair"
	"github.com/t1sk3/go-utils/types/vector"
)

func TestConverters(t *testing.T) {
	if got := PointFromVec2(vector.Of2(1, 2), "a"); got != (Point2D[int]{X: 1, Y: 2, Label: "a"}) {
		t.Errorf("PointFromVec2 = %v", got)
	}
	if got := PointFromVec3(vector.Of3(1.5, 2, 3), "b"); got != (Point3D[float64]{X: 1.5, Y: 2, Z: 3, Label: "b"}) {
		t.Errorf("PointFromVec3 = %v", got)
	}
	if got := PointFromPair(pair.Of(4, 5), ""); got != (Point2D[int]{X: 4, Y: 5}) {
		t.Errorf("PointFromPair = %v", got)
	}
	vs := stream.Range(0, 4).Map(func(i int) vector.Vec2[int] { return vector.Of2(i, i*i) })
	if got := Points2D(vs.All()); !slices.Equal(got, []Point2D[int]{{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 4}, {X: 3, Y: 9}}) {
		t.Errorf("Points2D = %v", got)
	}
	ws := stream.Of(vector.Of3(1, 2, 3))
	if got := Points3D(ws.All()); !slices.Equal(got, []Point3D[int]{{X: 1, Y: 2, Z: 3}}) {
		t.Errorf("Points3D = %v", got)
	}
	if _, err := Scatter2D(Points2D(vs.All()), ""); err != nil {
		t.Errorf("Scatter2D of converted points: %v", err)
	}
}
