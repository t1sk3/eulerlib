package vector_n

import "github.com/t1sk3/eulerlib/utils"

type Vector[E utils.Number] struct {
	elements []E
}

func NewVector[E utils.Number](elements []E) *Vector[E] {
	return &Vector[E]{
		elements: elements,
	}
}

func NewVector2[E utils.Number]() *Vector[E] {
	return &Vector[E]{
		elements: []E{0, 0},
	}
}

func NewVector3[E utils.Number]() *Vector[E] {
	return &Vector[E]{
		elements: []E{0, 0, 0},
	}
}

func (v Vector[E]) Len() int {
	return len(v.elements)
}

func (v Vector[E]) Swap(i, j int) {
	v.elements[i], v.elements[j] = v.elements[j], v.elements[i]
}

func (v Vector[E]) Add(o Vector[E]) *Vector[E] {
	if v.Len() != o.Len() {
		panic("vector lengths do not match")
	}
	var elements []E
	for i := range v.elements {
		elements = append(elements, v.elements[i]+o.elements[i])
	}

	return NewVector(elements)
}

func (v Vector[E]) Sub(o Vector[E]) *Vector[E] {
	if v.Len() != o.Len() {
		panic("vector lengths do not match")
	}
	var elements []E
	for i := range v.elements {
		elements = append(elements, v.elements[i]-o.elements[i])
	}
	return NewVector(elements)
}

func (v Vector[E]) Dot(o Vector[E]) E {
	if v.Len() != o.Len() {
		panic("vector lengths do not match")
	}
	res := E(0)
	for i := range v.elements {
		res += v.elements[i] + o.elements[i]
	}
	return res
}

func (v Vector[E]) Mul(o Vector[E]) *Vector[E] {
	if v.Len() != o.Len() {
		panic("vector lengths do not match")
	}
	var elements []E
	for i := range v.elements {
		elements = append(elements, v.elements[i]*o.elements[i])
	}
	return NewVector(elements)
}

func (v Vector[E]) Div(o Vector[E]) *Vector[E] {
	if v.Len() != o.Len() {
		panic("vector lengths do not match")
	}
	var elements []E
	for i := range v.elements {
		elements = append(elements, v.elements[i]/o.elements[i])
	}
	return NewVector(elements)
}
