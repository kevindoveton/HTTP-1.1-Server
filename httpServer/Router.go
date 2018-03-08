package httpServer

type routeFn func(req *Request, res *Response)
//type Route struct{
//  routePath string
//  routeFunc routeFn
//}

type Router struct {
  registeredRoutes map[string]routeFn
}

func (router *Router) Init() {
  router.registeredRoutes = make(map[string]routeFn)
}

func (router *Router) AddRoute(path string, f routeFn) {
  router.registeredRoutes[path] = f
  return
}

func (router *Router) GetResponse(req *Request, res *Response) {
  if p := router.registeredRoutes[req.Get("Path")]; p != nil {
    p(req, res)
  } else if p := router.registeredRoutes["*"]; p != nil {
    router.registeredRoutes["*"](req, res)
  } else {
    res.Send404()
  }
}
