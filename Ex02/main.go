//使用net.Listener开始fasthttp Serve的例子
//例如可以是UNIX Socket or TLS listener
package main

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"net"
	"log"
)

func main() {

	// Create network listener for accepting incoming requests.
	//
	// Note that you are not limited by TCP listener - arbitrary
	// net.Listener may be used by the server.
	// For example, unix socket listener or TLS listener.
	ln, err := net.Listen("tcp4", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("error in net.Listen: %s", err)
	}

	// This function will be called by the server for each incoming request.
	//
	// RequestCtx provides a lot of functionality related to http request
	// processing. See RequestCtx docs for details.
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "Hello, world! Requested path is %q", ctx.Path())
	}

	// Start the server with default settings.
	// Create Server instance for adjusting server settings.
	//
	// Serve returns on ln.Close() or error, so usually it blocks forever.
	if err := fasthttp.Serve(ln, requestHandler); err != nil {
		log.Fatalf("error in Serve: %s", err)
	}
}