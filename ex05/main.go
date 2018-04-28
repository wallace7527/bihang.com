package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)


//这个会以太坊etherscan Websocket接口，用于监控某个账号的交易事件。
var addr = flag.String("addr", "socket.etherscan.io", "http service address")

type EtherscanWS struct{
	url url.URL
}

func NewEtherscanWS() *EtherscanWS{
	return &EtherscanWS{
		url:url.URL{Scheme: "wss", Host: "socket.etherscan.io", Path: "/wshandler"},
	}
}

func (e * EtherscanWS)start() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Printf("connecting to %s", e.url.String())

	//连接
	c, _, err := websocket.DefaultDialer.Dial(e.url.String(), nil)
	if err != nil {
		log.Println("dial:", err)
		return err
	}
	defer c.Close()

	done := make(chan struct{})

	//启动接受历程
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	pingTicker := time.NewTicker(time.Second*20)
	defer pingTicker.Stop()

	//订阅交易记录
	err = c.WriteMessage(websocket.TextMessage, []byte("{\"event\": \"txlist\", \"address\": \"0x2a65aca4d5fc5b5c859090a6c34d164135398226\"}"))
	if err != nil {
		log.Println("subscribe write:", err)
		//订阅错误输出不退出
	}

	//发送Ping直到连接断开
	//或者用户中断
	for {
		select {
		case <-done:
			return errors.New("Restart..")
		case <-pingTicker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte("{\"event\":\"ping\"}"))
			if err != nil {
				log.Println("ping write:", err)
				return err
			}
		case <-interrupt: // 用户中断自动停止
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}

func (e * EtherscanWS)Run() {
	for {
		if err := e.start(); err == nil {
			break
		}

		log.Println("The network is disconnected and reconnected after 10 seconds.")
		//等一段时间
		//重新连接
		time.Sleep(time.Second*10)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	ews := NewEtherscanWS()
	ews.Run()
}