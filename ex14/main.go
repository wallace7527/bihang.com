package main

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"encoding/json"
	"github.com/wanglei-ok/evms"
)

type Result struct {
	Pass bool			`json:"pass"`
	Message string		`json:"message"`
	Id int				`json:"id"`
}

// RequestHandler 类型，使用 RequestCtx 传递 HTTP 的数据
func httpHandle(ctx *fasthttp.RequestCtx) {
	addr := string(ctx.FormValue("addr"))
	sig := string(ctx.FormValue("sig"))
	msg := string(ctx.FormValue("msg"))

	err, id := evms.VerifyMessage(addr, sig, msg)
	var result Result

	result.Id = id
	if err != nil {
		result.Pass = false
		result.Message = fmt.Sprintf("%v", err)
	}else {
		result.Pass = true
		result.Message = fmt.Sprintf("Message Signature Verified, Pass(%d)", id)
	}

	json.NewEncoder(ctx).Encode(result)
}

func main() {
	// 一定要写 httpHandle，否则会有 nil pointer 的错误，没有处理 HTTP 数据的函数
	if err := fasthttp.ListenAndServe("localhost:8080", httpHandle); err != nil {
		fmt.Println("start fasthttp fail:", err.Error())
	}
}