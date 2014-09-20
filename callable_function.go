package coroutine

// Copyright © 2012 Popog
import "reflect"

type CallableFunction reflect.Value

// Shorthand for Create(CreateCallableFunction(f))
func CreateFromFunc(f interface{}) *Thread {
	return Create(CreateCallableFunction(f))
}

// Shorthand for Wrap(CreateCallableFunction(f))
func WrapFromFunc(f interface{}) func(val ...interface{}) []interface{} {
	return Wrap(CreateCallableFunction(f))
}

// f must be a function. f must have a first parameter of type Thread.
// The return values of f will be wrapped into an []interface{}.
func CreateCallableFunction(f interface{}) CallableFunction {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		panic("CreateCallableFunction: f must be a function.")
	}

	if v.Type().NumIn() < 1 {
		panic("CreateCallableFunction: f must take at least 1 argument.")
	}

	if v.Type().In(0) != reflect.TypeOf(&Thread{}) {
		panic("CreateCallableFunction: First parameter of f must be Thread.")
	}

	return CallableFunction(v)
}

func (cf CallableFunction) Call(t *Thread, val ...interface{}) []interface{} {
	// convert the values
	vval := make([]reflect.Value, len(val)+1)
	vval[0] = reflect.ValueOf(t)
	for i, v := range val {
		vval[i+1] = reflect.ValueOf(v)
	}

	vret := reflect.Value(cf).Call(vval)

	ret := make([]interface{}, len(vret))
	for i, r := range vret {
		ret[i] = r.Interface()
	}
	return ret
}
