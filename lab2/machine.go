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

func allConnections(first bool) bool {
  if first{
    fmt.Println("Checking connections for first wave")
  } else {
    fmt.Println("Checking connections for second wave")
  }
  for _, conn := range connections {
      if first {
        if ! conn.firstWave {return false}
      } else {
        if ! conn.secondWave {return false}
      }
  }
  return true
}

/*
  Handles a connection (conn) for the server.
*/
func serverConnection(conn net.Conn, ch chan string){
  // will listen for message to process ending in newline (\n)
  if firstWave{
    for {
      if len(connections) > 0 {break}
    }
  }
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
  id := temp[1] // De moment la id serà la ip de cada server
  fmt.Println("From server: Message received from: "+ id + " -> " + message)
  if message == "hello"{
    if firstWave {
      if len(connections) > 1 { // This means that it has neighbours and not only a parent
        fmt.Println("From server HELLO: has neighbors, sending message to neighbors")
         ch <- message
      }
      firstWave = false
      // marco aquesta com a parent
      fmt.Print("From server:")
      fmt.Println(len(connections))
      fmt.Println("From server: marking as parent: " + id)
      //time.Sleep(time.Second*10)
      parent = *connections[id]
      fmt.Println("From server: removing parent")
      delete(connections, id) // Remove parent from list of connections
      fmt.Println("From server: parent removed")
    } else {
      fmt.Println("From server: Received message from a neighbor, not a parent.")
      // marco la ip com que ja he rebut la primera wave daquest neig
      connections[id].firstWave = true
    }
    if allConnections(true) {
      // envio missatge al parent
      if !initial {
        fmt.Println("From server: all neighbors have sent me a hello message")
        ch <- "parent;hello"
      } else {
        fmt.Println("From server: I'm init and all neighbors have sent me a hello message")
        fmt.Println("Decision event")
        ch <- "finish"
      }
      secondWave = true
    }

  } else if message == "finish"{
    if secondWave && !initial {
      fmt.Println("From server FINISH: secondWave & !initial")
      if len(connections) > 1 { // This means that it has neighbours and not only a parent
        fmt.Println("From server FINISH: has neighbors, sending message to neighbors")
         ch <- message
      }
      secondWave = false
    }
    if len(connections) > 0 {
      connections[id].secondWave = true
    }
    if allConnections(false) {
      // envio missatge al parent
      if !initial {
        fmt.Println("From server: all neighbors have sent me a finish message")
        ch <- "parent;finish"
        time.Sleep(time.Second/4)
        os.Exit(0)
      } else {
        fmt.Println("From server: I'm init and all neighbors have sent me a finish message")
        fmt.Println("I'm going to finish.")
        os.Exit(0)
      }
    }

  }
  /*if strings.ToLower(message) == "stop"{
    ip := strings.TrimSpace(string(temp[1]))
    fmt.Println("client " + id + " left..")
    ch <- ip
    conn.Close()
    return
  } else { */
  //fmt.Println("\n--- Message Received from " + id + ": " + string(message))
  // sample process for string received
  newmessage := strings.ToUpper(message)
  // send new string back to client
  conn.Write([]byte(newmessage + "\n"))
  serverConnection(conn, ch)
  //}
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
func getConfig(myFile string) (string, string, string, bool, []string){
  file, err := os.Open(myFile)
  if err != nil {
      fmt.Println(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  first := true
  serverIP, serverPort, id := "", "", ""
  init := false
  ips := []string{}
  for scanner.Scan() {
      if first {
        //serverIP = scanner.Text() // The first line is the server ip
        texts := strings.Split(scanner.Text(), ":")
        if len(texts) == 4 {
          fmt.Println("This is an initial node")
          init = true
        }
        serverIP = texts[0]
        serverPort = texts[1]
        id = texts[2]
        first = false
      } else {
        ips = append(ips, scanner.Text()) // Saving every client's ip
      }
  }
  return serverIP, serverPort, id, init, ips
}


/*
  Defines a connection between this machine and another one.
*/
type Connection struct {
    Channel chan string
    Connection net.Conn
    ID string // ID of the connected machine
    firstWave bool
    secondWave bool
}

/*
  Connects with a server with the IP "ip" (which already contains the port).
*/
func clients(ips []string, ch chan string, myIP string) {
    // Connect with the given ips and save their connections
    // and channels in a map.
    //connections := map[string]Connection{}
    for _, ip := range ips {
      chIP := make(chan string)
      c := Connection{chIP, connect(ip), ip, false, false}
      connections[ip] = &c
      fmt.Println("Connected ip:" + ip)
    }

    // We handle every client (connected ip) concurrently.
    for ip, c := range connections{
        go handleClient(*c, myIP, ip)
    }

    // Checking whether an "stop" message has arrived
    //go waitForStop(ch, connections)


    // Reading messages from terminal and sending them to the corresponding servers.
    for {
        // read in input from stdin
        /*reader := bufio.NewReader(os.Stdin)
        fmt.Print("Text to send: ")
        text, _ := reader.ReadString('\n')*/
        if initial && firstWave {
          for _, c := range connections {
            c.Channel <- "hello"
          }
          firstWave = false
        }

        message, _ := <- ch
        messages := strings.Split(message, ";")
        if messages[0] == "parent" {
          fmt.Println("Sending message to parent: " + messages[1])
          parent.Channel <- messages[1]

        } else {
          // Send to connections
          fmt.Println("Sending message to neighbours: " + messages[0])
          for _, c := range connections {
            c.Channel <- messages[0]
          }
        }

        /*// If the input message is "stop", we close the program.
        temp = strings.ToLower(temp)
        if temp == "stop"{
          //time.Sleep(time.Second/4)
          fmt.Println("Program has been stopped")
          os.Exit(1)
        }*/
      }
}


/*
  Send the written message to its corresponding server.
  c = connection with a client.
  myIP = the IP corresponding to the server of this program (we use it as an id)
*/
func handleClient(c Connection, myIP string, ip string){
    fmt.Println("Handle Client: " + ip + " myIP: " + myIP)
    for {
      text, _ := <- c.Channel
      fmt.Println("HANDLE CLIENT: Message received: " + text)
      conn := c.Connection
      temp := strings.TrimSpace(string(text)) + ";" + myIP
      fmt.Fprintf(conn, temp + "\n")
      fmt.Println("HANDLE CLIENT: Sent to connection")
      // listen for reply
      message, _ := bufio.NewReader(conn).ReadString('\n')
      //if err != nil { break }
      fmt.Println("HANDLE CLIENT: Message from server " + ip + ": " + message)
    }
}

var firstWave bool = true
var secondWave bool = false
var connections = map[string]*Connection{}
var parent Connection
var initial bool = false

func main() {
    if len(os.Args) != 2{
      usage()
    } else {
      var myIP, port, id string
      var ips []string
      myIP, port, id, initial, ips = getConfig(os.Args[1])
      fmt.Println(myIP, port, id, initial, ips)
      /*port := strings.Split(myIP, ":")[1] */
      ch := make(chan string)
      go server(port, ch)
      clients(ips, ch, myIP + ":" + port)

  }
}
