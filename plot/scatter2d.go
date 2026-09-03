package plot

import (
	"bytes"

	"github.com/go-echarts/go-echarts/v2/charts"
	eopts "github.com/go-echarts/go-echarts/v2/opts"

	"github.com/t1sk3/eulerlib/v2/utils"
)

// Scatter2D renders points as an interactive HTML scatter plot — built
// with Apache ECharts — and returns the generated HTML document. If path
// is non-empty, the document is also written to that file; open it in any
// browser to hover a point for its coordinates and label, drag the bottom
// slider (or scroll) to zoom, and pan around.
//
// Each point may carry an optional Label, drawn as a caption next to its
// marker; leave it empty to plot a bare point (it still shows its
// coordinates on hover). Use WithTitle, WithXLabel, and WithYLabel to
// caption the graph itself.
//
//	pts := []plot.Point2D[int]{
//	    {X: 1, Y: 1, Label: "A"},
//	    {X: 2, Y: 4},
//	    {X: 3, Y: 9, Label: "C"},
//	}
//	html, err := plot.Scatter2D(pts, "points.html", plot.WithTitle("y = x^2"))
func Scatter2D[E utils.RealNumber](points []Point2D[E], path string, opts ...Option) (string, error) {
	cfg := newConfig(opts...)

	data := make([]eopts.ScatterData, len(points))
	for i, p := range points {
		data[i] = eopts.ScatterData{
			Name:  p.Label,
			Value: []interface{}{float64(p.X), float64(p.Y)},
		}
	}

	chart := charts.NewScatter()
	chart.SetGlobalOptions(
		charts.WithInitializationOpts(eopts.Initialization{
			Width:     cfg.width,
			Height:    cfg.height,
			Theme:     cfg.theme,
			PageTitle: pageTitle(cfg.title),
		}),
		charts.WithTitleOpts(eopts.Title{Title: cfg.title}),
		charts.WithXAxisOpts(eopts.XAxis{Type: "value", Name: cfg.xLabel, Scale: eopts.Bool(true)}),
		charts.WithYAxisOpts(eopts.YAxis{Type: "value", Name: cfg.yLabel, Scale: eopts.Bool(true)}),
		charts.WithTooltipOpts(eopts.Tooltip{Show: eopts.Bool(true), Trigger: "item"}),
		charts.WithDataZoomOpts(
			eopts.DataZoom{Type: "inside"},
			eopts.DataZoom{Type: "slider"},
		),
		charts.WithToolboxOpts(eopts.Toolbox{
			Show: eopts.Bool(true),
			Feature: &eopts.ToolBoxFeature{
				SaveAsImage: &eopts.ToolBoxFeatureSaveAsImage{Show: eopts.Bool(true)},
				DataZoom:    &eopts.ToolBoxFeatureDataZoom{Show: eopts.Bool(true)},
				Restore:     &eopts.ToolBoxFeatureRestore{Show: eopts.Bool(true)},
			},
		}),
	)

	seriesOpts := []charts.SeriesOpts{
		// The {b} token renders each point's Name (our Label); it's simply
		// blank for unlabeled points, so this is safe to leave on always.
		charts.WithLabelOpts(eopts.Label{Show: eopts.Bool(true), Formatter: "{b}", Position: "top"}),
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
