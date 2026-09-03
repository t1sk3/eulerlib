package plot

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScatter2DBasic(t *testing.T) {
	pts := []Point2D[int]{
		{X: 1, Y: 1, Label: "A"},
		{X: 2, Y: 4},
		{X: 3, Y: 9, Label: "C"},
	}
	html, err := Scatter2D(pts, "", WithTitle("y = x^2"), WithXLabel("x"), WithYLabel("y"))
	if err != nil {
		t.Fatalf("Scatter2D returned error: %v", err)
	}

	if !strings.Contains(html, "<html") {
		t.Fatalf("output does not look like a full HTML document")
	}
	if !strings.Contains(html, "echarts.min.js") {
		t.Fatalf("expected the ECharts script tag in output")
	}
	if !strings.Contains(html, "y = x^2") {
		t.Fatalf("expected title text in output")
	}
	for _, want := range []string{`"A"`, `"C"`} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected label %s to appear in the series data", want)
		}
	}
	if strings.Count(html, "1,1") == 0 && strings.Count(html, "1, 1") == 0 {
		// coordinates are embedded as a JSON array; just make sure some
		// point values made it into the document.
		if !strings.Contains(html, "9") {
			t.Fatalf("expected point values to appear in output")
		}
	}
}

func TestScatter2DWritesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.html")
	pts := []Point2D[float64]{{X: 0, Y: 0}, {X: 1, Y: 1}}

	html, err := Scatter2D(pts, path)
	if err != nil {
		t.Fatalf("Scatter2D returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected file to be written: %v", err)
	}
	if string(data) != html {
		t.Fatalf("file content does not match returned HTML string")
	}
}

func TestScatter2DEmptyPath(t *testing.T) {
	pts := []Point2D[int]{{X: 0, Y: 0}}
	if _, err := Scatter2D(pts, ""); err != nil {
		t.Fatalf("Scatter2D returned error: %v", err)
	}
}

func TestScatter2DEmptyPoints(t *testing.T) {
	html, err := Scatter2D([]Point2D[int]{}, "")
	if err != nil {
		t.Fatalf("Scatter2D returned error: %v", err)
	}
	if !strings.Contains(html, "<html") {
		t.Fatalf("expected a valid HTML document even with no points")
	}
}

func TestScatter2DCustomOptions(t *testing.T) {
	pts := []Point2D[int]{{X: 1, Y: 2}}
	html, err := Scatter2D(pts, "",
		WithSize(400, 300),
		WithPointSize(20),
		WithColor("#ff0000"),
		WithTheme("vintage"),
	)
	if err != nil {
		t.Fatalf("Scatter2D returned error: %v", err)
	}
	for _, want := range []string{"400px", "300px", "#ff0000", "themes/vintage.js"} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected %q in output", want)
		}
	}
}

func TestScatter3DBasic(t *testing.T) {
	pts := []Point3D[float64]{
		{X: 1, Y: 0, Z: 0, Label: "A"},
		{X: 0, Y: 1, Z: 0, Label: "B"},
		{X: 0, Y: 0, Z: 1, Label: "C"},
	}
	html, err := Scatter3D(pts, "", WithTitle("unit points"), WithXLabel("lat"), WithYLabel("lon"), WithZLabel("alt"))
	if err != nil {
		t.Fatalf("Scatter3D returned error: %v", err)
	}
	if !strings.Contains(html, "<html") {
		t.Fatalf("output does not look like a full HTML document")
	}
	if !strings.Contains(html, "echarts-gl.min.js") {
		t.Fatalf("expected the echarts-gl script tag for a 3D chart")
	}
	for _, want := range []string{`"A"`, `"B"`, `"C"`, "lat", "lon", "alt"} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected %q in output", want)
		}
	}
}

func TestScatter3DWritesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out3d.html")
	pts := []Point3D[int]{{X: 1, Y: 2, Z: 3}}

	if _, err := Scatter3D(pts, path); err != nil {
		t.Fatalf("Scatter3D returned error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to be written: %v", err)
	}
}

func TestScatter3DEmptyPoints(t *testing.T) {
	html, err := Scatter3D([]Point3D[int]{}, "")
	if err != nil {
		t.Fatalf("Scatter3D returned error: %v", err)
	}
	if !strings.Contains(html, "<html") {
		t.Fatalf("expected a valid HTML document even with no points")
	}
}

func TestScatter3DCustomOptions(t *testing.T) {
	pts := []Point3D[int]{{X: 1, Y: 2, Z: 3}}
	html, err := Scatter3D(pts, "", WithSize(500, 400), WithPointSize(15), WithColor("#00ff00"))
	if err != nil {
		t.Fatalf("Scatter3D returned error: %v", err)
	}
	for _, want := range []string{"500px", "400px", "#00ff00"} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected %q in output", want)
		}
	}
}
