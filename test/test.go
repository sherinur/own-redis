package main

import (
	"fmt"
	"time"
)

func f1(x int) {
	for i := 0; i <= x; i++ {
		fmt.Println("XXX:", i, "[", time.Now().Second(), ":", time.Now().UnixMilli(), "]")
	}
}

func f2(x int) {
	for i := 0; i <= x; i++ {
		fmt.Println("YYY:", i, "[", time.Now().Second(), ":", time.Now().UnixMilli(), "]")
	}
}

func main() {
	f1(10)
	f2(10)

	time.Sleep(time.Millisecond)
}
