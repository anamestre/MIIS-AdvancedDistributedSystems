package main

import ("net"
        "fmt"
        "bufio"
        "os"
        "time"
        "strings")


/*
  Creates a server using the Port "port".
*/
func server(port string){
  fmt.Println("Launching server...")
  // listen on all interfaces,
  ln, _ := net.Listen("tcp", port)
  // accept connection on port
  conn, _ := ln.Accept()
  // run loop forever (or until ctrl-c)
  for {
    // will listen for message to process ending in newline (\n)
    message, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil { break }
    // output message received
    if strings.ToLower(message) == "stop\n"{
      os.Exit(1)
    } else {
      fmt.Print("--- Message Received:", string(message))
      // sample process for string received
      newmessage := strings.ToUpper(message)
      // send new string back to client
      conn.Write([]byte(newmessage + "\n"))
    }

  }
}


/*
  Connects with a server with the IP "ip" (which already contains the port).
*/
func client(ip string){
  // connect to this socket -> Dial = creates the connection
  fmt.Print("Checking if server has started" + "\n")
  conn, err := net.Dial("tcp", ip)
  // Handle error: try to connect again until server is up and running
  if err != nil {
    fmt.Print("Error: ", err)
    fmt.Print("\n")
    time.Sleep(time.Second)
    client(ip)
  }  else {
    fmt.Print("Server Up and Running" + "\n")
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
      fmt.Print("Message from server: " + message)
    }
  }
}



func main() {
    port := ":6002" // Port used for creating our own server
    ip := "127.0.0.1:6001" // Ip used for connecting to the other server
    go server(port)
    client(ip)
}
