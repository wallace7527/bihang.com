/*
Copyright 2018 The go-eam Authors
This file is part of the go-eam library.

go-eam
以太坊地址交易监控程序，使用以太坊api接口，查询指定以太坊地址下的交易记录。
并将记录数据保存到数据库中


wanglei.ok@foxmail.com

1.0
版本时间：2018年5月7日18:32:12

*/

package main

import (
	"strconv"
	"time"
	"strings"
	"log"
)

const (
	ETHERSCANAPI_ADDR   = "https://etherscan.io/address/"
	ETHERSCANAPI_TX     = "https://etherscan.io/tx/"
)

func init() {
	logSetup()
}

func main() {
	//打开数据库
	if err := OpenDatabase(); err != nil {
		log.Println("Error OpenDatabase:", err.Error())
		return
	}
	//程序结束前关闭数据库
	defer CloseDatabase()

	//取得地址和最后一次更新的块号
	eais, err := GetEthAddress()
	if err != nil {
		log.Println("Error GetEthAddress:", err.Error())
		return
	}

	if len(eais) == 0 {
		log.Println("Error not found ethereum address.")
		return
	}
	for _,eai := range(eais) {
		retrieve(eai.Address,eai.LastBlock)
	}
}

//取回指定地址和开始块号
func retrieve(addr string, startBlock int) {
	proc := 0
	inc := 0
	skip := 0

	t1 := time.Now()
	log.Printf("Begin retrieve.(Address: %s%s, StartBlock:%d)\n", ETHERSCANAPI_ADDR, addr, startBlock)
	defer func (){
		elapsed := time.Since(t1)
		log.Printf("End retrieve.(Elapsed:%v, Process:%d, Increase:%d, Skip:%d)\n", elapsed, proc, inc, skip)
	}()

	maxBlock := startBlock

	//使用Etherscan API检索交易列表
	txlistJson, err := Retrieve(addr, startBlock, true)
	//检索失败处理
	if err != nil {
		log.Println("Error retrieve:", err.Error())
		return
	}

	//API返回错误处理
	if txlistJson.Status != "1" {
		log.Println("Etherscan api:", txlistJson.Message)
		return
	}

	//没有交易记录
	if txlistJson.Result == nil {
		return
	}

	//遍历交易记录插入数据库
	//开启事务
	trans, err := TxBegin()
	for _ , tx := range txlistJson.Result {
		proc++
		err := trans.InsertTx(&tx)
		if err != nil {
			//插入失败显示日志
			//txString, _ := json.Marshal(tx)
			errString := err.Error()
			if strings.Contains(errString, "Duplicate entry") {
				log.Printf("Skip Duplicate tx: %s%s\n", ETHERSCANAPI_TX,tx.Hash)
			} else {
				log.Println("Error insertTx:", errString, tx.Hash)
			}
			skip++
		}else{
			inc++
		}


		//最后块编号
		b, ok := strconv.Atoi(tx.BlockNumber )
		if ok == nil {
			if b > maxBlock {
				maxBlock = b
			}
		}
	}

	//保存最后块编号
	if maxBlock > startBlock {
		if trans.UpdateLastBlock(addr, maxBlock) == 0 {
			trans.Rollback()
		}
	}

	trans.Commit()
}
