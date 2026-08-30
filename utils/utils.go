package utils

import "reflect"

type signedIntegerToken struct{}
type unsignedIntegerToken struct{}
type integerToken struct{}
type floatToken struct{}
type realNumberToken struct{}
type complexToken struct{}
type numberToken struct{}
type comparableToken struct{}

// SameType returns whether A and B are of the same type or type family token.
func SameType[A, B any]() bool {
	ta, tb := reflect.TypeOf((*A)(nil)).Elem(), reflect.TypeOf((*B)(nil)).Elem()
	if ta == tb {
		return true
	}

	checks := map[reflect.Type]func(reflect.Type) bool{
		typeOf[signedIntegerToken]():   IsSignedInteger,
		typeOf[unsignedIntegerToken](): IsUnsignedInteger,
		typeOf[integerToken]():         IsInteger,
		typeOf[floatToken]():           IsFloat,
		typeOf[realNumberToken]():      IsRealNumber,
		typeOf[complexToken]():         IsComplex,
		typeOf[numberToken]():          IsNumber,
		typeOf[comparableToken]():      IsComparable,
	}

	if check, ok := checks[ta]; ok {
		return check(tb)
	}
	if check, ok := checks[tb]; ok {
		return check(ta)
	}
	return false
}

// typeOf returns the reflect.Type of T without needing a value of T.
func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// IsSignedInteger returns true if t's kind is a signed integer type.
func IsSignedInteger(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

// IsUnsignedInteger returns true if t's kind is an unsigned integer type.
func IsUnsignedInteger(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return true
	default:
		return false
	}
}

// IsInteger returns true if t's kind is a signed or unsigned integer type.
func IsInteger(t reflect.Type) bool {
	return IsSignedInteger(t) || IsUnsignedInteger(t)
}

// IsFloat returns true if t's kind is a floating-point type.
func IsFloat(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// IsRealNumber returns true if t's kind is an integer or floating-point type.
func IsRealNumber(t reflect.Type) bool {
	return IsInteger(t) || IsFloat(t)
}

// IsComplex returns true if t's kind is a complex number type.
func IsComplex(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Complex64, reflect.Complex128:
		return true
	default:
		return false
	}
}

// IsNumber returns true if t's kind is a real or complex number type.
func IsNumber(t reflect.Type) bool {
	return IsRealNumber(t) || IsComplex(t)
}

// IsComparable returns true if t's kind is a number, string, or bool type.
func IsComparable(t reflect.Type) bool {
	return IsNumber(t) || t.Kind() == reflect.String || t.Kind() == reflect.Bool
}
