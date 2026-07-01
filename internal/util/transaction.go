package util

import (
	"encoding/json"
	"fmt"
	"os"
)

type Tx struct {
	TxType               []byte
	ChainId              []byte
	Nonce                []byte
	MaxPriorityFeePerGas []byte
	MaxFeePerGas         []byte
	GasLimit             []byte
	To                   []byte
	Value                []byte
	Data                 []byte
	AccessList           []byte
	YParity              []byte
	R                    []byte
	S                    []byte
}

func PreSigRLP(tx Tx) []byte {
	txlist := []byte{}
	txlist = append(txlist, RLPconv(tx.ChainId, false)...)
	txlist = append(txlist, RLPconv(tx.Nonce, false)...)
	txlist = append(txlist, RLPconv(tx.MaxPriorityFeePerGas, false)...)
	txlist = append(txlist, RLPconv(tx.MaxFeePerGas, false)...)
	txlist = append(txlist, RLPconv(tx.GasLimit, false)...)
	txlist = append(txlist, RLPconv(tx.To, true)...)
	txlist = append(txlist, RLPconv(tx.Value, false)...)
	txlist = append(txlist, RLPconv(tx.Data, true)...)
	txlist = append(txlist, RLPlistConv(tx.AccessList)...)
	bin := RLPlistConv(txlist)
	return append(tx.TxType, bin...)
}

func SignedRLP(tx Tx) []byte {
	txlist := []byte{}
	txlist = append(txlist, RLPconv(tx.ChainId, false)...)
	txlist = append(txlist, RLPconv(tx.Nonce, false)...)
	txlist = append(txlist, RLPconv(tx.MaxPriorityFeePerGas, false)...)
	txlist = append(txlist, RLPconv(tx.MaxFeePerGas, false)...)
	txlist = append(txlist, RLPconv(tx.GasLimit, false)...)
	txlist = append(txlist, RLPconv(tx.To, true)...)
	txlist = append(txlist, RLPconv(tx.Value, false)...)
	txlist = append(txlist, RLPconv(tx.Data, false)...)
	txlist = append(txlist, RLPlistConv(tx.AccessList)...)
	txlist = append(txlist, RLPconv(tx.YParity, false)...)
	txlist = append(txlist, RLPconv(tx.R, true)...)
	txlist = append(txlist, RLPconv(tx.S, true)...)
	bin := RLPlistConv(txlist)
	return append(tx.TxType, bin...)
}

func NewTx() Tx {
	return Tx{}
}

func LoadTx() Tx {
	MkdirOrNothing("temp")
	var tx Tx
	data, _ := os.ReadFile(RelativeToAbsolute("temp", "ETH_transaction.json"))
	json.Unmarshal(data, &tx)
	return tx
}

func SaveTx(tx Tx) {
	MkdirOrNothing("temp")
	data, err := json.MarshalIndent(tx, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	os.WriteFile(RelativeToAbsolute("temp", "ETH_transaction.json"), data, 0644)
}

func InitNativeTx(tx *Tx, chain string) {
	tx.TxType = []byte{0x02}
	tx.ChainId = GetChainId(chain)
	tx.GasLimit = []byte{0x52, 0x08}
}

func InitTokenTx(tx *Tx, chain string) {
	tx.TxType = []byte{0x02}
	tx.ChainId = GetChainId(chain)
	tx.GasLimit = []byte{0x01, 0x86, 0xa0}
}
