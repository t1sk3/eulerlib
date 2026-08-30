// Package pythagoras provides small geometry helpers such as Pythagorean
// triplet checking and degree-to-radian conversion.
package pythagoras

import (
	"github.com/t1sk3/eulerlib/v2/utils"
)

// IsTriplet returns true if (a, b, c) is a Pythagorean triplet, i.e. a²+b²=c².
func IsTriplet[E utils.Integer](a E, b E, c E) bool {
	return c*c == b*b+a*a
}
