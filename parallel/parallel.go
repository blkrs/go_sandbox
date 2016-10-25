package main

import (
//    "bufio"
    "fmt"
    "time"
//    "math"
//    "os"
//    "runtime"
)

func calculate(input chan int) {
        go func() {
          liczba := <- input
	  fmt.Printf("Starting calculation for %d ", liczba)
          for i:=2;i < liczba; i++ {
            if liczba % i ==0 {
             fmt.Printf("Not prime, divided by %d\n", i)
             return
            }
          }
          fmt.Printf(" XXXXXXXXXXXXXXXXXXXXXX Prime\n")
        }()
}


func main() {
 fmt.Printf("Hello world\n")
 inputs := make(chan int)
 for i:=500;i < 600; i++ {
   calculate(inputs)
 }
 for i:=500;i < 600; i++ {
   inputs <- i*i - 1
 }
 fmt.Printf("Created 100 inputs\n")
 time.Sleep(3000*time.Millisecond)
}
