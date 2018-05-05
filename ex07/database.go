package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

var db *sql.DB

const POOL_MAXOPENCONNS = 10

const POOL_MAXIDLECONNS = 2

func OpenDatabase() error {
	db1, err := sql.Open("mysql", "wanglei:123123@tcp(192.168.1.192:3306)/btc_dealing?charset=utf8")
	if err != nil {
		return err
	}

	//连接池
	db1.SetMaxOpenConns(POOL_MAXOPENCONNS)
	db1.SetMaxIdleConns(POOL_MAXIDLECONNS)
	//连接
	if  err = db1.Ping(); err != nil {
		return err
	}
	db = db1
	return nil
}

func CloseDatabase() {
	db.Close()
}

type MyTx struct {
	Tx *sql.Tx
}

func TxBegin() (*MyTx, error) {
	tx, err := db.Begin()
	return &MyTx{tx}, err
}

func (x *MyTx)Commit() error {
	return x.Tx.Commit()
}

func (x *MyTx)Rollback() error {
	return x.Tx.Rollback()
}

func (x *MyTx) InsertTx(tx *TxJson) error {
	stmt, err := x.Tx.Prepare("INSERT ethdata SET block_number=?,time_stamp=?,tx_hash=?,	nonce=?, block_hash=?, tx_index=?, from_addr=?, to_addr=?, contract_addr=?, amount=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(tx.BlockNumber, tx.TimeStamp, tx.Hash, tx.Nonce, tx.BlockHash, tx.TransactionIndex, tx.From, tx.To, tx.ContractAddress, tx.Value )
	if err != nil {
		return err
	}

	return err
}

func (x *MyTx) UpdateLastBlock(addr string, block int) (affect int64) {

	affect = 0

	stmt, err := x.Tx.Prepare("update address_log set last_block = ? where address = ?")
	if err != nil {
		return
	}

	res, err := stmt.Exec(strconv.Itoa(block), addr)
	if err != nil {
		return
	}

	affect, err = res.RowsAffected()

	return
}


func GetEthAddress() (addr string, block int){

	addr = ""
	block = 0

	//查询数据
	rows, err := db.Query("select address, last_block from address_log where type = 'eth' limit 1")
	if err != nil {
		return
	}


	if rows.Next() {
		a := ""
		b := 0
		err = rows.Scan(&a, &b)
		if err != nil {
			return
		}
		addr, block = a,b
	}

	return
}
