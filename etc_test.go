package eulerlib

import "testing"

func TestReverseString(t *testing.T) {
	input := "abcdefg"

	output := ReverseString(input)
	if output != "gfedcba" {
		t.Errorf("ReverseString(%s) returned %s, expected %s", input, output, "gfedcba")
	}
}
