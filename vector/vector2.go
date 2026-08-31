package vector

import "github.com/t1sk3/eulerlib/v2/utils"

// Vector2 is a 2-dimensional vector. It embeds Vector to inherit the
// element-wise arithmetic operations, and adds named X/Y accessors and
// 2D-specific operations such as Cross and Perp.
type Vector2[E utils.SignedNumber] struct {
	Vector[E]
}

// NewVector2 creates a 2D vector from its x and y components.
func NewVector2[E utils.SignedNumber](x, y E) *Vector2[E] {
	return &Vector2[E]{Vector[E]{elements: []E{x, y}}}
}

// AsVector2 adopts an existing Vector (no copy — they share backing storage).
func AsVector2[E utils.SignedNumber](v *Vector[E]) *Vector2[E] {
	if v.Len() != 2 {
		panic("vector is not of length 2")
	}
	return &Vector2[E]{*v}
}

// X returns the x component.
func (v *Vector2[E]) X() E { return v.elements[0] }

// Y returns the y component.
func (v *Vector2[E]) Y() E { return v.elements[1] }

// SetX sets the x component.
func (v *Vector2[E]) SetX(x E) { v.elements[0] = x }

// SetY sets the y component.
func (v *Vector2[E]) SetY(y E) { v.elements[1] = y }

// Cross is the z-component of the 3D cross product (the "perp dot" product).
func (v *Vector2[E]) Cross(o Vector2[E]) E {
	return v.X()*o.Y() - v.Y()*o.X()
}

// Perp returns the perpendicular vector (-y, x).
func (v *Vector2[E]) Perp() *Vector2[E] { return NewVector2(-v.Y(), v.X()) }

// Add shadows the promoted Vector.Add so chaining stays in Vector2.
func (v *Vector2[E]) Add(o Vector2[E]) *Vector2[E] {
	return &Vector2[E]{*v.Vector.Add(o.Vector)}
}

// Sub shadows the promoted Vector.Sub so chaining stays in Vector2.
func (v *Vector2[E]) Sub(o Vector2[E]) *Vector2[E] {
	return &Vector2[E]{*v.Vector.Sub(o.Vector)}
}

// Mul shadows the promoted Vector.Mul so chaining stays in Vector2.
func (v *Vector2[E]) Mul(o Vector2[E]) *Vector2[E] {
	return &Vector2[E]{*v.Vector.Mul(o.Vector)}
}

// Div shadows the promoted Vector.Div so chaining stays in Vector2.
func (v *Vector2[E]) Div(o Vector2[E]) *Vector2[E] {
	return &Vector2[E]{*v.Vector.Div(o.Vector)}
}
