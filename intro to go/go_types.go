
/*
NOTES
- Go does not have garbage collection, therefore we have to control the array sizes.
*/


package main

import (
	"fmt"
)

// There are no classes on Go, therefore we use Structs.
type Vertex struct {
	X int
	Y int
}


func pointers(){
  fmt.Println(" \n--- Pointer function")
  i, j := 42, 2701

	p := &i         // point to i
	fmt.Println("First value of *p:", *p) // read i through the pointer
	*p = 21         // set i through the pointer
	fmt.Println("Value should have changed, i:", i)  // see the new value of i

	p = &j         // point to j
	*p = *p / 37   // divide j through the pointer
	fmt.Println("New value of j:", j) // see the new value of j
}


func structs(){
  fmt.Println(" \n--- Structs function")
  v := Vertex{1, 2}
	v.X = 4
	fmt.Println(v.X)
}


func pointer_to_struct(){
  fmt.Println(" \n--- Pointer to struct function")
  v := Vertex{1, 2}
	p := &v
	p.X = 1e9 // instead of using (*p).X, we can use p.X
	fmt.Println(v)
}


func struct_literals(){
  fmt.Println(" \n--- Struct literals function")
  var (
  	v1 = Vertex{1, 2}  // has type Vertex
  	v2 = Vertex{X: 1}  // Y:0 is implicit
  	v3 = Vertex{}      // X:0 and Y:0
  	p  = &Vertex{1, 2} // has type *Vertex
  )
  fmt.Println(v1, p, v2, v3)
}


// their size is not mutable, you have to allocate memory in advance
func arrays(){
  fmt.Println(" \n--- Arrays function")
  var a [2]string
	a[0] = "Hello"
	a[1] = "World"
	fmt.Println(a[0], a[1])
	fmt.Println(a)

	primes := [6]int{2, 3, 5, 7, 11, 13}
	fmt.Println(primes)

  // Slicing an array: slices are like references to arrays
  var s []int = primes[1:4] // This is not a copy, if we change s, we change also primes

	fmt.Println(s)
}


func slicing(){
  fmt.Println(" \n--- Slicing function")
  s := []int{2, 3, 5, 7, 11, 13}

	s = s[1:4]
	fmt.Println(s)

	s = s[:2]
	fmt.Println(s)

	s = s[1:]
	fmt.Println(s)
}


func len_capacity(){
  fmt.Println(" \n--- Length and capacity function")
  s := []int{2, 3, 5, 7, 11, 13}
	printSlice(s)

	// Slice the slice to give it zero length.
	s = s[:0]
	printSlice(s)

	// Extend its length.
	s = s[:4]
	printSlice(s)

	// Drop its first two values.
	s = s[2:]
	printSlice(s)
}


/* Slices can be created with the built-in make function;
   this is how you create dynamically-sized arrays.
*/
func make_(){
  fmt.Println(" \n--- Make function")
  a := make([]int, 5)
	printSlice(a)

	b := make([]int, 0, 5)
	printSlice( b)

	c := b[:2]
	printSlice(c)

	d := c[2:5]
	printSlice(d)
}


func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}


/* If the backing array of s is too small to fit all the given values a bigger
   array will be allocated. The returned slice will point to the newly allocated array.
*/
func appending(){
  fmt.Println(" \n--- Appending function")
  var s []int
	printSlice(s)

	// append works on nil slices.
	s = append(s, 0)
	printSlice(s)

	// The slice grows as needed.
	s = append(s, 1)
	printSlice(s)

	// We can add more than one element at a time.
	s = append(s, 2, 3, 4)
	printSlice(s)
}


func range_(){
  fmt.Println(" \n--- Range function")
  var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}
	for i, v := range pow {
		fmt.Printf("2**%d = %d\n", i, v)
	}
}

type Vertexx struct {
	Lat, Long float64
}


func mapping(){
  fmt.Println(" \n--- Mapping function")
  var m = map[string]Vertexx{
  	"Bell Labs": Vertexx{
  		40.68433, -74.39967,
  	},
  	"Google": Vertexx{
  		37.42202, -122.08408,
  	},
  }
  fmt.Println(m)
}


func mutating_maps(){
  fmt.Println(" \n--- Mutating map function")
  m := make(map[string]int)

	m["Answer"] = 42
	fmt.Println("The value:", m["Answer"])

	m["Answer"] = 48
	fmt.Println("The value:", m["Answer"])

	delete(m, "Answer")
	fmt.Println("The value:", m["Answer"])

	v, ok := m["Answer"]
	fmt.Println("The value:", v, "Present?", ok)
}


func main(){
  pointers()
  structs()
  pointer_to_struct()
  struct_literals()
  arrays()
  slicing()
  len_capacity()
  make_()
  appending()
  range_()
  mapping()
  mutating_maps()
}
