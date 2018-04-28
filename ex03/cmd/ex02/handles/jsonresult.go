package handles

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
	"time"
)


func jsonResult(ctx *fasthttp.RequestCtx, state int32, msg interface{}) {
	json.NewEncoder(ctx).Encode(DataResult{JsonTime(time.Now()),state,msg})
}

func JsonMsgResult(ctx *fasthttp.RequestCtx, msg interface{} ) {
	jsonResult(ctx, ERROR_SUCCESS, msg)
}

func JsonSuccResult(ctx *fasthttp.RequestCtx) {
	jsonResult(ctx, ERROR_SUCCESS, "success")
}

func JsonErrorResult(ctx *fasthttp.RequestCtx, state int32, msg string) {
	jsonResult(ctx, state, msg)
}
