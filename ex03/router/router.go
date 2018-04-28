package router

import (
	"github.com/buaazp/fasthttprouter"
)

func NewRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()
	for _, r := range rs {
		handler := Logger(r.handle,r.name)
		router.Handle( r.method, r.path, handler )
	}
	return router
}