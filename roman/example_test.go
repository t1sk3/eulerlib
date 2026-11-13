package roman_test

import (
	"fmt"

	"github.com/t1sk3/eulerlib/v4/roman"
)

func ExampleFormat() {
	fmt.Println(roman.Format(1994))
	// Output: MCMXCIV
}

func ExampleParse() {
	n, err := roman.Parse("XXXXVIIII") // valid, but not minimal
	fmt.Println(n, err, roman.Format(n))
	// Output: 49 <nil> XLIX
}
