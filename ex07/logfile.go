/*
Copyright 2018 The go-eam Authors
This file is part of the go-eam library.

logfile
设置标准日志的输出格式并输出到文件


wanglei.ok@foxmail.com

1.0
版本时间：2018年4月13日18:32:12

*/

package main

import (
	"os"
	"io"
	"path/filepath"
	"strings"
	"fmt"
	"time"
	"log"
)

//取得当前可执行程序路径
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//创建path执行的文件夹
func createDir(path string) bool {
	// check
	if _, err := os.Stat(path); err != nil {
		//不存在创建
		err := os.MkdirAll(path, 0711)

		if err != nil {
			return false
		}
	}
	return true
}

//取得当前可执行程序名称
func baseName() string {
	return filepath.Base(os.Args[0])
}

//在当前可执行程序路径下创建log文件夹存放日志
func createLogDir() string {
	//当前目录加log
	logPath := fmt.Sprintf("%s/log", getCurrentDirectory())
	if createDir(logPath) {
		return logPath
	}
	return  ""
}

//指定路径、前缀、后缀，返回格式化的日志文件路径
//格式化方式为 /path/to/file/<prefix>YYYYMMDD<suffix>
func logFilePath(path, prefix, suffix string) string {
	return fmt.Sprintf("%s/%s%s%s", path, prefix, time.Now().Format("20060102"), suffix)
}

//设置默认logger的
//标志位log.Ldate | log.Lmicroseconds
//设置日志输出到文件 ./log/baseName_YYYYMMDD.log
func logSetup() {
	log.SetFlags(log.Ldate | log.Lmicroseconds )

	//创建日志文件夹
	logDir := createLogDir()
	if len(logDir) != 0 {
		//baseName_YYYYMMDD.log
		filePath := logFilePath(logDir, baseName()+"_",".log")
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return
		}

		log.SetOutput(io.MultiWriter(file, os.Stderr))
	}
}


