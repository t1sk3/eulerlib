package plot

import (
	"strings"
	"testing"
)

// Caller-supplied strings end up inside a <script> element, so they must
// not be able to close it or open a new one.
func TestScatterEscapesMarkup(t *testing.T) {
	const evil = `</script><script>alert(1)</script>&amp;`
	opts := []Option{WithTitle(evil), WithXLabel(evil), WithYLabel(evil), WithZLabel(evil), WithColor(evil), WithTheme(evil)}

	html2, err := Scatter2D([]Point2D[int]{{X: 1, Y: 2, Label: evil}}, "", opts...)
	if err != nil {
		t.Fatal(err)
	}
	html3, err := Scatter3D([]Point3D[int]{{X: 1, Y: 2, Z: 3, Label: evil}}, "", opts...)
	if err != nil {
		t.Fatal(err)
	}
	for name, html := range map[string]string{"Scatter2D": html2, "Scatter3D": html3} {
		for _, bad := range []string{"<script>alert", "alert(1)</script>", "﷐", "﷑", "﷒"} {
			if strings.Contains(html, bad) {
				t.Errorf("%s output contains %q", name, bad)
			}
		}
		// The label is still there, as a JSON string that decodes to the
		// original text.
		// (The escapes are spelled with a separate backslash so that no tool
		// that touches this file can decode them.)
		const bs = "\\"
		lt, gt, amp := bs+"u003c", bs+"u003e", bs+"u0026"
		want := `"name":"` + lt + "/script" + gt + lt + "script" + gt + "alert(1)" + lt + "/script" + gt + amp + `amp;"`
		if !strings.Contains(html, want) {
			t.Errorf("%s output is missing the escaped label %s", name, want)
		}
	}

	// Plain labels are untouched.
	html, err := Scatter2D([]Point2D[int]{{X: 1, Y: 2, Label: "a b"}}, "", WithTitle("y = x^2"))
	if err != nil || !strings.Contains(html, `"name":"a b"`) || !strings.Contains(html, `"text":"y = x^2"`) {
		t.Errorf("plain strings changed: %v", err)
	}
}

func TestSaveCSVReportsErrors(t *testing.T) {
	if err := SaveCSV2D([]Point2D[int]{{X: 1, Y: 2}}, t.TempDir()); err == nil {
		t.Error("SaveCSV2D to a directory should fail")
	}
	if err := SaveCSV3D([]Point3D[int]{{X: 1, Y: 2, Z: 3}}, t.TempDir()); err == nil {
		t.Error("SaveCSV3D to a directory should fail")
	}
}
