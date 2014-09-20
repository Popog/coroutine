// Copyright © 2012 Popog
package coroutine

import "testing"

func foo(test *testing.T, t *Thread, a int) []interface{} {
	test.Log("foo", a)
	return t.Yield(2 * a)
}

func TestLua(test *testing.T) {
	test.Log()
	co := CreateFromFunc(func(t *Thread, a, b int) (interface{}, interface{}) {
		test.Log("co-body", a, b)
		r1 := foo(test, t, a+1)
		test.Log("co-body", r1)
		r2 := t.Yield(a+b, a-b)
		test.Log("co-body", r2)
		return b, "end"
	})

	test.Log("main", co.MustResume(1, 10))
	test.Log("main", co.MustResume("r"))
	test.Log("main", co.MustResume("x", "y"))
	if _, err := co.Resume("x", "y"); err != nil {
		test.Log("main", err)
	}
}

func TestWrap(test *testing.T) {
	test.Log()
	co := WrapFromFunc(func(t *Thread, a, b int) (interface{}, interface{}) {
		test.Log("co-body", a, b)
		r1 := foo(test, t, a+1)
		test.Log("co-body", r1)
		r2 := t.Yield(a+b, a-b)
		test.Log("co-body", r2)
		return b, "end"
	})

	test.Log("main", co(1, 10))
	test.Log("main", co("r"))
	test.Log("main", co("x", "y"))
}


func TestStop1(test *testing.T) {
	test.Log()
	co := CreateFromFunc(func(t *Thread, a, b int) (interface{}, interface{}) {
		test.Log("co-body", a, b)
		r1 := foo(test, t, a+1)
		test.Log("co-body", r1)
		r2 := t.Yield(a+b, a-b)
		test.Log("co-body", r2)
		return b, "end"
	})

	test.Log("main", co.MustResume(1, 10))
	co.Stop()
	if _, err := co.Resume("x", "y"); err != nil {
		test.Log("main", err)
	}
}


func TestStop2(test *testing.T) {
	test.Log()
	co := CreateFromFunc(func(t *Thread, a, b int) (interface{}, interface{}) {
		test.Log("co-body", a, b)
		r1 := foo(test, t, a+1)
		test.Log("co-body", r1)
		r2 := t.Yield(a+b, a-b)
		test.Log("co-body", r2)
		return b, "end"
	})

	co.Stop()
	if _, err := co.Resume("x", "y"); err != nil {
		test.Log("main", err)
	}
}