package main

import ("net"
        "fmt"
        "bufio"
        "os"
        "time"
        "strings")


/*
  Turns up a server with the port "port"
*/
func turnServerUp(port string) (net.Listener){
  fmt.Println("--- Launching server in port: " + port + "...")
  // listen on all interfaces,
  ln, _ := net.Listen("tcp", ":" + port)
  return ln
}

/*
  Creates a connection with all neighbours from the config file and stores them
  in "connections"
*/
func connectClients(ips []string){
  for _, ip := range ips {
    chIP := make(chan string)
    c := Connection{chIP, connect(ip), false, false, ip}
    connections[ip] = &c
  }
}


/*
  Creates a server using the Port "port".
*/
func server(ln net.Listener, ch chan string) {
  // run loop forever (or until ctrl-c)
  for {
    // accept connection on port
    conn, _ := ln.Accept()
    go serverConnection(conn, ch)
  }
}

/*
  Checks whether all connections have finished the corresponding wave or not.
  First marks if this is the first or the second wave.
*/
func allConnections(first bool) bool {
  for _, conn := range connections {
      if first {
        if ! conn.firstWave {return false}
      } else {
        if ! conn.secondWave {return false}
      }
  }
  return true
}


func deepCopy(conn Connection) Connection{
  newConn := Connection{conn.Channel,
                        conn.Connection,
                        conn.firstWave,
                        conn.secondWave,
                        conn.ip}
  return newConn
}

func handleWave(ip, id string, ch chan string){
  fmt.Println("Received first wave from " + ip + " with id: " + id)
  // If this is the first time the node gets a wave
  if waveID == "" {
    waveID = id
    // Send to neighbours
    fmt.Println("1. My parent is " + id)
    parent = *connections[ip]
    delete(connections, ip) // Remove parent from list of neighbours
    fmt.Println("1. My neighbours are:")
    for i, c := range connections {
      fmt.Println(" 1. Neighbour: " + i + " " + c.ip)
    }
    if len(connections) > 0 { // This means that it has neighbours and not only a parent
       ch <- "wave;" + id // Send wave to neighbours
    }
  } else {
    // If this wave has a larger ID than the wave that has already hit me
    if id > waveID {
      // Mark this new node as a parent
      waveID = id
      fmt.Println("2. My parent is " + id)
      fmt.Println("2. My current parent's IP is " + parent.ip)
      fmt.Println("2. My IP is " + myIP)
      if parent.ip != myIP {
        fmt.Println("My IP is not the same as my parent's")
        oldParent := deepCopy(parent)
        connections[parent.ip] = &oldParent // Afegeixo l'antic pare com a neighbour
      }
      parent = deepCopy(*connections[ip])
      //parent = *connections[ip] // Assigno el nou parent
      delete(connections, ip) // Remove parent from list of neighbours
      fmt.Println("2. My neighbours are:")
      for i, c := range connections {
        fmt.Println(" 2. Neighbour: " + i + " " + c.ip)
      }
      if len(connections) > 0 { // This means that it has neighbours and not only a parent
         ch <- "wave;" + id // Send wave to neighbours
         // Treure el rebut de tots els neighbours
         for _, c := range connections {
           c.firstWave = false
         }
      }
    } else if id == waveID{
      // Marco com a rebut
      connections[ip].firstWave = true
    }
    // Miro si he rebut de tots els meus neighbors aquest missatge
    // Si es que si: Envio al meu parent
    if allConnections(true) {
      // If this is not an initial node, send message to parent
      fmt.Println("All connections are true")
      if myIP == parent.ip {
        fmt.Println("Decision event")
        ch <- "finish"
      } else {
        ch <- "parent;wave;" + id
      }

      secondWave = true
    }
  }
}

