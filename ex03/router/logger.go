package router

import (
	"time"
	"log"
	"github.com/valyala/fasthttp"
)

func Logger(inner fasthttp.RequestHandler, name string) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		start := time.Now()

		inner(ctx)

		log.Printf(
			"%s\t%s\t%s\t%s",
			string(ctx.Method()),
			string(ctx.RequestURI()),
			name,
			time.Since(start),
		)
	})
}
