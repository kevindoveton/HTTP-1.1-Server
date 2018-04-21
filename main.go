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

  router.AddStaticRoute("/static", "docs")

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
