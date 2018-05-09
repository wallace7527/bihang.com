package main

import (
	"gopkg.in/gcfg.v1"
	"path/filepath"
	"os"
	"strings"
)

var config Config

type (
	Config struct{
		Mysql struct{
			Username string
			Password string
			Host string
			Port string
			Database string
		}

		EtherscanApi struct {
			Apikey string
		}
	}
)

func iniFileName() string {
	exePath := os.Args[0]
	base := filepath.Base(exePath)
	suffix := filepath.Ext(exePath)
	return strings.TrimSuffix(base, suffix) + ".ini"
}

func readConfig(){
	gcfg.ReadFileInto(&config, iniFileName())
}