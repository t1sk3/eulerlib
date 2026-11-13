package paths

import (
	"math/rand/v2"
	"testing"
)

// The 5×5 matrix of Project Euler #81, #82 and #83.
var euler = [][]int{
	{131, 673, 234, 103, 18},
	{201, 96, 342, 965, 150},
	{630, 803, 746, 422, 111},
	{537, 699, 497, 121, 956},
	{805, 732, 524, 37, 331},
}

func TestTriangleMax(t *testing.T) {
	if got := TriangleMax([][]int{{3}, {7, 4}, {2, 4, 6}, {8, 5, 9, 3}}); got != 23 {
		t.Errorf("TriangleMax = %d, want 23", got)
	}
	// Euler #18.
	tri := [][]int{
		{75},
		{95, 64},
		{17, 47, 82},
		{18, 35, 87, 10},
		{20, 4, 82, 47, 65},
		{19, 1, 23, 75, 3, 34},
		{88, 2, 77, 73, 7, 63, 67},
		{99, 65, 4, 28, 6, 16, 70, 92},
		{41, 41, 26, 56, 83, 40, 80, 70, 33},
		{41, 48, 72, 33, 47, 32, 37, 16, 94, 29},
		{53, 71, 44, 65, 25, 43, 91, 52, 97, 51, 14},
		{70, 11, 33, 28, 77, 73, 17, 78, 39, 68, 17, 57},
		{91, 71, 52, 38, 17, 14, 91, 43, 58, 50, 27, 29, 48},
		{63, 66, 4, 68, 89, 53, 67, 30, 73, 16, 69, 87, 40, 31},
		{4, 62, 98, 27, 23, 9, 70, 98, 73, 93, 38, 53, 60, 4, 23},
	}
	if got := TriangleMax(tri); got != 1074 {
		t.Errorf("Euler #18 = %d, want 1074", got)
	}
	if tri[0][0] != 75 || tri[14][0] != 4 {
		t.Error("TriangleMax modified its input")
	}
	// Against trying every path.
	for range 200 {
		n := 1 + rand.IntN(9)
		tri := make([][]float64, n)
		for i := range tri {
			tri[i] = make([]float64, i+1)
			for j := range tri[i] {
				tri[i][j] = float64(rand.IntN(41) - 20)
			}
		}
		var best func(i, j int) float64
		best = func(i, j int) float64 {
			if i == n-1 {
				return tri[i][j]
			}
			return tri[i][j] + max(best(i+1, j), best(i+1, j+1))
		}
		if got := TriangleMax(tri); got != best(0, 0) {
			t.Fatalf("TriangleMax(%v) = %v, want %v", tri, got, best(0, 0))
		}
	}
	if TriangleMax[int](nil) != 0 || TriangleMax([][]int{{-5}}) != -5 {
		t.Error("TriangleMax edge cases")
	}
	mustPanic(t, "ragged triangle", func() { TriangleMax([][]int{{1}, {2}}) })
}

func TestGridEuler(t *testing.T) {
	if got := GridMinTwoWays(euler); got != 2427 {
		t.Errorf("Euler #81 example = %d, want 2427", got)
	}
	if got := GridMinThreeWays(euler); got != 994 {
		t.Errorf("Euler #82 example = %d, want 994", got)
	}
	if got := GridMinFourWays(euler); got != 2297 {
		t.Errorf("Euler #83 example = %d, want 2297", got)
	}
	if euler[0][0] != 131 || euler[4][4] != 331 {
		t.Error("a grid function modified its input")
	}
}

// bellmanFord returns the cheapest path sum from any of the start cells to
// each cell, relaxing along the allowed moves until nothing changes.
func bellmanFord(g [][]int, moves [][2]int, starts [][2]int) [][]int {
	const inf = 1 << 60
	d := make([][]int, len(g))
	for i := range d {
		d[i] = make([]int, len(g[i]))
		for j := range d[i] {
			d[i][j] = inf
		}
	}
	for _, s := range starts {
		d[s[0]][s[1]] = g[s[0]][s[1]]
	}
	for changed := true; changed; {
		changed = false
		for i := range g {
			for j := range g[i] {
				for _, m := range moves {
					ni, nj := i+m[0], j+m[1]
					if ni < 0 || ni >= len(g) || nj < 0 || nj >= len(g[i]) {
						continue
					}
					if v := d[i][j] + g[ni][nj]; d[i][j] < inf && v < d[ni][nj] {
						d[ni][nj], changed = v, true
					}
				}
			}
		}
	}
	return d
}

