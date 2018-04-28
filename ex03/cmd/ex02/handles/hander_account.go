package handles

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"bihang.com/ex03/router"
)

func init() {
	router.Add("getaccount","/account/:uid","GET", getaccount)
	router.Add("addaccount","/account/","POST",	addaccount)
	router.Add("deleteaccount","/account/:uid","DELETE",deleteaccount)
	router.Add("modifyaccount","/account/:uid","PUT",modifyaccount)
	router.Add("getallaccount","/account","GET",getallaccount)
}

func getallaccount(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "you are get all account.")
}

func getaccount(ctx *fasthttp.RequestCtx) {
	uid := ctx.UserValue("uid")
	fmt.Fprintf(ctx, "you are get account %s", uid)
}
func modifyaccount(ctx *fasthttp.RequestCtx) {
	uid := ctx.UserValue("uid")
	fmt.Fprintf(ctx, "you are modify account %s", uid)
}
func deleteaccount(ctx *fasthttp.RequestCtx) {
	uid := ctx.UserValue("uid")
	fmt.Fprintf(ctx, "you are delete account %s", uid)
}
func addaccount(ctx *fasthttp.RequestCtx) {
	uid := ctx.FormValue("uid")
	fmt.Fprintf(ctx, "you are add account %s", string(uid))
}

