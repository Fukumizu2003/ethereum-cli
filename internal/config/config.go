package config

import (
	"errors"
	"ethereum-cli/internal/util"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type State struct {
	Name    string
	Chain   string
	Address string
	Key     string
}

func ChangeMainAccount(address string) (*State, error) {
	var state State

	accounts := util.LoadAccounts()
	flag := false
	for _, ac := range accounts {
		if address == ac[1] {
			state.Address = address
			state.Name = ac[0]
			state.Key = ac[2]
			flag = true
			break
		}
	}
	if !flag {
		return nil, errors.New("このアドレスは登録されていません")
	}
	return &state, nil
}

func SetMainChain(chain string) (*State, error) {
	if !util.ValidChain(chain) {
		return nil, errors.New("非対応チェーンです。")
	}
	chain = strings.ToUpper(chain)
	state := GetMainAccount()
	state.Chain = chain
	return state, nil
}

func GetMainAccount() *State {
	godotenv.Load()
	var state State
	state.Name = os.Getenv("NAME_ETH")
	state.Chain = os.Getenv("CHAIN_ETH")
	state.Address = os.Getenv("ADDRESS_ETH")
	state.Key = os.Getenv("PRIVKEY_ENCRYPTED_ETH")
	return &state
}

func SaveConfig(st State) {
	curr, err := godotenv.Read(".env")
	if err != nil {
		curr = make(map[string]string)
	}
	curr["NAME_ETH"] = st.Name
	curr["CHAIN_ETH"] = st.Chain
	curr["ADDRESS_ETH"] = st.Address
	curr["PRIVKEY_ENCRYPTED_ETH"] = st.Key
	godotenv.Write(curr, ".env")
}

func AnkrAPIKey() string {
	godotenv.Load()
	return os.Getenv("ANKR_API")
}

func AnkrAPIKeySet(key string) {
	curr, err := godotenv.Read(".env")
	if err != nil {
		curr = make(map[string]string)
	}
	curr["ANKR_API"] = key
	godotenv.Write(curr, ".env")
}
