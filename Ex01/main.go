//直接使用ListenAndServe使用fasthttp的例子
package main

import (
	"fmt"
	"github.com/valyala/fasthttp"

)



type Element interface{}

type Queue interface {
	Offer(e Element)    //向队列中添加元素
	Poll()   Element    //移除队列中最前面的元素
	Clear()  bool       //清空队列
	Size()  int            //获取队列的元素个数
	IsEmpty() bool   //判断队列是否是空
}

type  sliceEntry struct{
	element []Element
}

func NewQueue() *sliceEntry {
	return &sliceEntry{}
}

//向队列中添加元素
func (entry *sliceEntry) Offer(e Element) {
	entry.element = append(entry.element,e)
}

//移除队列中最前面的额元素
func (entry *sliceEntry) Poll() Element {
	if entry.IsEmpty() {
		fmt.Println("queue is empty!")
		return nil
	}

	firstElement := entry.element[0]
	entry.element = entry.element[1:]
	return firstElement
}

func (entry *sliceEntry) Clear() bool {
	if entry.IsEmpty() {
		fmt.Println("queue is empty!")
		return false
	}
	for i:=0 ; i< entry.Size() ; i++ {
		entry.element[i] = nil
	}
	entry.element = nil
	return true
}

func (entry *sliceEntry) Size() int {
	return len(entry.element)
}

func (entry *sliceEntry) IsEmpty() bool {
	if len(entry.element) == 0 {
		return true
	}
	return false
}


var taskQueue = NewQueue()

// RequestHandler 类型，使用 RequestCtx 传递 HTTP 的数据
func httpHandle(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	if  path == "/get"	{
		deviceid := string(ctx.FormValue("deviceid"))
		devicemodel := string(ctx.FormValue("devicemodel"))
		fmt.Printf("deviceid=%s, devicemodel=%s\n",deviceid,devicemodel)
		fmt.Println("Has task:",!taskQueue.IsEmpty())

		if !taskQueue.IsEmpty() {
			task := taskQueue.Poll().(string)
			fmt.Fprintf(ctx, task) // *RequestCtx 实现了 io.Writer
			fmt.Println("Reponse:",task)
		}
	}else if path == "/set" {
		t := string(ctx.FormValue("task"))
		m := string(ctx.FormValue("market"))
		if len(m) <= 0 {
			m = "QQDownloader"
		}
		if len(t) > 0 {
			task := fmt.Sprintf("{\"task\":\"%s\",\"market\":\"%s\"}", t,m)
			taskQueue.Offer(task)
			fmt.Fprintf(ctx, task) // *RequestCtx 实现了 io.Writer
			fmt.Println("Reponse:",task)
		}
	}else if path == "/" {
		ctx.SendFile("home.html")
	}
}

func main() {
	// 一定要写 httpHandle，否则会有 nil pointer 的错误，没有处理 HTTP 数据的函数
	if err := fasthttp.ListenAndServe("0.0.0.0:12345", httpHandle); err != nil {
		fmt.Println("start fasthttp fail:", err.Error())
	}
}


