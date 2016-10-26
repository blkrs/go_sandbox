package main

import (
	//    "bufio"
	"fmt"
	"time"
	//    "math"
	//    "os"
	//    "runtime"
)

var min int
var max int

func calculate(input chan int, output chan int) {
	go func() {
		liczba := <-input
		for i := 2; i < liczba; i++ {
			if liczba%i == 0 {
				output <- 0
				return
			}
		}
		time.Sleep(5000)
		output <- liczba
	}()
}

func collect(output chan int) {
	go func() {
		for i := min; i < max; i++ {
			liczba := <-output
			fmt.Printf("[%d] collected: %d\n", i, liczba)
		}
	}()
}

func main() {
	fmt.Printf("Parallel prime number calculator\n")
	min, max = 500, 600
	inputs := make(chan int)
	output := make(chan int)
	// create goroutine for result collection
	collect(output)

	// create goroutines for calculation
	for i := min; i < max; i++ {
		calculate(inputs, output)
	}

	// distribute numbers through channels
	for i := min; i < max; i++ {
		inputs <- i
	}

	fmt.Printf("Created %d inputs\n", max-min)
	time.Sleep(3000 * time.Millisecond)
}
