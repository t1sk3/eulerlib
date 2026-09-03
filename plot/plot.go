// Package plot renders small sets of 2D and 3D data points as standalone
// SVG scatter plots — useful for eyeballing brute-force results, lattice
// points, or geometric constructions without pulling in a full charting
// dependency. It has no external dependencies; everything is plain SVG
// markup built with the standard library.
package plot

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/t1sk3/eulerlib/v2/utils"
)

// Point2D is a single labeled data point to plot in two dimensions. Label
// is optional — leave it empty to draw the point without a caption.
type Point2D[E utils.RealNumber] struct {
	X, Y  E
	Label string
}

// Point3D is a single labeled data point to plot in three dimensions.
// Label is optional — leave it empty to draw the point without a caption.
type Point3D[E utils.RealNumber] struct {
	X, Y, Z E
	Label   string
}

// config holds the rendering options built up by Option funcs.
type config struct {
	title         string
	xLabel        string
	yLabel        string
	zLabel        string
	width, height int
	pointRadius   float64
	showAxes      bool
	pointColor    string
	labelColor    string
}

// Option configures a plot. Pass zero or more to Scatter2D/Scatter3D; any
// left unset fall back to sane defaults.
type Option func(*config)

// WithTitle sets the title drawn above the plot.
func WithTitle(title string) Option {
	return func(c *config) { c.title = title }
}

// WithXLabel sets the caption drawn under the x-axis.
func WithXLabel(label string) Option {
	return func(c *config) { c.xLabel = label }
}

// WithYLabel sets the caption drawn beside the y-axis.
func WithYLabel(label string) Option {
	return func(c *config) { c.yLabel = label }
}

// WithZLabel sets the caption used for the z-axis. It has no effect on
// Scatter2D.
func WithZLabel(label string) Option {
	return func(c *config) { c.zLabel = label }
}

// WithSize sets the rendered SVG canvas size in pixels.
func WithSize(width, height int) Option {
	return func(c *config) {
		c.width, c.height = width, height
	}
}

// WithPointRadius sets the radius, in pixels, of the circle drawn for each
// point (for Scatter3D, points nearer the viewer are drawn larger and
// points further away smaller, relative to this base radius).
func WithPointRadius(r float64) Option {
	return func(c *config) { c.pointRadius = r }
}

// WithAxes toggles the axis frame, tick labels, and axis captions. Axes are
// drawn by default.
func WithAxes(show bool) Option {
	return func(c *config) { c.showAxes = show }
}

// WithColor sets the fill color used for points and their labels. Any
// valid SVG color (name or hex code) is accepted.
func WithColor(color string) Option {
	return func(c *config) {
		c.pointColor, c.labelColor = color, color
	}
}

