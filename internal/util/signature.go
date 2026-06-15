package util

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
)

func GetSignature(tx Tx, privb []byte) ([]byte, []byte, []byte) {
	priv, _ := btcec.PrivKeyFromBytes(privb)

	presig := PreSigRLP(tx)
	hashed := Keccak(presig)

	sig := ecdsa.SignCompact(priv, hashed, true)

	headerByte := sig[0]
	parity := []byte{byte((headerByte + 1) & 1)}
	r := sig[1:33]
	s := sig[33:65]

	return r, s, parity
}

func Sign(tx *Tx, priv []byte) {
	r, s, parity := GetSignature(*tx, priv)
	tx.R = r
	tx.S = s
	tx.YParity = parity
}
