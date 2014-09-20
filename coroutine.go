// Copyright © 2012 Popog
package coroutine

import (
	"errors"
)

type yieldData struct {
	val []interface{}
	err interface{}
}

type Thread struct {
	resume chan []interface{} // the data sent from a call to resume
	yield  chan yieldData     // the data sent from a call to yield
	status Status
}

type Status int

const (
	Suspended Status = iota
	Running
	Dead
)

type closeThreadError struct{}

type Callable interface {
	Call(t *Thread, val ...interface{}) []interface{}
}

// Creates a coroutine
func Create(c Callable) *Thread {
	t := &Thread{
		resume: make(chan []interface{}),
		yield:  make(chan yieldData),
		status: Suspended,
	}

	// Kick off the coroutine
	go func() {
		var yval yieldData

		// Handle and forward panics and set the state to dead
		defer func() {
			t.status = Dead

			switch err := recover(); err.(type) {
			case nil, closeThreadError:
			default:
				yval = yieldData{err: err}
			}

			t.yield <- yval
		}()

		// wait for the parameters
		val := <-t.resume
		t.status = Running

		// get the results
		yval = yieldData{val: c.Call(t, val...)}
	}()
	return t
}

func Wrap(c Callable) func(val ...interface{}) []interface{} {
	co := Create(c)
	return func(val ...interface{}) []interface{} {
		return co.MustResume(val...)
	}
}

func (t *Thread) Status() Status {
	return t.status
}

func (t *Thread) Resume(val ...interface{}) (yval[]interface{}, err interface{}) {
	if t.Status() == Dead {
		return nil, errors.New("coroutine: cannot resume dead coroutine.")
	}

	t.resume <- val
	yd := <-t.yield

	// if we're dead, then clearly we're not suspended
	if t.Status() != Dead {
		t.status = Suspended
	}

	return yd.val, yd.err
}

// panics if err != nil
func (t *Thread) MustResume(val ...interface{}) []interface{} {
	v, err := t.Resume(val...)
	if err != nil {
		panic(err)
	}
	return v
}

func (t *Thread) Stop() {
	if t.Status() == Dead {
		return
	}
	close(t.resume)
	t.status = Dead
}

// Yields results from the coroutine
func (t *Thread) Yield(val ...interface{}) []interface{} {
	t.yield <- yieldData{val: val}

	if rval, ok := <-t.resume; ok {
		t.status = Running
		return rval
	}

	panic(closeThreadError{})
	return nil
}

func (t *Thread) YieldError(err error) []interface{} {
	t.yield <- yieldData{err: err}

	if rval, ok := <-t.resume; ok {
		t.status = Running
		return rval
	}

	panic(closeThreadError{})
	return nil
}
