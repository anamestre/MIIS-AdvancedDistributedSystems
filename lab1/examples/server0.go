package main

import "net"
import "fmt"
import "bufio"
import "strings" // only needed below for sample processing

func main() {

  fmt.Println("Launching server...")

  // listen on all interfaces,
  // listener object waiting for connection in this particular port
  // it is opening the port
  ln, _ := net.Listen("tcp", ":6001")

  // accept connection on port
  conn, _ := ln.Accept()

  // run loop forever (or until ctrl-c)
  for {
    // will listen for message to process ending in newline (\n)
    message, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil { break }
    // output message received
    fmt.Print("Message Received:", string(message))
    // sample process for string received
    newmessage := strings.ToUpper(message)

    // send new string back to client
    conn.Write([]byte(newmessage + "\n"))
  }
}
