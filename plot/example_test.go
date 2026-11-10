package plot_test

import (
	"fmt"
	"strings"

	"github.com/t1sk3/eulerlib/v4/plot"
	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/vector"
)

func ExampleScatter2D() {
	pts := plot.Points2D(stream.RangeClosed(1, 5).
		Map(func(n int) vector.Vec2[int] { return vector.Of2(n, n*n) }).All())
	html, err := plot.Scatter2D(pts, "", plot.WithTitle("y = x^2"))
	fmt.Println(err, strings.Contains(html, "y = x^2"))
	// Output: <nil> true
}
