package vector_n

import "github.com/t1sk3/eulerlib/v2/utils"

type Vector[E utils.Number] struct {
	elements []E
}

func NewVector[E utils.Number](elements []E) *Vector[E] {
	e := make([]E, len(elements))
	copy(e, elements)
	return &Vector[E]{elements: e}
}

func (v *Vector[E]) Len() int {
	return len(v.elements)
}

func (v *Vector[E]) At(i int) E {
	return v.elements[i]
}

func (v *Vector[E]) Set(i int, e E) {
	v.elements[i] = e
}

func (v *Vector[E]) Elements() []E {
	e := make([]E, len(v.elements))
	copy(e, v.elements)
	return e
}

func (v *Vector[E]) Swap(i, j int) {
	v.elements[i], v.elements[j] = v.elements[j], v.elements[i]
}

func (v *Vector[E]) Add(o Vector[E]) *Vector[E] {
	v.mustMatch(o)
	elements := make([]E, v.Len())
	for i := range v.elements {
		elements[i] = v.elements[i] + o.elements[i]
	}

	return &Vector[E]{elements: elements}
}

func (v *Vector[E]) Sub(o Vector[E]) *Vector[E] {
	v.mustMatch(o)
	elements := make([]E, v.Len())
	for i := range v.elements {
		elements[i] = v.elements[i] - o.elements[i]
	}
	return &Vector[E]{elements: elements}
}

func (v *Vector[E]) Dot(o Vector[E]) E {
	v.mustMatch(o)
	res := E(0)
	for i := range v.elements {
		res += v.elements[i] * o.elements[i]
	}
	return res
}

func (v *Vector[E]) Mul(o Vector[E]) *Vector[E] {
	v.mustMatch(o)
	elements := make([]E, v.Len())
	for i := range v.elements {
		elements[i] = v.elements[i] * o.elements[i]
	}
	return &Vector[E]{elements: elements}
}

func (v *Vector[E]) Div(o Vector[E]) *Vector[E] {
	v.mustMatch(o)
	elements := make([]E, v.Len())
	for i := range v.elements {
		elements[i] = v.elements[i] / o.elements[i]
	}
	return &Vector[E]{elements: elements}
}

func (v *Vector[E]) mustMatch(o Vector[E]) {
	if v.Len() != o.Len() {
		panic("vector lengths do not match")
	}
}
