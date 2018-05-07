//这个例子展示了比特币区块链区块HASH的计算
//块HASH计算需要注意的是
// 小端字节序
// 和两次HASH运算

package main

import (
	"fmt"
	"time"
	"bihang.com/ex04/wire"
	"bihang.com/ex04/chainhash"
)

func newHashFromStr(hexStr string) *chainhash.Hash {
	hash, err := chainhash.NewHashFromStr(hexStr)
	if err != nil {
		panic(err)
	}
	return hash
}


func main() {

	header := &wire.BlockHeader{
		Version:0x02,
		PrevBlock:*newHashFromStr("0000000000000002a7bbd25a417c0374cc55261021e8a9ca74442b01284f0569"),
		MerkleRoot:*newHashFromStr("c91c008c26e50763e9f548bb8b2fc323735f73577effbc55502c51eb4cc7cf2e"),
		Timestamp:time.Unix(1388185914,0),
		Bits:0x1903a30c,
		Nonce:0x371c2688,
	}

	fmt.Printf("%s\n", header.BlockHash())

	tm, _ := time.Parse("2006-01-02 15:04:05","2018-04-09 06:04:43")

	header = &wire.BlockHeader{
		Version:0x20000000,
		PrevBlock:*newHashFromStr("00000000000000000004cfea2aa9013e71e6606239cf87bf242b25f49bb77b4c"),
		MerkleRoot:*newHashFromStr("42510a9129a8d55d4ed9d9b1477bf3de588a7bcee71dec23f7cf5d9901992586"),
		Timestamp:tm,
		Bits:0x17502ab7,
		Nonce:0x4b487c50,
	}

	fmt.Printf("%s\n", header.BlockHash())


}
