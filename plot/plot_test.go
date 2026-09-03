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
	svg, err := Scatter2D(pts, "", WithTitle("y = x^2"), WithXLabel("x"), WithYLabel("y"))
	if err != nil {
		t.Fatalf("Scatter2D returned error: %v", err)
	}

	if !strings.HasPrefix(svg, "<svg") {
		t.Fatalf("output does not look like SVG markup: %q", svg[:min(20, len(svg))])
	}
	if !strings.HasSuffix(strings.TrimSpace(svg), "</svg>") {
		t.Fatalf("output is not a closed SVG document")
	}
	if strings.Count(svg, "<circle") != len(pts) {
		t.Fatalf("expected %d points drawn, got %d circles", len(pts), strings.Count(svg, "<circle"))
	}
	if !strings.Contains(svg, ">A<") || !strings.Contains(svg, ">C<") {
		t.Fatalf("expected labels A and C to appear in the SVG")
	}
	if !strings.Contains(svg, "y = x^2") {
		t.Fatalf("expected title text in output")
	}
}

func TestScatter2DEmptyLabelNotDrawn(t *testing.T) {
	pts := []Point2D[int]{{X: 0, Y: 0}}
	svg, err := Scatter2D(pts, "")
	if err != nil {
		t.Fatalf("Scatter2D returned error: %v", err)
	}
	if strings.Contains(svg, "></text>") {
		t.Fatalf("unlabeled point should not draw an empty label")
	}
}

func TestScatter2DWritesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.svg")
	pts := []Point2D[float64]{{X: 0, Y: 0}, {X: 1, Y: 1}}

	svg, err := Scatter2D(pts, path)
	if err != nil {
		t.Fatalf("Scatter2D returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected file to be written: %v", err)
	}
	if string(data) != svg {
		t.Fatalf("file content does not match returned SVG string")
	}
}

func TestScatter2DConstantCoordinates(t *testing.T) {
	// All points share the same x and y — must not divide by zero.
	pts := []Point2D[int]{{X: 5, Y: 5}, {X: 5, Y: 5}, {X: 5, Y: 5}}
	svg, err := Scatter2D(pts, "")
	if err != nil {
		t.Fatalf("Scatter2D returned error: %v", err)
	}
	if strings.Contains(svg, "NaN") || strings.Contains(svg, "Inf") {
		t.Fatalf("degenerate bounding box produced NaN/Inf in output: %s", svg)
	}
}

func TestScatter2DEmptyPoints(t *testing.T) {
	svg, err := Scatter2D([]Point2D[int]{}, "")
	if err != nil {
		t.Fatalf("Scatter2D returned error: %v", err)
	}
	if strings.Contains(svg, "<circle") {
		t.Fatalf("expected no circles for an empty point set")
	}
}

func TestScatter3DBasic(t *testing.T) {
	pts := []Point3D[float64]{
		{X: 1, Y: 0, Z: 0, Label: "A"},
		{X: 0, Y: 1, Z: 0, Label: "B"},
		{X: 0, Y: 0, Z: 1, Label: "C"},
	}
	svg, err := Scatter3D(pts, "", WithTitle("unit points"))
	if err != nil {
		t.Fatalf("Scatter3D returned error: %v", err)
	}
	if strings.Count(svg, "<circle") != len(pts) {
		t.Fatalf("expected %d points drawn, got %d circles", len(pts), strings.Count(svg, "<circle"))
	}
	for _, label := range []string{"A", "B", "C"} {
		if !strings.Contains(svg, ">"+label+"<") {
			t.Fatalf("expected label %q in output", label)
		}
	}
	// The orientation gizmo should be present by default.
	if !strings.Contains(svg, ">X<") || !strings.Contains(svg, ">Y<") || !strings.Contains(svg, ">Z<") {
		t.Fatalf("expected default XYZ gizmo labels in output")
	}
}

func TestScatter3DAxesDisabled(t *testing.T) {
	pts := []Point3D[int]{{X: 1, Y: 2, Z: 3}}
	svg, err := Scatter3D(pts, "", WithAxes(false))
	if err != nil {
		t.Fatalf("Scatter3D returned error: %v", err)
	}
	if strings.Contains(svg, ">X<") || strings.Contains(svg, ">Y<") || strings.Contains(svg, ">Z<") {
		t.Fatalf("expected no gizmo when WithAxes(false)")
	}
}

func TestScatter3DCustomAxisLabels(t *testing.T) {
	pts := []Point3D[int]{{X: 1, Y: 2, Z: 3}}
	svg, err := Scatter3D(pts, "", WithXLabel("lat"), WithYLabel("lon"), WithZLabel("alt"))
	if err != nil {
		t.Fatalf("Scatter3D returned error: %v", err)
	}
	for _, label := range []string{"lat", "lon", "alt"} {
		if !strings.Contains(svg, ">"+label+"<") {
			t.Fatalf("expected custom gizmo label %q in output", label)
		}
	}
}

func TestScatter3DDepthOrdering(t *testing.T) {
	// Nearer points (larger x+y+z) should be drawn after farther ones so
	// they layer on top; verify via the order circles appear in the SVG.
	pts := []Point3D[int]{
		{X: 10, Y: 10, Z: 10, Label: "near"},
		{X: -10, Y: -10, Z: -10, Label: "far"},
	}
	svg, err := Scatter3D(pts, "")
	if err != nil {
		t.Fatalf("Scatter3D returned error: %v", err)
	}
	farIdx := strings.Index(svg, ">far<")
	nearIdx := strings.Index(svg, ">near<")
	if farIdx == -1 || nearIdx == -1 {
		t.Fatalf("expected both labels present")
	}
	if farIdx > nearIdx {
		t.Fatalf("expected farther point drawn before nearer point (painter's algorithm)")
	}
}

func TestScatter3DWritesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out3d.svg")
	pts := []Point3D[int]{{X: 1, Y: 2, Z: 3}}

	if _, err := Scatter3D(pts, path); err != nil {
		t.Fatalf("Scatter3D returned error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to be written: %v", err)
	}
}

func TestIsometricProjectOrigin(t *testing.T) {
	sx, sy := isometricProject(0, 0, 0)
	if sx != 0 || sy != 0 {
		t.Fatalf("isometricProject(0,0,0) = (%v, %v), want (0, 0)", sx, sy)
	}
}

func TestDepthScaleConstantDepth(t *testing.T) {
	// When every point has the same depth, scale must not divide by zero.
	if got := depthScale(5, 5, 5); got != 1 {
		t.Fatalf("depthScale with equal min/max = %v, want 1", got)
	}
}
