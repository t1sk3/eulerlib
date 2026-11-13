package roman

import "testing"

func TestFormat(t *testing.T) {
	for n, want := range map[int]string{
		1: "I", 4: "IV", 9: "IX", 14: "XIV", 40: "XL", 49: "XLIX", 90: "XC", 400: "CD", 606: "DCVI",
		900: "CM", 1994: "MCMXCIV", 2024: "MMXXIV", 3999: "MMMCMXCIX", 4000: "MMMM", 5000: "MMMMM",
	} {
		if got := Format(n); got != want {
			t.Errorf("Format(%d) = %q, want %q", n, got, want)
		}
	}
	defer func() {
		if recover() == nil {
			t.Error("Format(0) should panic")
		}
	}()
	Format(0)
}

func TestParse(t *testing.T) {
	for n := 1; n <= 12000; n++ {
		if got, err := Parse(Format(n)); err != nil || got != n {
			t.Fatalf("Parse(Format(%d) = %q) = %d, %v", n, Format(n), got, err)
		}
	}
	// The non-minimal spellings of Project Euler #89.
	for s, want := range map[string]int{
		"IIII": 4, "VIIII": 9, "XIIIIII": 16, "XVI": 16, "XIIIIIIIII": 19, "XVIIII": 19, "XIX": 19,
		"XXXXVIIII": 49, "XXXXIX": 49, "XLVIIII": 49, "XLIX": 49, "MCCCCCCVI": 1606, "MDCVI": 1606,
		"DCCCC": 900, "LXXXX": 90, "MMMMCCCCLXXXXVIIII": 4499, "IIIIIIIIIIIIIIII": 16,
	} {
		if got, err := Parse(s); err != nil || got != want {
			t.Errorf("Parse(%q) = %d, %v, want %d", s, got, err, want)
		}
	}
	for _, s := range []string{
		"", "A", "iv", "X I", "IIX", "IL", "IC", "ID", "IM", "VX", "VL", "XD", "XM", "LC", "DM",
		"IXX", "IXC", "XCC", "IVX", "XLX", "CMM", "IIV", "IXI", "IVI", "XCX", "CMC", "CMD", "XCL",
	} {
		if got, err := Parse(s); err == nil {
			t.Errorf("Parse(%q) = %d, want an error", s, got)
		}
	}
}

// Project Euler #89 asks how many characters are saved by rewriting each
// numeral in its minimal form.
func TestSavings(t *testing.T) {
	for s, want := range map[string]int{
		"IIII": 2, "VIIII": 3, "XXXXVIIII": 5, "MCCCCCCVI": 4, "MMMMDCCCCI": 3, "MCMXCIV": 0, "XVI": 0,
	} {
		n, err := Parse(s)
		if err != nil {
			t.Fatalf("Parse(%q): %v", s, err)
		}
		if saved := len(s) - len(Format(n)); saved != want {
			t.Errorf("%q -> %q saves %d characters, want %d", s, Format(n), saved, want)
		}
	}
}
