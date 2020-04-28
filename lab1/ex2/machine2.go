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
  fmt.Println("Launching server in port: " + port + "...")
  // listen on all interfaces,
  ln, _ := net.Listen("tcp", ":"+port)

  // run loop forever (or until ctrl-c)
  for {
    // accept connection on port
    conn, _ := ln.Accept()
    go serverConnection(conn)
  }
}


func serverConnection(conn net.Conn){
  // will listen for message to process ending in newline (\n)
  message, err := bufio.NewReader(conn).ReadString('\n')
  if err != nil {
        fmt.Println("client left..")
        conn.Close()
        return
    }
  //if err != nil { break }
  // output message received
  if strings.ToLower(message) == "stop\n"{
    //os.Exit(1)
    conn.Close()
    return
  } else {
    fmt.Print("--- Message Received:", string(message))
    // sample process for string received
    newmessage := strings.ToUpper(message)
    // send new string back to client
    conn.Write([]byte(newmessage + "\n"))
    serverConnection(conn)
  }
}

func connect(ip string) net.Conn{
  fmt.Print("Checking if server " + ip + " has started" + "\n")
  conn, err := net.Dial("tcp", ip)
  // Handle error: try to connect again until server is up and running
  if err != nil {
      fmt.Print("Error: ", err)
      fmt.Print("\n")
      time.Sleep(time.Second)
      return connect(ip)
  }  else {
      fmt.Print("Server Up and Running" + "\n")
      return conn
  }
}


/*
  Connects with a server with the IP "ip" (which already contains the port).
*/
func client(ips []string){
  connections := make([]net.Conn, len(ips))
  for i, ip := range ips {
    connections[i] = connect(ip)
  }

  for {
      // read in input from stdin
      reader := bufio.NewReader(os.Stdin)
      fmt.Print("Text to send: ")
      text, _ := reader.ReadString('\n')

      for _, conn := range connections {
        // send to socket
        fmt.Fprintf(conn, text + "\n")
        // listen for reply
        message, err := bufio.NewReader(conn).ReadString('\n')
        if err != nil { break }
        fmt.Print("Message from server: " + message)
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
        serverIP = strings.Split(serverIP, ":")[1]
        first = false
      } else {
        ips = append(ips, scanner.Text())
      }
  }
  return serverIP, ips

}



func main() {
    if len(os.Args) != 2{
      usage()
    } else {
      //port := ":6001" // Port used for creating our own server
      //ip := "127.0.0.1:6002" // Ip used for connecting to the other server
      port, ips := getConfig(os.Args[1])
      go server(port)
      //length := len(ips) - 1
      client(ips)
      //client(ip)
  }
}
