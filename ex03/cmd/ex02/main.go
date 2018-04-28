package main

import (
	"log"
	"github.com/valyala/fasthttp"
	"bihang.com/ex03/router"
	_ "bihang.com/ex03/cmd/ex02/handles"
)



func main() {
	router := router.NewRouter()
	log.Fatal(fasthttp.ListenAndServe("0.0.0.0:8080", router.Handler))
}