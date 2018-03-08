package httpServer

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

// CRLF - Carriage Return, Line Feed "\r\n"
const CRLF = "\r\n"



// SetWebRoot - set the web root dir
func (server *Server) SetWebRoot(path string) {
	server.WebRoot = path
}

// GetWebRoot - return set web root
func (server *Server) GetWebRoot() string {
	return server.WebRoot
}



type Server struct {
	WebRoot string
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
