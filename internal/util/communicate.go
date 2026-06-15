package util

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"time"
)

func PostEthNode(payload string, cur string) ([]byte, error) {
	url := GetNodeURL(cur)
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(payload)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body, nil
}

func Broadcast(raw []byte, cur string) ([]byte, error) {
	payload := `{
	"jsonrpc":"2.0",
	"method":"eth_sendRawTransaction",
	"params":[
	"0x` + hex.EncodeToString(raw) + `"
	],
	"id":1
	}`
	body, err := PostEthNode(payload, cur)
	if err != nil {
		return nil, err
	}
	var buf map[string]interface{}
	json.Unmarshal(body, &buf)
	if buf["result"] != nil {
		result := buf["result"].(string)
		if result != "" {
			return []byte("SUCCEED: " + result), nil
		}
	}
	return body, nil
}

func GetBalance(addr string, chain string) (*big.Int, error) {
	payload := `{
	"jsonrpc":"2.0",
	"method":"eth_getBalance",
	"params":[
	"` + addr + `",
		"latest"
	],
	"id":1
	}`
	body, err := PostEthNode(payload, chain)
	if err != nil {
		return nil, err
	}

	var js map[string]interface{}
	json.Unmarshal(body, &js)
	if js["result"] == nil {
		return nil, errors.New("取得失敗")
	}
	n3, _ := new(big.Int).SetString(js["result"].(string), 0)
	return n3, nil
}

func GetTokenBalance(addr string, chain string, token string) (string, error) {
	decimal, id := TokenInfo(chain, token)
	idhex := hex.EncodeToString(id)
	payload := `{
	"jsonrpc":"2.0",
	"method":"eth_call",
	"params":[
		{
			"to":"0x` + idhex + `",
			"data":"0x70a08231000000000000000000000000` + PureHex(addr) + `"
		},
		"latest"
	],
	"id":1
	}`
	body, err := PostEthNode(payload, chain)
	if err != nil {
		return "", err
	}

	var js map[string]interface{}
	json.Unmarshal(body, &js)
	if js["result"] == nil {
		jsbyte, _ := json.Marshal(js)
		return string(jsbyte), errors.New("取得失敗")
	}
	n3, _ := new(big.Int).SetString(js["result"].(string), 0)
	balstr := n3.String()
	return IntstrToFloatstr(balstr, decimal), nil
}

func GetNonce(addr string, chain string) ([]byte, error) {
	payload := `{
	"jsonrpc":"2.0",
	"method":"eth_getTransactionCount",
	"params":[
	"` + addr + `",
		"latest"
	],
	"id":1
	}`
	body, err := PostEthNode(payload, chain)
	if err != nil {
		return nil, err
	}

	var js map[string]interface{}
	json.Unmarshal(body, &js)
	if js["result"] == "" {
		return nil, errors.New("取得失敗")
	}
	ans, _ := hex.DecodeString(PureHex(js["result"].(string)))
	return ans, nil
}

func GetChainInfo(chain string) ([]byte, error) {
	payload := `{
	"jsonrpc":"2.0",
	"method":"eth_getBlockByNumber",
	"params":[
		"latest",
		false
	],
	"id":1
	}`
	body, err := PostEthNode(payload, chain)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func ReadBaseFee(chaininfo []byte) (uint64, error) {
	var data map[string]interface{}
	json.Unmarshal(chaininfo, &data)
	result := data["result"].(map[string]interface{})
	if result == nil {
		return 0, errors.New("情報取得失敗")
	}
	bfpg := result["baseFeePerGas"].(string)
	basefeehex := PureHex(bfpg)
	basefeeBytes, _ := hex.DecodeString(basefeehex)
	return BytesToInt(basefeeBytes), nil
}
