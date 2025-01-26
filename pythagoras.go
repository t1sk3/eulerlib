package eulerlib

// checks if the given triplet is a pythagorean triplet
func IsTriplet[E Integer](a E, b E, c E) bool {
	return c*c == b*b+a*a
}
