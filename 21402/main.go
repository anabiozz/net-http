package main

import (
	"fmt"
	"runtime"
	"time"
)

type Test struct {
	data string
}

func main() {
	fmt.Println("A")

	o := &Test{"some data"}
	runtime.SetFinalizer(o, func(o *Test) {
		fmt.Printf("Finalized %p\n", o)
	})

	defer runtime.KeepAlive(o)

	runtime.GC()
	time.Sleep(1 * time.Second)

	fmt.Println("B")
}
