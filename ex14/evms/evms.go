// Ethereum Verified Message Signature
package evms

import (
	"encoding/hex"
	"fmt"
	"log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var signhashers = []struct {
	fn func([]byte) []byte
	msg string
}{
	{func(data[]byte) []byte{
		msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
		return crypto.Keccak256([]byte(msg))
	},"Pass(1)"},
	{func(data[]byte) []byte{
		return crypto.Keccak256([]byte(data))
	},"pass(2)"},
	{func(data[]byte) []byte{
		d := crypto.Keccak256([]byte(data))
		msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(d), d)
		return crypto.Keccak256([]byte(msg))
	},"Pass geth prefix(3)"},
}

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func VerifyMessage(addr, sign, msg string) (bool, string) {
	addressStr := addr
	signatureHex := sign
	message := []byte(msg)

	if !common.IsHexAddress(addressStr) {
		return false, fmt.Sprintf("Invalid address: %s", addressStr )
	}
	address := common.HexToAddress(addressStr)

	var signature []byte
	var err error
	if has0xPrefix(signatureHex) {
		signature, err = hex.DecodeString(signatureHex[2:])
	}else {
		signature, err = hex.DecodeString(signatureHex)
	}

	if err != nil {
		return false, fmt.Sprintf("Signature encoding is not hexadecimal: %v", err )
	}

	if len(signature) != 65 {
		return false, fmt.Sprintf("Signature must be 65 bytes long")
	}
	if signature[64] == 27 || signature[64] == 28 {
		signature[64] -= 27 // Transform yellow paper V from 27/28 to 0/1
	}

	for _, sh := range signhashers {
		recoveredPubkey, err := crypto.SigToPub(sh.fn(message), signature)
		if err != nil || recoveredPubkey == nil {
			log.Fatalf("Signature verification failed: %v", err)
		}

		recoveredAddress := crypto.PubkeyToAddress(*recoveredPubkey)
		if address == recoveredAddress {
			return true, sh.msg
		}
	}

	return false, ""
}
