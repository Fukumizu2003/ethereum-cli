package config

import (
	"encoding/json"
	"errors"
	"ethereum-cli/internal/util"
	"os"
	"path/filepath"
	"strings"
)

type State struct {
	Name    string
	Chain   string
	Address string
	Key     string
}

type Config struct{}

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
	var state State
	util.MkdirOrNothing("ref")
	f, _ := os.ReadFile(filepath.Join("ref", "state.json"))
	json.Unmarshal(f, &state)
	state.Chain = chain
	return &state, nil
}

func GetMainAccount() *State {
	var state State
	f, _ := os.ReadFile(filepath.Join("ref", "state.json"))
	json.Unmarshal(f, &state)
	return &state
}

func SaveConfig(st State) {
	stateSave, _ := json.MarshalIndent(st, "", "    ")
	os.WriteFile(filepath.Join("ref", "state.json"), stateSave, 0644)
}
