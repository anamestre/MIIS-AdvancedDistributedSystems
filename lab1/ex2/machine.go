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
  //temp := strings.TrimSpace(string(message))
  temp := strings.Split(string(message), ";")
  if strings.ToLower(temp[0]) == "stop"{
    ip := strings.TrimSpace(string(temp[1]))
    fmt.Println("client left.." + ip)
    conn.Close()
    ch <- ip
    return
  } else {
    fmt.Print("--- Message Received:", string(message))
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

func waitForStop(ch chan string, connections map[string]net.Conn){
    elem, _ := <- ch
    delete(connections, elem)
}


/*
  Connects with a server with the IP "ip" (which already contains the port).
*/
func client(ips []string, ch chan string, myIP string){
  connections := map[string]net.Conn{}
  for _, ip := range ips {
    connections[ip] = connect(ip)
  }

  for {
      // read in input from stdin
      reader := bufio.NewReader(os.Stdin)
      fmt.Print("Text to send: ")
      text, _ := reader.ReadString('\n')
      go waitForStop(ch, connections)
      //fmt.Print("Client loop, connections length: ")
      //fmt.Println(len(connections))
      if len(connections) == 0 {
        //fmt.Println("Len connections == 0")
        os.Exit(1)
      }
      temp := strings.TrimSpace(string(text))
      temp = strings.ToLower(temp)

      for _, conn := range connections {
        if temp == "stop"{
          fmt.Println("Client ip " + myIP + " senting stop")
          text = "stop;" + myIP
        }
          fmt.Fprintf(conn, text + "\n")
          // listen for reply
          message, err := bufio.NewReader(conn).ReadString('\n')
          if err != nil { break }
          fmt.Print("Message from server: " + message)
      }
      if temp == "stop"{
        fmt.Println("Client: Exit due to stop")
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
        serverIP = scanner.Text()
        first = false
      } else {
        ips = append(ips, scanner.Text())
      }
  }
  return serverIP, ips
}

/*

type Connection struct {
	IP string
	Port string
  Ch chan string
}


func clients(connections){
    for {
      fmt.Print("Text to send: ")
      text, _ := reader.ReadString('\n')
      for every connection {
          envia al seu channel el text
      }
    }
} */

func clients(ips []string, ch chan string, myIP string){
  connections := map[string]net.Conn{}
  for _, ip := range ips {
    connections[ip] = connect(ip)
  }
  for {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Text to send: ")
    text, _ := reader.ReadString('\n')
    go waitForStop(ch, connections)
    //fmt.Print("Client loop, connections length: ")
    //fmt.Println(len(connections))
    if len(connections) == 0 {
      //fmt.Println("Len connections == 0")
      os.Exit(1)
    }
    for ip, conn := range connections {
      go handleClient(conn, text, ip)
    }
    temp := strings.TrimSpace(string(text))
    temp = strings.ToLower(temp)
    if temp == "stop"{
      fmt.Println("Client: Exit due to stop")
      os.Exit(1)
    }
  }
}


func handleClient(conn net.Conn, text string, myIP string){
    temp := strings.TrimSpace(string(text))
    temp = strings.ToLower(temp)
    if temp == "stop"{
      fmt.Println("Client ip " + myIP + " senting stop")
      text = "stop;" + myIP
    }
    fmt.Fprintf(conn, text + "\n")
    // listen for reply
    message, _ := bufio.NewReader(conn).ReadString('\n')
    //if err != nil { break }
    fmt.Print("Message from server: " + message)
    /*for {
      //llegeixo del meu channel
      //envio missatge al server al qual estic connectat
    }*/
  }

func main() {
    if len(os.Args) != 2{
      usage()
    } else {
      myIP, ips := getConfig(os.Args[1])
      port := strings.Split(myIP, ":")[1]
      ch := make(chan string)
      go server(port, ch)
      //client(ips, ch, myIP)
      client(ips, ch, myIP)
  }
}
