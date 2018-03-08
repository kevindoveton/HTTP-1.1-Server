# HTTP/1.1 Server
A HTTP/1.1 server implemented in Go, there are definitely far better implementations than this, but I wanted to have a crack at implementing a HTTP server from the original spec using only inbuilt libraries. This is also my first attempt at using go in a project.  

## Demo Usage
```go
package main

import (
  "github.com/kevindoveton/httpServer/httpServer"
)

func main() {

  router := &httpServer.Router{}
  router.Init()

  // home
  router.AddRoute("/", func(req *httpServer.Request, res *httpServer.Response) {
    res.SendString("Hello, World!")
  })

  router.AddRoute("/file", func(req *httpServer.Request, res *httpServer.Response) {
    res.SendFile("docs/HelloWorld.html")
  })

  // error page
  router.AddRoute("*", func(req *httpServer.Request, res *httpServer.Response) {
    res.SetStatusCode(404)
    res.SendString("Oh no! The page can't be found!")
  })

  server := &httpServer.Server{
    Port:    8081,
    Router:  router,
  }

  server.Run()
}
```