package util

import (
	"crypto/rand"
	"encoding/hex"
	"unicode"

	"github.com/btcsuite/btcd/btcec/v2"
	"golang.org/x/crypto/sha3"
)

func GenKey(len int) []byte {
	buf := make([]byte, len)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf
}

func NewKeypair() (*btcec.PrivateKey, *btcec.PublicKey) {
	privKey, _ := btcec.NewPrivateKey()
	pubKey := privKey.PubKey()
	return privKey, pubKey
}

func BytesToKeypair(priv []byte) (*btcec.PrivateKey, *btcec.PublicKey) {
	privKey, pubKey := btcec.PrivKeyFromBytes(priv)
	return privKey, pubKey
}

func Keccak(orig []byte) []byte {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(orig)
	return hasher.Sum(nil)
}

func PubkeyToAddress(pub []byte) string {
	addrBytes := []rune{'0', 'x'}
	hashed := Keccak(pub)
	info := hashed[len(hashed)-20:]
	infoHex := hex.EncodeToString(info)
	infoHexRunes := []rune(infoHex)
	infoHexBytes := []byte(infoHex)
	infoHashed := Keccak(infoHexBytes)
	infoHashedHex := hex.EncodeToString(infoHashed)
	for i, digit := range []rune(infoHashedHex) {
		if i == 40 {
			break
		}
		switch digit {
		case '0', '1', '2', '3', '4', '5', '6', '7':
			addrBytes = append(addrBytes, infoHexRunes[i])
		case '8', '9', 'a', 'b', 'c', 'd', 'e', 'f':
			addrBytes = append(addrBytes, unicode.ToUpper(infoHexRunes[i]))
		}
	}
	return string(addrBytes)
}

func NewPrivAddress() ([]byte, string) {
	privkey, pubkey := NewKeypair()
	priv := privkey.Serialize()
	pub := pubkey.SerializeUncompressed()[1:]
	address := PubkeyToAddress(pub)
	return priv, address
}
