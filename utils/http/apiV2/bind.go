package apiV2

import (
	"net/http"
	"reflect"
)

// Bind decodes a request body and executes the Binder method of the
// payload structure.
// https://github.com/go-chi/render/blob/master/render.go
// Binder interface for managing request payloads.
type Binder interface {
	// Bind(r *http.Request) error
}

var (
	binderType = reflect.TypeOf(new(Binder)).Elem()
)

// Bind decodes a request body and executes the Binder method of the
// payload structure.
func Bind(r *http.Request, v Binder) error {
	if err := Decode(r, v); err != nil {
		return err
	}
	return binder(r, v)
}

// Executed bottom-up
func binder(r *http.Request, v Binder) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	// Call Binder on non-struct types right away
	if rv.Kind() != reflect.Struct {
		return nil // v.Bind(r)
	}

	// For structs, we call Bind on each field that implements Binder
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		if f.Type().Implements(binderType) {

			if isNil(f) {
				continue
			}

			fv := f.Interface().(Binder)
			if err := binder(r, fv); err != nil {
				return err
			}
		}
	}

	// We call it bottom-up
	//if err := v.Bind(r); err != nil {
	//	return err
	//}

	return nil
}

func isNil(f reflect.Value) bool {
	switch f.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return f.IsNil()
	default:
		return false
	}
}