func newConfig(opts ...Option) config {
	cfg := config{
		width:       640,
		height:      480,
		pointRadius: 4,
		showAxes:    true,
		pointColor:  "#1f77b4",
		labelColor:  "#333333",
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

// plottedPoint is a point already reduced to 2D screen-space coordinates,
// ready to render, along with an optional depth-based size scale (used by
// Scatter3D to give nearer points a larger radius; Scatter2D always uses
// scale 1).
type plottedPoint struct {
	x, y  float64
	scale float64
	label string
}

// renderSVG lays out points inside the plot area and returns the complete
// SVG document as a string. extra, if non-empty, is raw SVG markup
// inserted just before the closing </svg> tag (used by Scatter3D to add an
// orientation gizmo).
func renderSVG(points []plottedPoint, cfg config, extra string) string {
	const padding = 50.0
	const topMargin = 40.0

	plotLeft, plotTop := padding, topMargin
	plotRight, plotBottom := float64(cfg.width)-padding, float64(cfg.height)-padding

	minX, maxX, minY, maxY := boundingBox(points)

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" font-family="sans-serif">`+"\n", cfg.width, cfg.height)
	fmt.Fprintf(&b, `<rect width="%d" height="%d" fill="white"/>`+"\n", cfg.width, cfg.height)

	if cfg.title != "" {
		fmt.Fprintf(&b, `<text x="%.2f" y="24" text-anchor="middle" font-size="18" fill="#111">%s</text>`+"\n",
			float64(cfg.width)/2, escape(cfg.title))
	}

	toScreen := func(x, y float64) (float64, float64) {
		sx := plotLeft + (x-minX)/spanOrOne(maxX-minX)*(plotRight-plotLeft)
		sy := plotBottom - (y-minY)/spanOrOne(maxY-minY)*(plotBottom-plotTop)
		return sx, sy
	}

	if cfg.showAxes {
		fmt.Fprintf(&b, `<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#888" stroke-width="1"/>`+"\n",
			plotLeft, plotBottom, plotRight, plotBottom)
		fmt.Fprintf(&b, `<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#888" stroke-width="1"/>`+"\n",
			plotLeft, plotTop, plotLeft, plotBottom)

		fmt.Fprintf(&b, `<text x="%.2f" y="%.2f" font-size="11" fill="#555">%s</text>`+"\n", plotLeft, plotBottom+16, formatFloat(minX))
		fmt.Fprintf(&b, `<text x="%.2f" y="%.2f" font-size="11" fill="#555" text-anchor="end">%s</text>`+"\n", plotRight, plotBottom+16, formatFloat(maxX))
		fmt.Fprintf(&b, `<text x="%.2f" y="%.2f" font-size="11" fill="#555" text-anchor="end">%s</text>`+"\n", plotLeft-6, plotBottom, formatFloat(minY))
		fmt.Fprintf(&b, `<text x="%.2f" y="%.2f" font-size="11" fill="#555" text-anchor="end">%s</text>`+"\n", plotLeft-6, plotTop+10, formatFloat(maxY))

		if cfg.xLabel != "" {
			fmt.Fprintf(&b, `<text x="%.2f" y="%.2f" text-anchor="middle" font-size="13" fill="#333">%s</text>`+"\n",
				(plotLeft+plotRight)/2, float64(cfg.height)-12, escape(cfg.xLabel))
		}
		if cfg.yLabel != "" {
			fmt.Fprintf(&b, `<text x="14" y="%.2f" text-anchor="middle" font-size="13" fill="#333" transform="rotate(-90 14 %.2f)">%s</text>`+"\n",
				(plotTop+plotBottom)/2, (plotTop+plotBottom)/2, escape(cfg.yLabel))
		}
	}

	for _, p := range points {
		sx, sy := toScreen(p.x, p.y)
		r := cfg.pointRadius * p.scale
		fmt.Fprintf(&b, `<circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s"/>`+"\n", sx, sy, r, cfg.pointColor)
		if p.label != "" {
			fmt.Fprintf(&b, `<text x="%.2f" y="%.2f" font-size="11" fill="%s">%s</text>`+"\n",
				sx+r+3, sy-r-3, cfg.labelColor, escape(p.label))
		}
	}

	if extra != "" {
		b.WriteString(extra)
		b.WriteString("\n")
	}

	b.WriteString("</svg>\n")
	return b.String()
}

func boundingBox(points []plottedPoint) (minX, maxX, minY, maxY float64) {
	if len(points) == 0 {
		return 0, 1, 0, 1
	}
	minX, maxX = points[0].x, points[0].x
	minY, maxY = points[0].y, points[0].y
	for _, p := range points[1:] {
		minX, maxX = math.Min(minX, p.x), math.Max(maxX, p.x)
		minY, maxY = math.Min(minY, p.y), math.Max(maxY, p.y)
	}
	return
}

// spanOrOne avoids a division by zero when every point shares the same
// x (or y) coordinate, in which case they're simply centered.
func spanOrOne(span float64) float64 {
	if span == 0 {
		return 1
	}
	return span
}

func formatFloat(f float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.4f", f), "0"), ".")
}

var svgEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")

func escape(s string) string {
	return svgEscaper.Replace(s)
}

// writeFile writes svg to path if path is non-empty; the empty path is a
// no-op so callers can request just the returned string.
func writeFile(path, svg string) error {
	if path == "" {
		return nil
	}
	return os.WriteFile(path, []byte(svg), 0644)
}
