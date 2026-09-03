package plot

import "github.com/t1sk3/eulerlib/v2/utils"

// Scatter2D renders points as an SVG scatter plot and returns the SVG
// markup. If path is non-empty, the markup is also written to that file
// (as a .svg, viewable in any browser or image viewer).
//
// Each point may carry an optional Label, drawn next to it; leave it empty
// to plot bare points. Use WithTitle, WithXLabel, and WithYLabel to
// caption the graph itself.
//
//	pts := []plot.Point2D[int]{
//	    {X: 1, Y: 1, Label: "A"},
//	    {X: 2, Y: 4},
//	    {X: 3, Y: 9, Label: "C"},
//	}
//	svg, err := plot.Scatter2D(pts, "points.svg", plot.WithTitle("y = x^2"))
func Scatter2D[E utils.RealNumber](points []Point2D[E], path string, opts ...Option) (string, error) {
	cfg := newConfig(opts...)

	plotted := make([]plottedPoint, len(points))
	for i, p := range points {
		plotted[i] = plottedPoint{
			x:     float64(p.X),
			y:     float64(p.Y),
			scale: 1,
			label: p.Label,
		}
	}

	svg := renderSVG(plotted, cfg, "")
	if err := writeFile(path, svg); err != nil {
		return "", err
	}
	return svg, nil
}
