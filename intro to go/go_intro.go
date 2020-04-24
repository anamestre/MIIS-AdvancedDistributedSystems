// https://tour.golang.org/basics


package main

import (
	"fmt"
	"math/rand"
  "math"
)

// func add(x, y int) int { // We can also write it like this since both are int
func add(x int, y int) int {
	return x + y
}

func swap(x, y string) (string, string) {
	return y, x
}

func split(sum int) (x, y int) { // We can write the variables we are going to return
	x = sum * 4 / 9
	y = sum - x
	return
}

var c, python, java bool // Global variable, duh.
var i, j int = 1, 2

const Pi = 3.14

func main() {
  var i int // Local variable, duuuh.
  /* If you don't assign anything to the variable, it has a default value,
     int = 0, string = "", bool = false ...
  */
  fmt.Println(i)

	fmt.Println("My favorite number is", rand.Intn(10))
  fmt.Println(math.Pi) // We have to capitalize "Pi" if we are calling it from the outside
  fmt.Println(add(42, 13))
  fmt.Println(swap("hello", "world"))
  fmt.Println(split(17))

  var c, python, java = true, false, "no!" // The system defines the type of the variable according to what you are assigning
  // c, python, java := true, false, "no!" We can user ":=" instead of "var"
  fmt.Println(i, j, c, python, java)

  // Casting:
  var k int = 42
  var f float64 = float64(k)
  var u uint = uint(f)
  fmt.Println(u)

}

/*
Notes:
- Basic types
    bool

    string

    int  int8  int16  int32  int64
    uint uint8 uint16 uint32 uint64 uintptr

    byte // alias for uint8

    rune // alias for int32
         // represents a Unicode code point

    float32 float64

    complex64 complex128
*/
