package pythagoras

import (
	"github.com/t1sk3/eulerlib/utils"
)

// checks if the given triplet is a pythagorean triplet
func IsTriplet[E utils.Integer](a E, b E, c E) bool {
	return c*c == b*b+a*a
}