/*
  Handles messages received from the second wave.
  If message from parent has arrived, sends message to neighbours.
  If all messages from neighbours have arrived, sends message to parent.
  If parent has received all messages from neighbours, finishes.
*/
func handleSecondWave(ip, id string, ch chan string){
  fmt.Println("Received second wave from " + ip)
  if secondWave && !initial {
    if len(connections) > 0 { // This means that it has neighbours and not only a parent
       ch <- "finish"
    }
    secondWave = false
  }
  if len(connections) > 0 {
    if _, ok := connections[ip]; ok {
        connections[ip].secondWave = true
    }
  }
  if allConnections(false) {
    // If this is not an initial node, send message to parent
    if !initial {
      ch <- "parent;finish"
      time.Sleep(time.Second/4)
    }
    fmt.Println("I'm going to finish.")
    os.Exit(0)
  }
}


/*
  Handles a connection (conn) for the server.
*/
func serverConnection(conn net.Conn, ch chan string){
  message, _ := bufio.NewReader(conn).ReadString('\n')
  // output message received
  text := strings.TrimSpace(string(message))
  temp := strings.Split(string(text), ";")
  if len(temp) == 3 {
    message = temp[0]
    id := temp[1]
    ip := temp[2]
    if message == "wave"{
      handleWave(ip, id, ch)
    } else if message == "finish"{
      handleSecondWave(ip, id, ch)
    }
  }
  serverConnection(conn, ch)
}


/*
  Makes the first connections with the corresponding servers.
  Doesn't stop until a connection is made.
*/
func connect(ip string) net.Conn{
  fmt.Print("--- Checking if server " + ip + " has started" + "\n")
  conn, err := net.Dial("tcp", ip)
  // Handle error: try to connect again until server is up and running
  if err != nil {
      fmt.Println("--- Error: ", err)
      time.Sleep(time.Second)
      return connect(ip)
  }  else {
      fmt.Print("----- Server " + ip + " Up and Running" + "\n")
      return conn
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
  Connects with a server with the IP "ip" (which already contains the port).
*/
func clients(ips []string, ch chan string, myIP string) {
    // We handle every client (connected ip) concurrently.
    for _, c := range connections{
        go handleClient(*c, myIP)
    }

    for {
        if initial && firstWave {
          waveID = myID
          parent = Connection{nil, nil, false, false, myIP}
          for _, c := range connections {
            //c.Channel <- "hello"
            c.Channel <- "wave;" + myID
          }
          fmt.Println("Init -> Sending first wave message to neighbours")
          firstWave = false
        }

        // This is the message that the server has to send
        message, _ := <- ch
        messages := strings.Split(message, ";")
        if messages[0] == "parent" {
          if messages[1] == "wave"{
            fmt.Println("Sending first wave message to parent")
          }
          parent.Channel <- strings.Join(messages[1:], ";")
        } else {
          // Send to connections
          if messages[0] == "wave"{
            fmt.Println("Sending first wave message to neighbours")
          }
          for _, c := range connections {
            c.Channel <- message
          }
        }
      }
}


/*
  Send the written message to its corresponding server.
  c = connection with a client.
  myIP = the IP corresponding to the server of this program (we use it as an id)
*/
func handleClient(c Connection, myIP string){
    for {
      text, _ := <- c.Channel
      conn := c.Connection
      temp := strings.TrimSpace(string(text)) + ";" + myIP
      fmt.Fprintf(conn, temp + "\n")
    }
}


/*
  Defines a connection between this machine and another one.
*/
type Connection struct {
    Channel chan string
    Connection net.Conn
    firstWave bool
    secondWave bool
    ip string
}


/*
  Initializations
*/
var firstWave bool = true
var secondWave bool = false
var connections = map[string]*Connection{} // Stores neighbours' connections
var parent Connection // stores parent's connection
var initial bool = false
var myID string
var waveID string = ""
var myIP string

func main() {
    if len(os.Args) != 2{
      usage()
    } else {
      var port string
      var ips []string
      var IP string
      IP, port, myID, initial, ips = getConfig(os.Args[1])
      fmt.Println("This is the server num " + port)
      myIP = IP + ":" + port

      ln := turnServerUp(port)
      connectClients(ips)

      ch := make(chan string)
      go server(ln, ch)
      clients(ips, ch, myIP)
  }
}
