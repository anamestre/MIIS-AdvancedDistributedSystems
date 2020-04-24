/*
	A TOUR OF GO:
	- Exercise: Equivalent Binary Trees
		https://tour.golang.org/concurrency/8
*/



package main

import (
	"golang.org/x/tour/tree"
	"fmt"
)

/*
type Tree struct {
    Left  *Tree
    Value int
    Right *Tree
} */

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	if t != nil {
    Walk(t.Left, ch)
    ch <- t.Value
    Walk(t.Right, ch)
  }
}


// Gets elements from t1, sends them through ch and then closes ch
func Walk_and_close(t1 *tree.Tree, ch chan int) {
	Walk(t1, ch)
	close(ch)
}


// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk_and_close(t1, ch1)
	go Walk_and_close(t2, ch2)

	for {
		// if ok == false, there are no more values to check
		elem1, ok1 := <- ch1
		elem2, ok2 := <- ch2

		if ok1 != ok2 {
			return false
		}
		if !ok1 {
			return true
		}
		if elem1 != elem2 {
			return false
		}
	}
}


func main() {
	t1 := tree.New(1)
  t2 := tree.New(1)
	fmt.Println(Same(t1, t2))
}
