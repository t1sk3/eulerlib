// Package vector provides an N-dimensional vector type backed by a slice of
// elements, along with fixed-size Vector2 and Vector3 types that embed it to
// add named X/Y/Z accessors and dimension-specific operations such as Cross
// and Perp.
package vector

import "github.com/t1sk3/eulerlib/v2/utils"

// Vector is an N-dimensional vector backed by a slice of elements.
type Vector[E utils.Number] struct {
	elements []E
}

// NewVector creates a vector from a slice of elements. The slice is copied,
// so later changes to elements do not affect the returned Vector.
func NewVector[E utils.Number](elements []E) *Vector[E] {
	e := make([]E, len(elements))
	copy(e, elements)
	return &Vector[E]{elements: e}
}

// Len returns the number of elements in the vector.
func (v *Vector[E]) Len() int {
	return len(v.elements)
}

// At returns the element at index i.
func (v *Vector[E]) At(i int) E {
	return v.elements[i]
}

// Set sets the element at index i to e.
func (v *Vector[E]) Set(i int, e E) {
	v.elements[i] = e
}

// Elements returns a copy of the underlying elements as a slice.
func (v *Vector[E]) Elements() []E {
	e := make([]E, len(v.elements))
	copy(e, v.elements)
	return e
}

// Swap swaps the elements at indices i and j.
func (v *Vector[E]) Swap(i, j int) {
	v.elements[i], v.elements[j] = v.elements[j], v.elements[i]
}

// Add returns the element-wise sum of v and o. It panics if their lengths differ.
func (v *Vector[E]) Add(o Vector[E]) *Vector[E] {
	v.mustMatch(o)
	elements := make([]E, v.Len())
	for i := range v.elements {
		elements[i] = v.elements[i] + o.elements[i]
	}

	return &Vector[E]{elements: elements}
}

// Sub returns the element-wise difference of v and o. It panics if their lengths differ.
func (v *Vector[E]) Sub(o Vector[E]) *Vector[E] {
	v.mustMatch(o)
	elements := make([]E, v.Len())
	for i := range v.elements {
		elements[i] = v.elements[i] - o.elements[i]
	}
	return &Vector[E]{elements: elements}
}

// Dot returns the dot product of v and o. It panics if their lengths differ.
func (v *Vector[E]) Dot(o Vector[E]) E {
	v.mustMatch(o)
	res := E(0)
	for i := range v.elements {
		res += v.elements[i] * o.elements[i]
	}
	return res
}

// Mul returns the element-wise product of v and o. It panics if their lengths differ.
func (v *Vector[E]) Mul(o Vector[E]) *Vector[E] {
	v.mustMatch(o)
	elements := make([]E, v.Len())
	for i := range v.elements {
		elements[i] = v.elements[i] * o.elements[i]
	}
	return &Vector[E]{elements: elements}
}

// Div returns the element-wise quotient of v and o. It panics if their lengths differ.
func (v *Vector[E]) Div(o Vector[E]) *Vector[E] {
	v.mustMatch(o)
	elements := make([]E, v.Len())
	for i := range v.elements {
		elements[i] = v.elements[i] / o.elements[i]
	}
	return &Vector[E]{elements: elements}
}

// mustMatch panics if v and o do not have the same length.
func (v *Vector[E]) mustMatch(o Vector[E]) {
	if v.Len() != o.Len() {
		panic("vector lengths do not match")
	}
}
