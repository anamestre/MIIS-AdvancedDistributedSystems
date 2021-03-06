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
    c := Connection{chIP, connect(ip), false, false, ip, ""}
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
  First marks if this is the first (election) or the second (finish) wave.
*/
func allConnections(first bool, id string) bool {
  for ip, conn := range connections {
      if ip != parent.ip{
        if first {
          if ! conn.electionWave {return false}
          if id != conn.waveID {return false}
        } else {
          if ! conn.finishWave {return false}
        }
      }
  }
  return true
}


/*
Returns a deep copy of connection (a new position in memory)
*/
func deepCopy(conn Connection) Connection{
  newConn := Connection{conn.Channel,
                        conn.Connection,
                        conn.electionWave,
                        conn.finishWave,
                        conn.ip,
                        conn.waveID}
  return newConn
}

/*
  Handles messages received from the first wave (the election one).
  If this is the first election wave that hits this machine or it has a larger
  id than the one that had already hit the machine, save this new id. (waveID)
  This new id is going to be the parent of this machine. Then a wave to all neighbours
  is sent with this new id.
  If all neighbours have sent back the same id, send message to parent. But, if the
  machine has no parent, decides and sends a wave to finish processes on all machines.
*/
func handleWave(ipSender, idSender, idWave string, ch chan string){
  fmt.Println("Received election wave from " + ipSender + " with wave id: " + idWave)
  // If this is the first time the node gets a wave
  if waveID == "" {
    waveID = idWave
    fmt.Println("My parent is " + idSender)
    parent = *connections[ipSender]
    parent.waveID = idWave
    if len(connections) > 0 {
       ch <- "wave;" + idWave // Send wave to neighbours with wave id
    }
  } else {
    // If this new wave has a larger ID than the wave that has already hit me
    if idWave > waveID {
      fmt.Println("My parent is " + idSender)
      waveID = idWave
      parent = deepCopy(*connections[ipSender]) // Mark this new node as a parent
      parent.waveID = idWave
      if len(connections) > 0 {
         // Since the node is joining a new wave, we have to mark that
         // the other neighbours have not sent the node this wave.
         for _, c := range connections {
           c.electionWave = false
         }
         ch <- "wave;" + idWave // Send wave to neighbours
      }
    } else if idWave == waveID{
      connections[ipSender].electionWave = true // Mark this neighbour as received
      connections[ipSender].waveID = idWave
    }
  }
    // Checking whether all neighbours have sent me the same wave
    if allConnections(true, idWave) {
      // If this is an initial node and has no parent, it decides
      if myIP == parent.ip {
        fmt.Println("DECISION EVENT - I'M THE LEADER")
        ch <- "finish;" + idWave // sending finish wave
      } else { // If this node has a parent, send wave to parent
        ch <- "parent;wave;" + idWave
      }
      finishWave = true
    }
}

/*
  Handles messages received from the second wave.
  If message from parent has arrived, sends message to neighbours.
  If all messages from neighbours have arrived, sends message to parent.
  If parent has received all messages from neighbours, finishes.
*/
func handleFinishWave(ip, id string, ch chan string){
  fmt.Println("Received finish wave from " + ip)
  if finishWave && parent.ip != myIP {
    if len(connections) > 0 { // This means that it has neighbours and not only a parent
       ch <- "finish;" + myID
    }
    finishWave = false
  }
  if len(connections) > 0 {
    if _, ok := connections[ip]; ok {
        connections[ip].finishWave = true
    }
  }
  if allConnections(false, id) {
    // If this is not an initial node, send message to parent
    if parent.ip != myIP {
      ch <- "parent;finish;" + myID
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
  if len(temp) == 4 {
    message = temp[0]
    idWave := temp[1]
    ipSender := temp[2]
    idSender := temp[3]
    if message == "wave"{
      handleWave(ipSender, idSender, idWave, ch)
    } else if message == "finish"{
      handleFinishWave(ipSender, idWave, ch)
    }
  }
  newmessage := strings.ToUpper(message)
  conn.Write([]byte(newmessage + "\n"))
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
        if initial && electionWave {
          waveID = myID
          parent = Connection{nil, nil, false, false, myIP, waveID}
          for _, c := range connections {
            c.Channel <- "wave;" + myID
          }
          fmt.Println("Init -> Sending election wave message to neighbours")
          electionWave = false
        }

        // This is the message that the server has to send
        message, _ := <- ch
        messages := strings.Split(message, ";")
        if messages[0] == "parent" {
          if messages[1] == "wave"{
            fmt.Println("Sending election wave message to parent")
          } else {
            fmt.Println("Sending finish wave message to parent")
          }
          parent.Channel <- strings.Join(messages[1:], ";")
        } else {
          // Send to connections
          if messages[0] == "wave"{
            fmt.Println("Sending election wave message to neighbours")
          } else {
            fmt.Println("Sending finish wave message to neighbours")
          }
          for i, c := range connections {
            if i != parent.ip {
              c.Channel <- message
            }
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
      // A message is sent with the folloring structure:
      // "wave";ID of the wave;IP of the sender; ID of the sender
      temp := strings.TrimSpace(string(text)) + ";" + myIP + ";" + myID
      //ip := c.ip
      fmt.Fprintf(conn, temp + "\n")
      _, err := bufio.NewReader(conn).ReadString('\n')
      if err != nil { break }
      //fmt.Println("Message back from server " + ip + ": " + message)
    }
}


/*
  Defines a connection between this machine and another one.
*/
type Connection struct {
    Channel chan string // Channel for sending messages to client
    Connection net.Conn // Connection to communicate with the client
    electionWave bool // wether the machine has received an electionWave from this client
    finishWave bool // wether the machine has received a finish wave from this client
    ip string // ip of the client
    waveID string
}


/*
  Initializations
*/
var electionWave bool = true
var finishWave bool = false
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
