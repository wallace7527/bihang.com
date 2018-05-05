package main

import (
	"fmt"
	"net/http"
	"encoding/json"
)


const ETHERSCANAPI_TXLIST_FMT = "http://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=%d&endblock=99999999&sort=asc&apikey=GI5M7WYJXTGQYGF747KKMZKSRTF3PXS3PZ"


type (
	TxJson struct {
		BlockNumber	string	`json:"blockNumber"`
		TimeStamp	string	`json:"timeStamp"`
		Hash	string	`json:"hash"`
		Nonce	string	`json:"nonce"`
		BlockHash	string	`json:"blockHash"`
		TransactionIndex	string	`json:"transactionIndex"`
		From	string	`json:"from"`
		To	string	`json:"to"`
		Value	string	`json:"value"`
		Gas	string	`json:"gas"`
		GasPrice	string	`json:"gasPrice"`
		Input	string	`json:"input"`
		ContractAddress	string	`json:"contractAddress"`
		CumulativeGasUsed	string	`json:"cumulativeGasUsed"`
		GasUsed	string	`json:"gasUsed"`
		Confirmations	string	`json:"confirmations"`
		IsError string `json:"isError"`
	}

	TxlistJson struct {
		Status string  `json:"status"`
		Message string  `json:"message"`
		Result [] TxJson `json:"result"`
	}

)

func Retrieve(address string, startBlock int, skipLastBlock bool ) (*TxlistJson, error) {

	if skipLastBlock {
		startBlock++
	}
	// Retrieve the rss feed document from the web.
	uri := fmt.Sprintf(ETHERSCANAPI_TXLIST_FMT, address, startBlock )
	//fmt.Println("URI:", uri)
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	// Close the response once we return from the function.
	defer resp.Body.Close()

	// Check the status code for a 200 so we know we have received a
	// proper response.
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Response Error %d\n", resp.StatusCode)
	}

	// Decode the rss feed document into our struct type.
	// We don't need to check for errors, the caller can do this.
	var txlistJson TxlistJson
	err = json.NewDecoder(resp.Body).Decode(&txlistJson)
	return &txlistJson,err
}

