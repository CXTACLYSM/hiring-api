package routing

import (
	"log"
	"net/http"

	"github.com/CXTACLYSM/hiring-api/internal/middleware"
)

const ApiPrefix = "api"

const (
	V1 = "v1"
	V2 = "v2"
)

const (
	RouteInfo     = ""
	RouteRegister = "register"
	RouteLogin    = "login"
)

type HttpRoutable interface {
	http.Handler
	Middlewares() []middleware.Middleware
	RouteMetadata() *Metadata
}

type Router struct {
	routes []*Route
}

type Route struct {
	metadata    *Metadata
	middlewares []middleware.Middleware
	handler     http.Handler
}

type Metadata struct {
	prefix  string
	version string
	name    string
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Init(handlers []HttpRoutable) error {
	for _, handler := range handlers {
		route := newRoute()
		route.metadata = handler.RouteMetadata()
		route.middlewares = handler.Middlewares()
		route.handler = handler

		log.Printf("Router initializing: %s", route.String())

		r.addRoute(route)
	}

	return nil
}

func newRoute() *Route {
	return &Route{}
}

func NewRouteMetadata(prefix string, version string, name string) *Metadata {
	return &Metadata{
		prefix:  prefix,
		version: version,
		name:    name,
	}
}

func (r *Router) GetRoutes() []*Route {
	return r.routes
}

func (r *Router) addRoute(route *Route) {
	r.routes = append(r.routes, route)
}

func (route *Route) String() string {
	return "/" + route.metadata.prefix + "/" + route.metadata.version + "/" + route.metadata.name
}

func (route *Route) Handler() http.Handler {
	finalHandler := route.handler
	for i := len(route.middlewares) - 1; i >= 0; i-- {
		finalHandler = route.middlewares[i](finalHandler)
	}
	return finalHandler
}
