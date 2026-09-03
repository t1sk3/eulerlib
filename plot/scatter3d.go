package plot

import (
	"bytes"

	"github.com/go-echarts/go-echarts/v2/charts"
	eopts "github.com/go-echarts/go-echarts/v2/opts"

	"github.com/t1sk3/eulerlib/v2/utils"
)

// Scatter3D renders points as an interactive HTML 3D scatter plot — built
// with Apache ECharts' GL extension — and returns the generated HTML
// document. If path is non-empty, the document is also written to that
// file; open it in any browser to drag to orbit the view, scroll to zoom,
// and hover a point for its coordinates and label.
//
// Each point may carry an optional Label, drawn as a caption next to its
// marker; leave it empty to plot a bare point (it still shows its
// coordinates on hover). Use WithTitle, WithXLabel, WithYLabel, and
// WithZLabel to caption the graph itself.
//
//	pts := []plot.Point3D[float64]{
//	    {X: 1, Y: 0, Z: 0, Label: "A"},
//	    {X: 0, Y: 1, Z: 0, Label: "B"},
//	    {X: 0, Y: 0, Z: 1, Label: "C"},
//	}
//	html, err := plot.Scatter3D(pts, "points.html", plot.WithTitle("unit points"))
func Scatter3D[E utils.RealNumber](points []Point3D[E], path string, opts ...Option) (string, error) {
	cfg := newConfig(opts...)

	data := make([]eopts.Chart3DData, len(points))
	for i, p := range points {
		item := eopts.Chart3DData{
			Name:  p.Label,
			Value: []interface{}{float64(p.X), float64(p.Y), float64(p.Z)},
		}
		if p.Label != "" {
			item.Label = &eopts.Label{Show: eopts.Bool(true), Formatter: "{b}", Position: "top"}
		}
		data[i] = item
	}

	chart := charts.NewScatter3D()
	chart.SetGlobalOptions(
		charts.WithInitializationOpts(eopts.Initialization{
			Width:     cfg.width,
			Height:    cfg.height,
			Theme:     cfg.theme,
			PageTitle: pageTitle(cfg.title),
		}),
		charts.WithTitleOpts(eopts.Title{Title: cfg.title}),
		charts.WithXAxis3DOpts(eopts.XAxis3D{Type: "value", Name: cfg.xLabel}),
		charts.WithYAxis3DOpts(eopts.YAxis3D{Type: "value", Name: cfg.yLabel}),
		charts.WithZAxis3DOpts(eopts.ZAxis3D{Type: "value", Name: cfg.zLabel}),
		charts.WithGrid3DOpts(eopts.Grid3D{Show: eopts.Bool(true)}),
		charts.WithTooltipOpts(eopts.Tooltip{Show: eopts.Bool(true), Trigger: "item"}),
	)

	// Note: WithScatterChartOpts is deliberately not used here — it would
	// clobber the "cartesian3D" CoordSystem that AddSeries sets on this
	// series with its own zero-valued CoordSystem field. WithSeriesOpts
	// mutates SymbolSize directly and leaves everything else untouched.
	seriesOpts := []charts.SeriesOpts{
		charts.WithSeriesOpts(func(s *charts.SingleSeries) { s.SymbolSize = cfg.pointSize }),
	}
	if cfg.color != "" {
		seriesOpts = append(seriesOpts, charts.WithItemStyleOpts(eopts.ItemStyle{Color: cfg.color}))
	}
	chart.AddSeries("points", data, seriesOpts...)

	var buf bytes.Buffer
	if err := chart.Render(&buf); err != nil {
		return "", err
	}
	html := buf.String()
	if err := writeFile(path, html); err != nil {
		return "", err
	}
	return html, nil
}
