package util

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
)

const (
	ETH_URL = "https://ethereum-rpc.publicnode.com"
	BNB_URL = "https://bsc-rpc.publicnode.com"
	POL_URL = "https://polygon.publicnode.com"
	ETH_ID  = 0x01
	BNB_ID  = 0x38
	POL_ID  = 0x89
)

func GetChainId(cur string) []byte {
	curUp := strings.ToUpper(cur)
	switch curUp {
	case "ETH":
		return []byte{ETH_ID}
	case "BNB":
		return []byte{BNB_ID}
	case "POL":
		return []byte{POL_ID}
	}
	return nil
}

func GetNodeURL(cur string) string {
	curUp := strings.ToUpper(cur)
	switch curUp {
	case "ETH":
		return ETH_URL
	case "BNB":
		return BNB_URL
	case "POL":
		return POL_URL
	}
	return ""
}

func TokenInfo(chain string, token string) (int, []byte) {
	chain = strings.ToUpper(chain)
	token = strings.ToUpper(token)
	var data map[string]interface{}
	b, _ := os.ReadFile(RelativeToAbsolute("ref", "const.json"))
	json.Unmarshal(b, &data)
	info := data[chain].(map[string]interface{})
	pair := info[token].(map[string]interface{})
	var id []byte
	if pair["ID"] != nil {
		id, _ = hex.DecodeString(PureHex(pair["ID"].(string)))
	}
	return int(pair["DECIMAL"].(float64)), id
}

func ValidChain(chain string) bool {
	chain = strings.ToUpper(chain)
	switch chain {
	case "ETH":
		fallthrough
	case "BNB":
		fallthrough
	case "POL":
		return true
	}
	return false
}
