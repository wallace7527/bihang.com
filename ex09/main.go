//这个例子演示如何进行日志文件分割
package main

import (
	"log"
	"time"
)

func init(){
	logSetup()
}

func main(){
	for {
		log.Println("Helloworld!")
		time.Sleep(time.Second)
	}

}