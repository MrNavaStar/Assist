package ops

import "reflect"

// Ternary is the standard ternary operator which returns posVal
// if the boolean condition is true, otherwise negVal.
func Ternary(cond bool, posVal any, negVal any) any {
	if cond {
		return posVal
	}
	return negVal
}

func IsEmpty(i interface{}) bool {
	return reflect.ValueOf(i).IsZero()
}
