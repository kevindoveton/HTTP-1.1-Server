# HTTP/1.1 Server
A HTTP/1.1 server implemented in Go, there are definitely far better implementations than this, but I wanted to 
have a crack at implementing a HTTP server from the original spec using only inbuilt libraries. This is also my 
first attempt at using go in a project.  

Current tests performed using apache bench on 1GB or ram, with 1vCPU at 3GHz running on Ubuntu Server can handle 
between 2500 requests per second and 3500 requests per second sending a static page from file. I believe that 
it could handle significantly more RPS if a string was sent back, or if more ram was available as we were running 
out of ram. 

## Demo Usage
```go
package main

import (
  "github.com/kevindoveton/httpServer/httpServer"
)

func main() {

  router := &httpServer.Router{}
  router.Init()

  // example of SendString
  router.AddRoute("/", func(req *httpServer.Request, res *httpServer.Response) {
    res.SendString("Hello, World!")
  })


  // example of SendFile
  router.AddRoute("/file", func(req *httpServer.Request, res *httpServer.Response) {
    res.SendFile("docs/HelloWorld.html")
  })

  // example of static directory
  router.AddStaticRoute("/static", "docs")

  // example of catch all - 404 in this case
  router.AddRoute("*", func(req *httpServer.Request, res *httpServer.Response) {
    res.SetStatusCode(404)
    res.SendString("Oh no! The page can't be found!")
  })

  // start the server
  server := &httpServer.Server{
    Port:    8081,
    Router:  router,
  }

  server.Run()
}
```
