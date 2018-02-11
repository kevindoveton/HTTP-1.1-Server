package main

import (
  "fmt"
	"net"
	"bufio"
	// "strings"
)

func main() {
  fmt.Println("Server opened on port 8081.")
	// listen on all interfaces
  ln, _ := net.Listen("tcp", ":8081")

  // run loop forever (or until ctrl-c)
  for {
    // wait for a connection
    conn, _ := ln.Accept()

    // will listen for message to process ending in newline (\n)
    message, _ := bufio.NewReader(conn).ReadString('\n')
    
    // output message received
    fmt.Print("Message Received:", string(message))
    
    // sample response
    newmessage := "HTTP/1.1 200 OK\nContent-Length: 5\r\n\r\nhello\n"
    
    // send new string back to client
    conn.Write([]byte(newmessage))
    
    // we're done, close the connection
    conn.Close()
  }
}
