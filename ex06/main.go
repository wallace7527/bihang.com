package main

//这是一个电报机器人的例子

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"fmt"
	"os"
	"net/http"
	"golang.org/x/net/proxy"
)

func main() {

	var replies map[string]string
	replies = make(map[string]string)

	replies["France"] = "Paris"
	replies["Italy"] = "Rome"
	replies["Japan"] = "Tokyo"
	replies["India"] = "New Delhi"


	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}

	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial

	bot, err := tgbotapi.NewBotAPIWithClient("585719882:AAGL9jrU__-fNSjiyfr03D5wMJf9BKDhV48",httpClient)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		reply, ok := replies[update.Message.Text]
		if ok {
			msg.Text = reply
		}
		bot.Send(msg)
	}
}


