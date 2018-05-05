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

	proc := 0
	inc := 0
	skip := 0

	t1 := time.Now()
	log.Println("Begin retrieve.")
	defer func (){
		elapsed := time.Since(t1)
		log.Printf("End retrieve.(elapsed:%v, process:%d, increase:%d, skip:%d) ", elapsed, proc, inc, skip)
	}()

	if err := OpenDatabase(); err != nil {
		log.Println("Error OpenDatabase:", err.Error())
		return
	}
	defer CloseDatabase()

	//取得地址和最后一次更新的块号
	addr, lastBlock := GetEthAddress()
	if len(addr) == 0 {
		log.Println("Error ethereum address is empty.")
		return
	}

	maxBlock := lastBlock


	log.Printf("Retrieve: %s%s, lastBlock:%d", ETHERSCANAPI_ADDR, addr, lastBlock)
	//使用Etherscan API检索交易列表
	txlistJson, err := Retrieve(addr, lastBlock, false)
	//检索失败处理
	if err != nil {
		log.Println("Error retrieve:", err.Error())
		return
	}

	//API返回错误处理
	if txlistJson.Status != "1" {
		log.Println("etherscan api:", txlistJson.Message)
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


		//保存最后块编号
		b, ok := strconv.Atoi(tx.BlockNumber )
		if ok == nil {
			if b > maxBlock {
				maxBlock = b
			}
		}
	}

	//保存最后一次的地址
	if maxBlock > lastBlock {
		if trans.UpdateLastBlock(addr, maxBlock) == 0 {
			trans.Rollback()
		}
	}

	trans.Commit()
}
