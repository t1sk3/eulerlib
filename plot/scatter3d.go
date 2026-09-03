package plot

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/t1sk3/eulerlib/v2/utils"
)

// isoAngle is the classic 30° isometric projection angle.
const isoAngle = math.Pi / 6

// isometricProject maps a 3D point onto the 2D plane using a simple
// isometric (axonometric) projection: the viewer looks along the (1,1,1)
// direction, so the x, y, and z axes all remain visible and equally
// foreshortened.
func isometricProject(x, y, z float64) (sx, sy float64) {
	sx = (x - z) * math.Cos(isoAngle)
	sy = (x+z)*math.Sin(isoAngle) - y
	return
}

// Scatter3D renders points as an isometrically-projected SVG scatter plot
// and returns the SVG markup. If path is non-empty, the markup is also
// written to that file (as a .svg, viewable in any browser or image
// viewer).
//
// Depth is conveyed two ways: points nearer the viewer (along the (1,1,1)
// viewing direction) are drawn larger and are layered on top of points
// further away. A small XYZ orientation gizmo is drawn in the corner
// unless WithAxes(false) is used.
//
// Each point may carry an optional Label, drawn next to it; leave it empty
// to plot bare points. Use WithTitle, WithXLabel, WithYLabel, and
// WithZLabel to caption the graph itself (the axis labels are used only by
// the orientation gizmo, since the plot area itself is a 2D projection).
//
//	pts := []plot.Point3D[float64]{
//	    {X: 1, Y: 0, Z: 0, Label: "A"},
//	    {X: 0, Y: 1, Z: 0, Label: "B"},
//	    {X: 0, Y: 0, Z: 1, Label: "C"},
//	}
//	svg, err := plot.Scatter3D(pts, "points.svg", plot.WithTitle("unit points"))
func Scatter3D[E utils.RealNumber](points []Point3D[E], path string, opts ...Option) (string, error) {
	cfg := newConfig(opts...)

	type projected struct {
		sx, sy, depth float64
		label         string
	}

	proj := make([]projected, len(points))
	minDepth, maxDepth := math.Inf(1), math.Inf(-1)
	for i, p := range points {
		x, y, z := float64(p.X), float64(p.Y), float64(p.Z)
		sx, sy := isometricProject(x, y, z)
		depth := x + y + z // distance along the (1,1,1) viewing direction
		proj[i] = projected{sx, sy, depth, p.Label}
		minDepth = math.Min(minDepth, depth)
		maxDepth = math.Max(maxDepth, depth)
	}

	// Painter's algorithm: draw the farthest points first so nearer ones
	// are layered on top.
	sort.Slice(proj, func(i, j int) bool { return proj[i].depth < proj[j].depth })

	plotted := make([]plottedPoint, len(proj))
	for i, p := range proj {
		plotted[i] = plottedPoint{
			x:     p.sx,
			y:     p.sy,
			scale: depthScale(p.depth, minDepth, maxDepth),
			label: p.label,
		}
	}

	svg := renderSVG(plotted, cfg, axisGizmo(cfg))
	if err := writeFile(path, svg); err != nil {
		return "", err
	}
	return svg, nil
}

// depthScale maps a point's depth into a size multiplier in [0.6, 1.4], so
// nearer points render up to twice the size of farther ones.
func depthScale(depth, minDepth, maxDepth float64) float64 {
	if maxDepth == minDepth {
		return 1
	}
	t := (depth - minDepth) / (maxDepth - minDepth)
	return 0.6 + 0.8*t
}

// axisGizmo builds a small fixed-size XYZ orientation indicator anchored
// in the bottom-left corner of the canvas, independent of the data's own
// scale.
func axisGizmo(cfg config) string {
	if !cfg.showAxes {
		return ""
	}

	const armLength = 28.0
	ox, oy := 55.0, float64(cfg.height)-30.0

	labelX, labelY, labelZ := cfg.xLabel, cfg.yLabel, cfg.zLabel
	if labelX == "" {
		labelX = "X"
	}
	if labelY == "" {
		labelY = "Y"
	}
	if labelZ == "" {
		labelZ = "Z"
	}

	arms := []struct {
		x, y, z      float64
		color, label string
	}{
		{1, 0, 0, "#d62728", labelX},
		{0, 1, 0, "#2ca02c", labelY},
		{0, 0, 1, "#1f77b4", labelZ},
	}

	var b strings.Builder
	b.WriteString(`<g stroke-width="2">`)
	for _, a := range arms {
		dx, dy := isometricProject(a.x, a.y, a.z)
		ex, ey := ox+dx*armLength, oy+dy*armLength
		fmt.Fprintf(&b, `<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s"/>`, ox, oy, ex, ey, a.color)
		fmt.Fprintf(&b, `<text x="%.2f" y="%.2f" font-size="11" fill="%s">%s</text>`, ex, ey, a.color, escape(a.label))
	}
	b.WriteString(`</g>`)
	return b.String()
}
