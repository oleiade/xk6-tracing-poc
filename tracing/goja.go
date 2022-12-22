package tracing

import "github.com/dop251/goja"

// isInstanceOf returns true if the given value is an instance of one of the
// given types. The types are specified as strings, which are the names of the
// target types constructors.
func isInstanceOf(rt *goja.Runtime, v goja.Value, instanceOf ...string) bool {
	var valid bool
	for _, t := range instanceOf {
		instanceOfConstructor := rt.Get(t)
		if valid = v.ToObject(rt).Get("constructor").SameAs(instanceOfConstructor); valid {
			break
		}
	}
	return valid
}

// isNullish checks if the given value is nullish, i.e. nil, undefined or null.
// This helper function reproduces the behavior of Javascripts nullish coalescing
// operator (??).
func isNullish(value goja.Value) bool {
	return value == nil || goja.IsUndefined(value) || goja.IsNull(value)
}
