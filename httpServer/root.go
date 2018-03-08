package httpServer

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

// CRLF - Carriage Return, Line Feed "\r\n"
const CRLF = "\r\n"




type Server struct {
	Port int
	Router *Router
}

// Run - Run the server
func (server *Server) Run() {
	fmt.Println("Server opened on port " + strconv.Itoa(server.Port) + ".")

	// listen on all interfaces
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(server.Port))

	if (err != nil) {
		fmt.Println(err)
		return
	}


	// run loop forever (or more likely until ctrl-c)
	for {
		// wait for a connection
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println(err)
			break
		}

		go server.HandleConnection(conn)
	}
}


// HandleConnection - Handle a connection
func (server *Server) HandleConnection(conn net.Conn) {

	// Close connection when this function ends
	defer func() {
		//fmt.Println("Closing connection...")
		conn.Close()
	}()

	// output message received
	//fmt.Println("Message Received")
	connection := &Connection{
		netConn: conn,
		router: server.Router,
	}

	connection.Handle()
}


// PrettyPrint - make a map nice
func PrettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	println(string(b))
}
