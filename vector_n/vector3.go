package vector_n

import "github.com/t1sk3/eulerlib/v2/utils"

type Vector3[E utils.SignedNumber] struct {
	Vector[E]
}

func NewVector3[E utils.SignedNumber](x, y, z E) *Vector3[E] {
	return &Vector3[E]{Vector[E]{elements: []E{x, y, z}}}
}

// AsVector3 adopts an existing Vector (no copy — they share backing storage).
func AsVector3[E utils.SignedNumber](v *Vector[E]) *Vector3[E] {
	if v.Len() != 3 {
		panic("vector is not of length 3")
	}
	return &Vector3[E]{*v}
}

func (v *Vector3[E]) X() E { return v.elements[0] }
func (v *Vector3[E]) Y() E { return v.elements[1] }
func (v *Vector3[E]) Z() E { return v.elements[2] }

func (v *Vector3[E]) SetX(x E) { v.elements[0] = x }
func (v *Vector3[E]) SetY(y E) { v.elements[1] = y }
func (v *Vector3[E]) SetZ(z E) { v.elements[2] = z }

func (v *Vector3[E]) Cross(o Vector3[E]) *Vector3[E] {
	return NewVector3(
		v.Y()*o.Z()-v.Z()*o.Y(),
		v.Z()*o.X()-v.X()*o.Z(),
		v.X()*o.Y()-v.Y()*o.X(),
	)
}

// Add shadows the promoted Vector.Add so chaining stays in Vector3.
func (v *Vector3[E]) Add(o Vector3[E]) *Vector3[E] {
	return &Vector3[E]{*v.Vector.Add(o.Vector)}
}

// Sub shadows the promoted Vector.Sub so chaining stays in Vector3.
func (v *Vector3[E]) Sub(o Vector3[E]) *Vector3[E] {
	return &Vector3[E]{*v.Vector.Sub(o.Vector)}
}

// Mul shadows the promoted Vector.Mul so chaining stays in Vector3.
func (v *Vector3[E]) Mul(o Vector3[E]) *Vector3[E] {
	return &Vector3[E]{*v.Vector.Mul(o.Vector)}
}

// Div shadows the promoted Vector.Div so chaining stays in Vector3.
func (v *Vector3[E]) Div(o Vector3[E]) *Vector3[E] {
	return &Vector3[E]{*v.Vector.Div(o.Vector)}
}
