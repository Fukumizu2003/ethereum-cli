package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

func GetAddressTokensInfo(key string, addr string) ([]byte, error) {
	url := "https://rpc.ankr.com/multichain/" + key
	payload := `{
	"jsonrpc": "2.0",
	"method": "ankr_getAccountBalance",
	"params": {
		"walletAddress": "` + addr + `",
		"blockchain": ["eth", "bsc", "polygon"] 
	},
	"id": 1
	}`
	info, err := PostHTTP(payload, url)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func AddTokensInfo(info []byte) error {
	save := LoadTokensInfo()
	var buf map[string]interface{}
	json.Unmarshal(info, &buf)
	if buf["result"] == nil {
		return errors.New("\"result\"カラムなし")
	}
	result := buf["result"].(map[string]interface{})
	if result["assets"] == nil {
		return errors.New("トークン情報なし")
	}
	assets := result["assets"].([]interface{})
	for _, tkinfo := range assets {
		tkinfoAssert := tkinfo.(map[string]interface{})
		chain := tkinfoAssert["blockchain"]
		chainNativeUnit := ""
		switch chain {
		case "eth":
			chainNativeUnit = "ETH"
		case "bsc":
			chainNativeUnit = "BNB"
		case "polygon":
			chainNativeUnit = "POL"
		}
		sy, dc, id, e := readTokenInfo(tkinfoAssert)
		if e != nil {
			continue
		}
		if AddTokenInfo(save, chainNativeUnit, sy, dc, id) {
			fmt.Println("ADDED: " + strings.ToUpper(sy) + " on " + chainNativeUnit + " (" + id + ")")
		}
	}
	SaveTokensInfo(save)
	return nil
}

func AddTokenInfo(origJson *map[string]map[string]map[string]interface{}, chain string, token string, decimal int, id string) bool {
	if !ValidChain(chain) || id == "0x"+strings.Repeat("00", 20) {
		return false
	}
	chain = strings.ToUpper(chain)
	tokens := (*origJson)[chain]
	tokens[token] = make(map[string]interface{})
	tokens[token]["ID"] = id
	tokens[token]["DECIMAL"] = decimal
	(*origJson)[chain] = tokens
	return true
}

func readTokenInfo(assetscol map[string]interface{}) (string, int, string, error) {
	if assetscol["tokenDecimals"] == nil || assetscol["contractAddress"] == nil {
		return "", -1, "", errors.New("トークン情報カラムなし")
	}
	sym := assetscol["tokenSymbol"].(string)
	decimal := assetscol["tokenDecimals"].(float64)
	id := assetscol["contractAddress"].(string)
	return sym, int(decimal), id, nil
}

func LoadTokensInfo() *map[string]map[string]map[string]interface{} {
	var res map[string]map[string]map[string]interface{}
	b, _ := os.ReadFile(RelativeToAbsolute("ref", "ETH_const.json"))
	json.Unmarshal(b, &res)
	return &res
}

func SaveTokensInfo(info *map[string]map[string]map[string]interface{}) {
	buf, _ := json.MarshalIndent(*info, "", "    ")
	os.WriteFile(RelativeToAbsolute("ref", "ETH_const.json"), buf, 0644)
}
