/*
Copyright 2018 The go-execsql Authors
This file is part of the go-execsql library.

config
读取配置文件


wanglei.ok@foxmail.com

1.0
版本时间：16:22 2018/5/9

*/

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
			DSN string
		}

		Action struct {
			UseTrans bool
			SQL []string
		}
	}
)

//获取配置文件的名称
//当前应用程序名后缀.ini
func iniFileName() string {
	exePath := os.Args[0]
	base := filepath.Base(exePath)
	suffix := filepath.Ext(exePath)
	return strings.TrimSuffix(base, suffix) + ".ini"
}


//从配置文件中读取配置
func readConfig() error {
	return gcfg.ReadFileInto(&config, iniFileName())
}