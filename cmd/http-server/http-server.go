package main

import (
	"fmt"
	"net"

	"github.com/kevindoveton/http-server/pkg/http-server"
)

func main() {
	fmt.Println("Server opened on port 8081.")
	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8081")

	// run loop forever (or until ctrl-c)
	for {
		// wait for a connection
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}
		go connection.HandleConnection(conn)
	}
}
