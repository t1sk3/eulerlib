// Package plot renders small sets of 2D and 3D data points as interactive,
// browser-based scatter plots. Charts are built with Apache ECharts (via
// go-echarts/go-echarts, and its echarts-gl extension for 3D) and written
// out as a single self-contained HTML file — open it in a browser to hover
// points for a tooltip, scroll/drag to zoom and pan, and, for Scatter3D,
// drag to orbit the 3D view.
package plot

import (
	"fmt"
	"os"

	"github.com/t1sk3/eulerlib/v2/utils"
)

// Point2D is a single labeled data point to plot in two dimensions. Label
// is optional — leave it empty to plot a bare point (it still shows its
// coordinates in the hover tooltip).
type Point2D[E utils.RealNumber] struct {
	X, Y  E
	Label string
}

// Point3D is a single labeled data point to plot in three dimensions.
// Label is optional — leave it empty to plot a bare point (it still shows
// its coordinates in the hover tooltip).
type Point3D[E utils.RealNumber] struct {
	X, Y, Z E
	Label   string
}

// config holds the rendering options built up by Option funcs.
type config struct {
	title                  string
	xLabel, yLabel, zLabel string
	width, height          string
	pointSize              float32
	color                  string
	theme                  string
}

// Option configures a plot. Pass zero or more to Scatter2D/Scatter3D; any
// left unset fall back to sane defaults.
type Option func(*config)

// WithTitle sets the title drawn above the chart.
func WithTitle(title string) Option {
	return func(c *config) { c.title = title }
}

// WithXLabel sets the x-axis caption.
func WithXLabel(label string) Option {
	return func(c *config) { c.xLabel = label }
}

// WithYLabel sets the y-axis caption.
func WithYLabel(label string) Option {
	return func(c *config) { c.yLabel = label }
}

// WithZLabel sets the z-axis caption. It has no effect on Scatter2D.
func WithZLabel(label string) Option {
	return func(c *config) { c.zLabel = label }
}

// WithSize sets the rendered chart's canvas size in pixels.
func WithSize(width, height int) Option {
	return func(c *config) {
		c.width = fmt.Sprintf("%dpx", width)
		c.height = fmt.Sprintf("%dpx", height)
	}
}

// WithPointSize sets the diameter, in pixels, of the marker drawn for each
// point.
func WithPointSize(size float32) Option {
	return func(c *config) { c.pointSize = size }
}

// WithColor sets the marker color for the whole series. Any valid CSS
// color (name or hex code) is accepted; leave unset to use ECharts'
// default palette color.
func WithColor(color string) Option {
	return func(c *config) { c.color = color }
}

// WithTheme selects one of ECharts' built-in themes, e.g. "dark",
// "vintage", "macarons", "infographic", "shine", "roma", or "westeros".
// Left unset, the default light theme is used.
func WithTheme(theme string) Option {
	return func(c *config) { c.theme = theme }
}

func newConfig(opts ...Option) config {
	cfg := config{
		width:     "900px",
		height:    "600px",
		pointSize: 10,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

// pageTitle returns the browser tab title: the chart title if set, or a
// generic fallback.
func pageTitle(title string) string {
	if title != "" {
		return title
	}
	return "eulerlib plot"
}

// writeFile writes content to path if path is non-empty; the empty path is
// a no-op so callers can request just the returned HTML string.
func writeFile(path, content string) error {
	if path == "" {
		return nil
	}
	return os.WriteFile(path, []byte(content), 0644)
}
