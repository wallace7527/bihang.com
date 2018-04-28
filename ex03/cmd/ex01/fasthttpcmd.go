package main

import (
	"gopkg.in/urfave/cli.v1"
)

var (
	FastHttpCompressFlag = cli.BoolFlag{
		Name:  "compress",
		Usage: "Whether to enable transparent response compression",
	}

	FastHttpAddrFlag = cli.StringFlag{
		Name:  "addr",
		Usage: "TCP address to listen to",
		Value: ":8080",
	}

)
