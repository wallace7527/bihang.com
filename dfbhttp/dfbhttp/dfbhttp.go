package dfbhttp

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"net"
)


// RequestHandler 类型，使用 RequestCtx 传递 HTTP 的数据
func requestHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world! Requested path is %q", ctx.Path())
}

var (
	lnObject net.Listener
)

func fastHttpStart() {
	// Create network listener for accepting incoming requests.
	//
	// Note that you are not limited by TCP listener - arbitrary
	// net.Listener may be used by the server.
	// For example, unix socket listener or TLS listener.
	ln, err := net.Listen("tcp4", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("error in net.Listen: %s", err)
	}

	lnObject = ln;

	// Start the server with default settings.
	// Create Server instance for adjusting server settings.
	//
	// Serve returns on ln.Close() or error, so usually it blocks forever.
	if err := fasthttp.Serve(ln, requestHandler); err != nil {
		log.Fatalf("Start DFCHttp fail: %s", err)
	}
}

func DFCHttpStart() {
	go fastHttpStart()
}

func DFCHttpStop() {
	log.Println("entery DFCHttpStop")
	if lnObject != nil {
		err := lnObject.Close()
		if err != nil {
			log.Fatalf("Stop DFCHttp fail: %s", err)
		}
	}else
	{
		log.Println("lnObject is nil.");
	}


}