// Package paths finds the best path sum through a triangle or a grid of
// numbers by dynamic programming (and Dijkstra's algorithm where moves can
// double back): the family of Project Euler #18, #67, #81, #82 and #83.
//
// A grid is a slice of rows of equal length. Every function leaves its
// input untouched, returns 0 for an empty triangle or grid, and adds the
// values with ordinary arithmetic, so a sum that does not fit in E
// overflows silently.
package paths

import (
	"container/heap"

	"github.com/t1sk3/go-utils/types/constraints"
)

// TriangleMax returns the largest sum of a path from the top of the
// triangle to its bottom row, where each step goes to one of the two
// numbers below: from tri[i][j] to tri[i+1][j] or tri[i+1][j+1] (Project
// Euler #18 and #67). It panics unless row i has exactly i+1 numbers.
func TriangleMax[E constraints.Number](tri [][]E) E {
	for i, row := range tri {
		if len(row) != i+1 {
			panic("paths: TriangleMax needs row i to have i+1 numbers")
		}
	}
	if len(tri) == 0 {
		return 0
	}
	// Bottom-up: best[j] is the best sum from row i, position j, down.
	best := append([]E(nil), tri[len(tri)-1]...)
	for i := len(tri) - 2; i >= 0; i-- {
		for j := 0; j <= i; j++ {
			best[j] = tri[i][j] + max(best[j], best[j+1])
		}
	}
	return best[0]
}

// checkGrid panics unless every row of g has the same length, and returns
// the number of rows and columns (0, 0 if there are no cells).
func checkGrid[E any](g [][]E) (rows, cols int) {
	if len(g) == 0 {
		return 0, 0
	}
	cols = len(g[0])
	for _, row := range g {
		if len(row) != cols {
			panic("paths: grid rows have different lengths")
		}
	}
	if cols == 0 {
		return 0, 0
	}
	return len(g), cols
}

// GridMinTwoWays returns the smallest sum of a path from the top-left cell
// to the bottom-right cell that only moves right and down (Project Euler
// #81). The values may be negative. It panics if the rows have different
// lengths.
func GridMinTwoWays[E constraints.Number](g [][]E) E {
	rows, cols := checkGrid(g)
	if rows == 0 {
		return 0
	}
	// best[j] is the best sum to reach the current row's cell j.
	best := make([]E, cols)
	for i := range rows {
		for j := range cols {
			switch {
			case i == 0 && j == 0:
				best[j] = g[0][0]
			case i == 0:
				best[j] = best[j-1] + g[i][j]
			case j == 0:
				best[j] += g[i][j]
			default:
				best[j] = min(best[j], best[j-1]) + g[i][j]
			}
		}
	}
	return best[cols-1]
}

// GridMinThreeWays returns the smallest sum of a path from any cell of the
// left column to any cell of the right column that moves up, down and
// right (Project Euler #82). The values must not be negative (the result
// is unspecified otherwise). It panics if the rows have different lengths.
func GridMinThreeWays[E constraints.Number](g [][]E) E {
	rows, cols := checkGrid(g)
	if rows == 0 {
		return 0
	}
	// best[i] is the best sum of a path ending in the current column's
	// cell i.
	best := make([]E, rows)
	for i := range rows {
		best[i] = g[i][0]
	}
	for j := 1; j < cols; j++ {
		// Enter the column from the left, then let paths run down and up
		// it. With non-negative values one pass each way is enough: a path
		// never gains by turning around inside a column.
		for i := range rows {
			best[i] += g[i][j]
		}
		for i := 1; i < rows; i++ {
			best[i] = min(best[i], best[i-1]+g[i][j])
		}
		for i := rows - 2; i >= 0; i-- {
			best[i] = min(best[i], best[i+1]+g[i][j])
		}
	}
	m := best[0]
	for _, v := range best[1:] {
		m = min(m, v)
	}
	return m
}

// GridMinFourWays returns the smallest sum of a path from the top-left
// cell to the bottom-right cell that moves up, down, left and right
// (Project Euler #83), using Dijkstra's algorithm in O(n log n) for n
// cells. The values must not be negative (the result is unspecified
// otherwise). It panics if the rows have different lengths.
func GridMinFourWays[E constraints.Number](g [][]E) E {
	rows, cols := checkGrid(g)
	if rows == 0 {
		return 0
	}
	dist := make([]E, rows*cols)
	done := make([]bool, rows*cols)
	seen := make([]bool, rows*cols)
	dist[0], seen[0] = g[0][0], true
	q := &queue[E]{{cell: 0, dist: g[0][0]}}
	for q.Len() > 0 {
		cur := heap.Pop(q).(entry[E])
		if done[cur.cell] {
			continue // a stale entry: the cell was settled with less
		}
		done[cur.cell] = true
		if cur.cell == rows*cols-1 {
			break
		}
		i, j := cur.cell/cols, cur.cell%cols
		for _, d := range [4][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
			ni, nj := i+d[0], j+d[1]
			if ni < 0 || ni >= rows || nj < 0 || nj >= cols {
				continue
			}
			c := ni*cols + nj
			if nd := cur.dist + g[ni][nj]; !seen[c] || nd < dist[c] {
				dist[c], seen[c] = nd, true
				heap.Push(q, entry[E]{cell: c, dist: nd})
			}
		}
	}
	return dist[rows*cols-1]
}

// entry is a cell with the length of a path to it.
type entry[E constraints.Number] struct {
	cell int
	dist E
}

// queue is a min-heap of entries by dist.
type queue[E constraints.Number] []entry[E]

func (q queue[E]) Len() int           { return len(q) }
func (q queue[E]) Less(i, j int) bool { return q[i].dist < q[j].dist }
func (q queue[E]) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q *queue[E]) Push(x any)        { *q = append(*q, x.(entry[E])) }
func (q *queue[E]) Pop() any {
	old := *q
	x := old[len(old)-1]
	*q = old[:len(old)-1]
	return x
}
