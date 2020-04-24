package main

import "net"
import "fmt"
import "bufio"
import "os"

func main() {

  // connect to this socket -> Dial = creates the connection
  conn, _ := net.Dial("tcp", "127.0.0.1:6001")
  for {
    // read in input from stdin
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Text to send: ")
    text, _ := reader.ReadString('\n')
    // send to socket
    fmt.Fprintf(conn, text + "\n")
    // listen for reply
    message, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil { break }
    fmt.Print("Message from server: "+message)
  }
}
