package main

import (
	"fmt"
  "math"
)


func for_loop(){
  sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}
	fmt.Println(sum)
}


// "while"
func while_loop(){
  sum := 1
  for sum < 1000 {
		sum += sum
	}
	fmt.Println(sum)
}


func infinite_loop(){
  for {
  }
}


// How a "if" with a short statement works
func pow(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim {
		return v
	} else {
		fmt.Printf("%g >= %g\n", v, lim)
	}
	// can't use v here, though
	return lim
}


func main() {
  for_loop()
  while_loop()
  //infinite_loop()
  fmt.Println(pow(3, 2, 10), pow(3, 2, 20))

}
