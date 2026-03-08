package eulerlib

import (
	"reflect"
	"testing"
)

func TestSameType(t *testing.T) {
	if !SameType[int, int]() {
		t.Fatal("SameType[int, int]() = false, want true")
	}
	if SameType[int, string]() {
		t.Fatal("SameType[int, string]() = true, want false")
	}
	if !SameType[signedIntegerToken, int64]() {
		t.Fatal("SameType[signedIntegerToken, int64]() = false, want true")
	}
	if !SameType[integerToken, uint32]() {
		t.Fatal("SameType[integerToken, uint32]() = false, want true")
	}
	if !SameType[floatToken, float32]() {
		t.Fatal("SameType[floatToken, float32]() = false, want true")
	}
	if !SameType[realNumberToken, int16]() {
		t.Fatal("SameType[realNumberToken, int16]() = false, want true")
	}
	if !SameType[complexToken, complex128]() {
		t.Fatal("SameType[complexToken, complex128]() = false, want true")
	}
	if !SameType[numberToken, complex64]() {
		t.Fatal("SameType[numberToken, complex64]() = false, want true")
	}
	if !SameType[comparableToken, string]() {
		t.Fatal("SameType[comparableToken, string]() = false, want true")
	}
}

func TestTypePredicates(t *testing.T) {
	intType := typeOf[int]()
	uintType := typeOf[uint]()
	floatType := typeOf[float64]()
	complexType := typeOf[complex64]()
	stringType := typeOf[string]()
	boolType := typeOf[bool]()

	if !IsSignedInteger(intType) || IsSignedInteger(uintType) {
		t.Fatal("IsSignedInteger classification mismatch")
	}
	if !IsUnsignedInteger(uintType) || IsUnsignedInteger(intType) {
		t.Fatal("IsUnsignedInteger classification mismatch")
	}
	if !IsInteger(intType) || !IsInteger(uintType) || IsInteger(floatType) {
		t.Fatal("IsInteger classification mismatch")
	}
	if !IsFloat(floatType) || IsFloat(intType) {
		t.Fatal("IsFloat classification mismatch")
	}
	if !IsRealNumber(intType) || !IsRealNumber(floatType) || IsRealNumber(complexType) {
		t.Fatal("IsRealNumber classification mismatch")
	}
	if !IsComplex(complexType) || IsComplex(floatType) {
		t.Fatal("IsComplex classification mismatch")
	}
	if !IsNumber(complexType) || !IsNumber(intType) || IsNumber(stringType) {
		t.Fatal("IsNumber classification mismatch")
	}
	if !IsComparable(stringType) || !IsComparable(boolType) || IsComparable(reflect.TypeOf([]int{})) {
		t.Fatal("IsComparable classification mismatch")
	}
}

func TestTypeOf(t *testing.T) {
	if typeOf[int]().Kind() != reflect.Int {
		t.Fatalf("typeOf[int]().Kind() = %v, want %v", typeOf[int]().Kind(), reflect.Int)
	}
}
