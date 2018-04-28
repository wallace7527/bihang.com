package main

import (
	"bihang.com/ex03/cmd/utils"
	"os"
	"fmt"

	"sort"
	"gopkg.in/urfave/cli.v1"
	"bihang.com/ex03/log"
	"github.com/valyala/fasthttp"
)

const (
	clientIdentifier = "cmd01"
)

var (
	app = utils.NewApp("Example command line interface")

	rpcFlags = []cli.Flag{
		utils.RPCEnabledFlag,
		utils.RPCListenAddrFlag,
		utils.RPCPortFlag,
		utils.RPCApiFlag,
	}

	fastHttpFlags = []cli.Flag{
		FastHttpCompressFlag,
		FastHttpAddrFlag,
	}
)

func init() {
	app.Action = ex01
	app.HideVersion = true
	app.Copyright = "Copyright 2018 The cmd01 Authors"

	app.Commands = []cli.Command{
		versionCommand,
		licenseCommand,
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	app.Flags = append(app.Flags, rpcFlags...)
	app.Flags = append(app.Flags, fastHttpFlags...)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ex01( ctx *cli.Context) error {
	log.Root().SetHandler(log.StdoutHandler)
	log.Info("Program starting", "args", os.Args)

	h := requestHandler
	if ctx.GlobalBool(FastHttpCompressFlag.Name) {
		log.Info("using compress")
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(ctx.GlobalString(FastHttpAddrFlag.Name), h); err != nil {
		log.Info("Error in ListenAndServe: ", "error", err)
	}

	return nil
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world!\n\n")

	fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
	fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
	fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
	fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
	fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
	fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
	fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
	fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
	fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
	fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())

	fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)

	ctx.SetContentType("text/plain; charset=utf8")

	// Set arbitrary headers
	ctx.Response.Header.Set("X-My-Header", "my-header-value")

	// Set cookies
	var c fasthttp.Cookie
	c.SetKey("cookie-name")
	c.SetValue("cookie-value")
	ctx.Response.Header.SetCookie(&c)
}
