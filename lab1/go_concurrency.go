package main

import (
	"fmt"
	"time"
)


func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}


func goroutine(){
  fmt.Println(" \n--- Go-routine function")
  go say("world") // A goroutine is a lightweight thread managed by the Go runtime.
  say("hello")
}


func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}


/* Channels are a typed conduit through which you can send and receive
   values with the channel operator, <-.
   Instead of using shared memory, we are using channels to share information.
*/
func channels() {
  fmt.Println(" \n--- Channels function")
	s := []int{7, 2, 8, -9, 4, 0}

	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c // receive from c

	fmt.Println(x, y, x+y)
}


func buffered_channel(){
  ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	fmt.Println(<-ch)
	fmt.Println(<-ch)
  // fmt.Println(<-ch) If we add a third one, we'll get an error.
  /*
    We have to think about channels as if they were queues.
    How to deal with channel size:
    ch := make(chan int, 2)
  	ch <- 1
  	ch <- 2
  	fmt.Println(<-ch)
  	fmt.Println(<-ch) // Since we have removed an element from the channel, we can add another one.
  	ch <- 3
  	fmt.Println(<-ch)
  */
}


func main() {
  goroutine()
  channels()
  buffered_channel()

}