func TestGridAgainstBellmanFord(t *testing.T) {
	for range 300 {
		rows, cols := 1+rand.IntN(7), 1+rand.IntN(7)
		g := make([][]int, rows)
		for i := range g {
			g[i] = make([]int, cols)
			for j := range g[i] {
				g[i][j] = rand.IntN(30)
			}
		}
		right, down, left, up := [2]int{0, 1}, [2]int{1, 0}, [2]int{0, -1}, [2]int{-1, 0}

		want := bellmanFord(g, [][2]int{right, down}, [][2]int{{0, 0}})[rows-1][cols-1]
		if got := GridMinTwoWays(g); got != want {
			t.Fatalf("GridMinTwoWays(%v) = %d, want %d", g, got, want)
		}

		starts := make([][2]int, rows)
		for i := range starts {
			starts[i] = [2]int{i, 0}
		}
		d := bellmanFord(g, [][2]int{right, down, up}, starts)
		want = d[0][cols-1]
		for i := range d {
			want = min(want, d[i][cols-1])
		}
		if got := GridMinThreeWays(g); got != want {
			t.Fatalf("GridMinThreeWays(%v) = %d, want %d", g, got, want)
		}

		want = bellmanFord(g, [][2]int{right, down, left, up}, [][2]int{{0, 0}})[rows-1][cols-1]
		if got := GridMinFourWays(g); got != want {
			t.Fatalf("GridMinFourWays(%v) = %d, want %d", g, got, want)
		}
	}
}

func TestGridEdgeCases(t *testing.T) {
	// A wall of expensive cells that only a path doubling back avoids.
	snake := [][]int{
		{1, 1, 1, 1, 1},
		{9, 9, 9, 9, 1},
		{1, 1, 1, 1, 1},
		{1, 9, 9, 9, 9},
		{1, 1, 1, 1, 1},
	}
	if got := GridMinFourWays(snake); got != 17 {
		t.Errorf("GridMinFourWays(snake) = %d, want 17", got)
	}
	// Negative values are fine when moves only go right and down.
	if got := GridMinTwoWays([][]int{{1, -5, 2}, {4, -9, -1}}); got != -14 {
		t.Errorf("GridMinTwoWays with negatives = %d, want -14", got)
	}
	if got := GridMinFourWays([][]float64{{0.5, 0.25}, {8, 0.125}}); got != 0.875 {
		t.Errorf("GridMinFourWays[float64] = %v", got)
	}
	if got := GridMinFourWays([][]uint8{{0, 0}, {0, 0}}); got != 0 {
		t.Errorf("GridMinFourWays of zeros = %d", got)
	}
	one := [][]int{{7}}
	if GridMinTwoWays(one) != 7 || GridMinThreeWays(one) != 7 || GridMinFourWays(one) != 7 {
		t.Error("1×1 grid")
	}
	for _, g := range [][][]int{nil, {}, {{}}, {{}, {}}} {
		if GridMinTwoWays(g) != 0 || GridMinThreeWays(g) != 0 || GridMinFourWays(g) != 0 {
			t.Errorf("empty grid %v should give 0", g)
		}
	}
	ragged := [][]int{{1, 2}, {3}}
	mustPanic(t, "GridMinTwoWays ragged", func() { GridMinTwoWays(ragged) })
	mustPanic(t, "GridMinThreeWays ragged", func() { GridMinThreeWays(ragged) })
	mustPanic(t, "GridMinFourWays ragged", func() { GridMinFourWays(ragged) })
}

func mustPanic(t *testing.T, name string, f func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s: expected a panic", name)
		}
	}()
	f()
}
