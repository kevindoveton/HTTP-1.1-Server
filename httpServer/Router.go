package httpServer

import (
  "strings"
)

type routeFn func(req *Request, res *Response)

type Router struct {
  registeredRoutes map[string]routeFn
  //staticRoutes
}

func (router *Router) Init() {
  router.registeredRoutes = make(map[string]routeFn)
}

func (router *Router) AddRoute(path string, f routeFn) {
  router.registeredRoutes[path] = f
  return
}

func (router *Router) AddStaticRoute(path string, dir string) {
  router.registeredRoutes[path] = func(req *Request, res *Response) {
    // search for a static route
    parts := strings.Split(path, "/")
    url := ""
    for  _, p := range parts[1:] {
      url += "/" + p
      if p:= router.registeredRoutes[url]; p!= nil {
        res.SendStatic(url, dir, req.Get("Path"))
      }
    }
  }
}

func (router *Router) GetResponse(req *Request, res *Response) {
  path := req.Get("Path")

  // check exact route
  if p := router.registeredRoutes[path]; p != nil {
    p(req, res)
  } else {
    // search for a static route
    parts := strings.Split(path, "/")
    url := ""
    for  _, p := range parts[1:] {
      url += "/" + p
      if p:= router.registeredRoutes[url]; p!= nil {
        p(req, res)
        return
      }
    }

    // check the 404 route
    if p := router.registeredRoutes["*"]; p != nil {
      router.registeredRoutes["*"](req, res)
    } else {
      // we failed
      res.Send404()
    }
  }
}
