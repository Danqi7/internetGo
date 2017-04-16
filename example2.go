package main

import "fmt"

func fib(c chan int) {
	i, j := 0, 1
	for {
		i, j = i+j, i
		fmt.Printf("i=%d j=%d\n", i, j)
		c <- i
	}
}

func main() {
	c := make(chan int)
	go fib(c)
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}

}
