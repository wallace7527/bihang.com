// Bitcoin Verified Message Signature
package bvms

import (
	"encoding/base64"
	"bytes"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"fmt"
)

const (
	MainNetID = 1
	RegressionNetID = 2
	TestNet3ID = 3
	SimNetID = 4
)

var nets = []*chaincfg.Params {
	&chaincfg.MainNetParams,
	&chaincfg.RegressionNetParams,
	&chaincfg.TestNet3Params,
	&chaincfg.SimNetParams,
}

func VerifyMessage( addr, signature, message string ) (error, int) {
	// Decode base64 signature.
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err, 0
	}

	// Validate the signature - this just shows that it was valid at all.
	// we will compare it with the key next.
	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, "Bitcoin Signed Message:\n")
	wire.WriteVarString(&buf, 0, message)
	expectedMessageHash := chainhash.DoubleHashB(buf.Bytes())
	pk, wasCompressed, err := btcec.RecoverCompact(btcec.S256(), sig,
		expectedMessageHash)
	if err != nil {
		// Mirror Bitcoin Core behavior, which treats error in
		// RecoverCompact as invalid signature.
		return err, 0
	}

	// Reconstruct the pubkey hash.
	var serializedPK []byte
	if wasCompressed {
		serializedPK = pk.SerializeCompressed()
	} else {
		serializedPK = pk.SerializeUncompressed()
	}

	for i, net := range nets {
		address, err := btcutil.NewAddressPubKey(serializedPK, net)
		if err != nil {
			// Again mirror Bitcoin Core behavior, which treats error in public key
			// reconstruction as invalid signature.
			return err, 0
		}

		if address.EncodeAddress() == addr {
			return nil, i+1
		}
	}

	return fmt.Errorf("The Signature Message Verification Failed."), 0
}

func IsValidAddress( address string ) bool {
	// Decode the provided address.
	for _, net := range nets {
		addr, err := btcutil.DecodeAddress(address, net)
		if err != nil {
			continue
		}
		// Only P2PKH addresses are valid for signing.
		if _, ok := addr.(*btcutil.AddressPubKeyHash); ok {
			return true
		}
	}

	return false
}