
/*
Copyright 2018 The go-execsql Authors
This file is part of the go-execsql library.

main
//1.从配置文件读取数据库连接与SQL语句
//2.连接数据，执行SQL语句
//3.操作过程日志输出到log/下，以日期切分的日志文件

wanglei.ok@foxmail.com

1.0
版本时间：16:22 2018/5/9

*/

package main

import (
	"log"
	"time"
)

func main() {

	if err := readConfig(); err != nil {
		log.Println("Error readConfig:", err)
		return
	}

	if err := OpenDatabase(config.Mysql.DSN); err != nil {
		log.Println( "Error OpenDatabase:", err )
		return
	}
	defer CloseDatabase()

	//Action
	if config.Action.UseTrans {

		//使用事务，错误出现回滚错误
		trans, err := TxBegin()
		if err != nil {
			log.Println( "Error TxBegin:", err)
			return
		}

		for _, s := range config.Action.SQL {
			log.Println( "Exec:", s)
			t1 := time.Now()
			result, err := trans.ExecSQL(s)
			if err != nil {
				log.Println( "Error ExecSQL:", err)
				trans.Rollback()
				log.Printf( "Exec end.(Elapsed:%v,Rollbacked)", time.Since(t1))
				return
			}

			log.Printf( "Exec end.(Elapsed:%v, %s)", time.Since(t1) , ResultString(result))
		}
		trans.Commit()

	}else {
		//不使用事务,忽略错误继续执行
		for _, s := range config.Action.SQL {
			log.Println( "Exec:", s)
			t1 := time.Now()
			result, err := ExecSQL(s)
			if err != nil {
				log.Println( "Error ExecSQL:", err)
				log.Printf( "Exec end.(Elapsed:%v)",time.Since(t1))
				continue
			}

			log.Printf( "Exec end.(Elapsed:%v, %s)", time.Since(t1), ResultString(result) )
		}
	}

}