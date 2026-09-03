package plot

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/t1sk3/eulerlib/v2/utils"
)

// SaveCSV2D writes points to path as CSV — one row per point, columns
// "x,y,label" — for batch-processing elsewhere (a spreadsheet, pandas,
// another program). The label column is empty for points with no Label.
func SaveCSV2D[E utils.RealNumber](points []Point2D[E], path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.Write([]string{"x", "y", "label"}); err != nil {
		return err
	}
	for _, p := range points {
		if err := w.Write([]string{formatNumber(p.X), formatNumber(p.Y), p.Label}); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

// SaveCSV3D writes points to path as CSV — one row per point, columns
// "x,y,z,label" — for batch-processing elsewhere (a spreadsheet, pandas,
// another program). The label column is empty for points with no Label.
func SaveCSV3D[E utils.RealNumber](points []Point3D[E], path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.Write([]string{"x", "y", "z", "label"}); err != nil {
		return err
	}
	for _, p := range points {
		if err := w.Write([]string{formatNumber(p.X), formatNumber(p.Y), formatNumber(p.Z), p.Label}); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

// formatNumber renders n with Go's default numeric formatting. Unlike a
// float64 conversion, this keeps large integer types (e.g. int64) exact
// instead of losing precision above 2^53.
func formatNumber[E utils.RealNumber](n E) string {
	return fmt.Sprintf("%v", n)
}
