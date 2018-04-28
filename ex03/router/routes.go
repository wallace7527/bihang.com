package router

import (
	"github.com/valyala/fasthttp"
)

type route struct {
	name string
	path string
	method string
	handle fasthttp.RequestHandler
}

type routes []route
var rs routes

func Add(name string, path string, method string, handle fasthttp.RequestHandler) {
	rs = append( rs, route{name, path, method, handle} )
}