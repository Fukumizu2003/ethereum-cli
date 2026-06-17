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

func PostHTTP(payload string, url string) ([]byte, error) {
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

func PostEthNode(payload string, chain string) ([]byte, error) {
	url := GetNodeURL(chain)
	body, err := PostHTTP(payload, url)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func readResult(body []byte) (string, error) {
	var js map[string]interface{}
	json.Unmarshal(body, &js)
	if js["result"] == nil {
		return "", errors.New("\"result\"カラムなし")
	}
	return js["result"].(string), nil
}

func Broadcast(raw []byte, chain string) ([]byte, error) {
	payload := `{
	"jsonrpc":"2.0",
	"method":"eth_sendRawTransaction",
	"params":[
	"0x` + hex.EncodeToString(raw) + `"
	],
	"id":1
	}`
	body, err := PostEthNode(payload, chain)
	if err != nil {
		return nil, err
	}
	result, err := readResult(body)
	if err != nil {
		return body, err
	}
	return []byte("SUCCEED: " + result), nil
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
	result, err := readResult(body)
	if err != nil {
		return nil, errors.New("取得失敗")
	}
	n3, _ := new(big.Int).SetString(result, 0)
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
	result, err := readResult(body)
	if err != nil {
		jsbyte, _ := json.Marshal(body)
		return string(jsbyte), errors.New("取得失敗")
	}
	n3, _ := new(big.Int).SetString(result, 0)
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
	result, err := readResult(body)
	if err != nil {
		return nil, err
	}
	ans, _ := hex.DecodeString(PureHex(result))
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
	SaveResp(body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetBaseFee(chain string) (uint64, error) {
	chainInfo, err := GetChainInfo(chain)
	if err != nil {
		return 0, err
	}
	basefee, err := ReadBaseFee(chainInfo)
	if err != nil {
		return 0, err
	}
	return basefee, nil
}

func GetGasPrice(chain string) (uint64, error) {
	payload := `{
	"jsonrpc":"2.0",
	"method":"eth_gasPrice",
	"params":[],
	"id":1
	}`
	body, err := PostEthNode(payload, chain)
	if err != nil {
		return 0, err
	}
	result, err := readResult(body)
	if err != nil {
		return 0, err
	}
	ans, _ := hex.DecodeString(PureHex(result))
	return uint64(BytesToInt(ans)), nil
}
