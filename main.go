package main

import (
	"github.com/kevindoveton/httpServer/httpServer"
)

func main() {

	// set the path to the web root
	httpServer.SetWebRoot("/web")

	// run!
	httpServer.Run(8081)

}
