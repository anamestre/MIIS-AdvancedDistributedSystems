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
func server(port string, ch chan string){
  fmt.Println("Launching server in port: " + port + "...")
  // listen on all interfaces,
  ln, _ := net.Listen("tcp", ":"+port)

  // run loop forever (or until ctrl-c)
  for {
    // accept connection on port
    conn, _ := ln.Accept()
    go serverConnection(conn, ch)
  }
}

/*
  Handles a connection (conn) for the server.
*/
func serverConnection(conn net.Conn, ch chan string){
  // will listen for message to process ending in newline (\n)
  message, err := bufio.NewReader(conn).ReadString('\n')
  if err != nil {
        fmt.Println("client left..")
        conn.Close()
        return
    }
  // output message received
  text := strings.TrimSpace(string(message))
  temp := strings.Split(string(text), ";")
  message = temp[0]
  ip := temp[1]
  if strings.ToLower(message) == "stop"{
    ip := strings.TrimSpace(string(temp[1]))
    fmt.Println("client " + ip + " left..")
    ch <- ip
    conn.Close()
    return
  } else {
    fmt.Println("\n--- Message Received from " + ip + ": " + string(message))
    // sample process for string received
    newmessage := strings.ToUpper(message)
    // send new string back to client
    conn.Write([]byte(newmessage + "\n"))
    serverConnection(conn, ch)
  }
}


/*
  Makes the first connections with the corresponding servers.
  Doesn't stop until a connection is made.
*/
func connect(ip string) net.Conn{
  fmt.Print("Checking if server " + ip + " has started" + "\n")
  conn, err := net.Dial("tcp", ip)
  // Handle error: try to connect again until server is up and running
  if err != nil {
      fmt.Println("Error: ", err)
      time.Sleep(time.Second)
      return connect(ip)
  }  else {
      fmt.Print("Server Up and Running" + "\n")
      return conn
  }
}


/*
  Awaits for an stop message on channel from this node's server and
  updates the list of connected nodes.
*/
func waitForStop(ch chan string, connections map[string]Connection){
  for {
    elem, _ := <- ch
    delete(connections, elem) // remove ip from our list of connections

    // If there are no more connections, we shut down the process
    if len(connections) == 0 {
      os.Exit(1)
    }
  }

}



/*
  Shows the usage of the program. How should it be executed.
*/
func usage(){
    fmt.Println("The path to a configuration file is needed.")
    fmt.Println("-- for example: ./machine config.txt")
}


/*
  Reads the configuration file. Obtains the port and clients ip's.
*/
func getConfig(myFile string) (string, []string){
  file, err := os.Open(myFile)
  if err != nil {
      fmt.Println(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  first := true
  serverIP := ""
  ips := []string{}
  for scanner.Scan() {
      if first {
        serverIP = scanner.Text() // The first line is the server ip
        first = false
      } else {
        ips = append(ips, scanner.Text()) // Saving every client's ip
      }
  }
  return serverIP, ips
}


type Connection struct {
    Channel chan string
    Connection net.Conn
}

/*
  Connects with a server with the IP "ip" (which already contains the port).
*/
func clients(ips []string, ch chan string, myIP string) {
    // Connect with the given ips and save their connections
    // and channels in a map.
    connections := map[string]Connection{}
    for _, ip := range ips {
      chIP := make(chan string)
      c := Connection{chIP, connect(ip)}
      connections[ip] = c
    }

    // We handle every client (connected ip) concurrently.
    for ip, c := range connections{
        go handleClient(c, myIP, ip)
    }

    // Checking whether an "stop" message has arrived
    go waitForStop(ch, connections)


    // Reading messages from terminal and sending them to the corresponding servers.
    for {
        // read in input from stdin
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Text to send: ")
        text, _ := reader.ReadString('\n')

        // Send to connections
        for _, c := range connections {
          c.Channel <- text
        }

        // If the input message is "stop", we close the program.
        temp := strings.TrimSpace(string(text))
        temp = strings.ToLower(temp)
        if temp == "stop"{
          //time.Sleep(time.Second/4)
          fmt.Println("Program has been stopped")
          os.Exit(1)
        }
      }
}


/*
  Send the written message to its corresponding server.
  c = connection with a client.
  myIP = the IP corresponding to the server of this program (we use it as an id)
*/
func handleClient(c Connection, myIP string, ip string){
    for {
      text, _ := <- c.Channel
      conn := c.Connection
      temp := strings.TrimSpace(string(text)) + ";" + myIP
      fmt.Fprintf(conn, temp + "\n")
      // listen for reply
      message, _ := bufio.NewReader(conn).ReadString('\n')
      //if err != nil { break }
      fmt.Print("\nMessage from server " + ip + ": " + message)
    }
}


func main() {
    if len(os.Args) != 2{
      usage()
    } else {
      myIP, ips := getConfig(os.Args[1])
      port := strings.Split(myIP, ":")[1]
      ch := make(chan string)
      go server(port, ch)
      clients(ips, ch, myIP)
  }
}
