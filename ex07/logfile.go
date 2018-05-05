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

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

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


func baseName() string {
	return filepath.Base(os.Args[0])
}


func createLogDir() string {
	//当前目录加log
	logPath := fmt.Sprintf("%s/log", getCurrentDirectory())
	if createDir(logPath) {
		return logPath
	}
	return  ""
}


func logFilePath(path, prefix, suffix string) string {
	return fmt.Sprintf("%s/%s_%s%s", path, prefix, time.Now().Format("20060102"), suffix)
}

func logSetup() {
	log.SetFlags(log.Ldate | log.Lmicroseconds )

	//创建日志文件夹
	logDir := createLogDir()
	if len(logDir) != 0 {
		//format /path/to/file/progname_YYYYMMDD.log
		filePath := logFilePath(logDir, baseName(), ".log")
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return
		}

		log.SetOutput(io.MultiWriter(file, os.Stderr))
	}
}


