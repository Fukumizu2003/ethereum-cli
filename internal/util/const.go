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

func TokenInfo(cur string) (int, []byte) {
	cur = strings.ToUpper(cur)
	var data map[string]interface{}
	b, _ := os.ReadFile(RelativeToAbsolute("ref", "const.json"))
	json.Unmarshal(b, &data)
	info := data[cur].(map[string]interface{})
	id, _ := hex.DecodeString(PureHex(info["ID"].(string)))
	return int(info["DECIMAL"].(float64)), id
}
