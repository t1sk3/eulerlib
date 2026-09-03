package plot

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveCSV2D(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "points.csv")

	points := []Point2D[int]{
		{X: 1, Y: 2, Label: "A"},
		{X: -3, Y: 4},
	}
	if err := SaveCSV2D(points, path); err != nil {
		t.Fatalf("SaveCSV2D returned error: %v", err)
	}

	rows := readCSV(t, path)
	want := [][]string{
		{"x", "y", "label"},
		{"1", "2", "A"},
		{"-3", "4", ""},
	}
	assertRowsEqual(t, want, rows)
}

func TestSaveCSV3D(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "points.csv")

	points := []Point3D[float64]{
		{X: 1.5, Y: 2.5, Z: 3.5, Label: "A"},
		{X: 0, Y: 0, Z: 0},
	}
	if err := SaveCSV3D(points, path); err != nil {
		t.Fatalf("SaveCSV3D returned error: %v", err)
	}

	rows := readCSV(t, path)
	want := [][]string{
		{"x", "y", "z", "label"},
		{"1.5", "2.5", "3.5", "A"},
		{"0", "0", "0", ""},
	}
	assertRowsEqual(t, want, rows)
}

func TestSaveCSV2DLargeIntPrecision(t *testing.T) {
	// A float64 round-trip would lose precision here; formatNumber must
	// not go through float64 for integer types.
	dir := t.TempDir()
	path := filepath.Join(dir, "points.csv")

	const big int64 = 9007199254740993 // 2^53 + 1, not exactly representable as float64
	points := []Point2D[int64]{{X: big, Y: 1}}
	if err := SaveCSV2D(points, path); err != nil {
		t.Fatalf("SaveCSV2D returned error: %v", err)
	}

	rows := readCSV(t, path)
	if rows[1][0] != "9007199254740993" {
		t.Fatalf("got x = %q, want exact int64 value 9007199254740993", rows[1][0])
	}
}

func TestSaveCSV2DEmptyPoints(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "points.csv")

	if err := SaveCSV2D([]Point2D[int]{}, path); err != nil {
		t.Fatalf("SaveCSV2D returned error: %v", err)
	}
	rows := readCSV(t, path)
	assertRowsEqual(t, [][]string{{"x", "y", "label"}}, rows)
}

func TestSaveCSV2DInvalidPath(t *testing.T) {
	err := SaveCSV2D([]Point2D[int]{{X: 1, Y: 1}}, filepath.Join(t.TempDir(), "missing-dir", "points.csv"))
	if err == nil {
		t.Fatalf("expected an error writing to a nonexistent directory")
	}
}

func TestSaveCSV3DInvalidPath(t *testing.T) {
	err := SaveCSV3D([]Point3D[int]{{X: 1, Y: 1, Z: 1}}, filepath.Join(t.TempDir(), "missing-dir", "points.csv"))
	if err == nil {
		t.Fatalf("expected an error writing to a nonexistent directory")
	}
}

func readCSV(t *testing.T, path string) [][]string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open %s: %v", path, err)
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV: %v", err)
	}
	return rows
}

func assertRowsEqual(t *testing.T, want, got [][]string) {
	t.Helper()
	if len(want) != len(got) {
		t.Fatalf("got %d rows, want %d (got=%v)", len(got), len(want), got)
	}
	for i := range want {
		if len(want[i]) != len(got[i]) {
			t.Fatalf("row %d: got %v, want %v", i, got[i], want[i])
		}
		for j := range want[i] {
			if want[i][j] != got[i][j] {
				t.Fatalf("row %d col %d: got %q, want %q", i, j, got[i][j], want[i][j])
			}
		}
	}
}
