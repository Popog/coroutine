Implements collaborative multithreading functions as Lua coroutine. Description mostly copied from http://www.lua.org/manual/5.1/manual.html#2.11 as below.

A coroutine represents an independent thread of execution. Unlike threads in multithread systems, however, a coroutine only suspends its execution by explicitly calling a yield function.

You create a coroutine with a call to coroutine.Create. Its sole argument is a function that is the main function of the coroutine. The create function only creates a new coroutine and returns a handle to it (an object of type thread); it does not start the coroutine execution.

When you first call Thread.Resume, using the thread returned by coroutine.Create, the coroutine starts its execution, at the first line of its main function. Arguments passed to Thread.Resume are passed on to the coroutine main function. After the coroutine starts running, it runs until it terminates or yields.

A coroutine can terminate its execution in two ways: normally, when its main function returns (explicitly or implicitly, after the last instruction); and abnormally, if there is an unprotected panic. In the first case, Thread.Resume returns any values returned by the coroutine main function. In case of panics, Thread.Resume returns the value "recover"ed from the Thread. Thread.MustResume will panic instead of returning the recovered value.

A coroutine yields by calling Thread.Yield. When a coroutine yields, the corresponding Thread.Resume returns immediately, even if the yield happens inside nested function calls (that is, not in the main function, but in a function directly or indirectly called by the main function). In the case of a yield, Thread.Resume also returns any values passed to Thread.Yield. The next time you resume the same coroutine, it continues its execution from the point where it yielded, with the call to Thread.Yield returning any arguments passed to Thread.Resume.

Like coroutine.Create, the coroutine.Wrap function also creates a coroutine, but instead of returning the coroutine itself, it returns a function that, when called, resumes the coroutine. Any arguments passed to this function go as arguments to Thread.MustResume. It also returns all the values returned by Thread.MustResume. Like Thread.MustResume, it does not catch panics; panics are propagated to the caller.

As an example, consider the following code:

```go
func foo(t *Thread, a int) []interface{} {
	fmt.Println("foo", a)
	return t.Yield(2 * a)
}

func test() {
	co := CreateFromFunc(func(t *Thread, a, b int) (interface{}, interface{}) {
		fmt.Println("co-body", a, b)
		r1 := foo(t, a+1)
		fmt.Println("co-body", r1)
		r2 := t.Yield(a+b, a-b)
		fmt.Println("co-body", r2)
		return b, "end"
	})

	fmt.Println("main", co.MustResume(1, 10))
	fmt.Println("main", co.MustResume("r"))
	fmt.Println("main", co.MustResume("x", "y"))
	if _, err := co.Resume("x", "y"); err != nil {
		fmt.Println("main", err)
	}
}
```

When you run it, it produces the following output:

```
co-body 1 10
foo 2
main [4]
co-body [r]
main [11 -9]
co-body [x y]
main [10 end]
main coroutine: cannot resume dead coroutine.
```
