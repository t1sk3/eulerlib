package eulerlib

import "testing"

func TestTestFunction(t *testing.T) {
	if got := Test(); got != "test" {
		t.Fatalf("Test() = %q, want %q", got, "test")
	}
}
