package util

import (
	"encoding/csv"
	"errors"
	"os"
)

func MkdirOrNothing(dir string) {
	os.MkdirAll(dir, 0755)
}

func LoadAccounts() [][]string {
	MkdirOrNothing("ref")
	f, _ := os.Open(RelativeToAbsolute("ref", "keypair.csv"))
	defer f.Close()
	data, _ := csv.NewReader(f).ReadAll()
	return data
}

func LoadDestinations() [][]string {
	MkdirOrNothing("ref")
	f, _ := os.Open(RelativeToAbsolute("ref", "destinations.csv"))
	defer f.Close()
	data, _ := csv.NewReader(f).ReadAll()
	return data
}

func GetAddressFromName(name string) (string, error) {
	accounts := LoadAccounts()
	accounts = append(accounts, LoadDestinations()...)
	addr := ""
	flag := false
	for _, ac := range accounts {
		if ac[0] == name {
			addr = ac[1]
			flag = true
			break
		}
	}
	if !flag {
		return "", errors.New("指定のアカウント名は存在しません")
	}
	return addr, nil
}

func NameToAddress(name string) string {
	addr, err := GetAddressFromName(name)
	if err != nil {
		return name
	}
	return addr
}

func SaveKeypair(acname string, address string, priv []byte) {
	MkdirOrNothing("ref")
	privB64 := B64Encode(priv)
	row := []byte{}
	row = append(row, []byte(acname)...)
	row = append(row, byte(','))
	row = append(row, []byte(address)...)
	row = append(row, byte(','))
	row = append(row, []byte(privB64)...)
	row = append(row, []byte("\n")...)
	appendFile(RelativeToAbsolute("ref", "keypair.csv"), row)
}

func SaveAddress(acname string, address string) {
	MkdirOrNothing("ref")
	row := []byte{}
	row = append(row, []byte(acname)...)
	row = append(row, byte(','))
	row = append(row, []byte(address)...)
	row = append(row, byte('\n'))
	appendFile(RelativeToAbsolute("ref", "destinations.csv"), row)
}

func appendFile(path string, row []byte) {
	existing, _ := os.ReadFile(path)
	existing = append(existing, row...)
	os.WriteFile(path, existing, 0644)
}

func CheckName(acs [][]string, dss [][]string, name string) bool {
	for _, ac := range acs {
		if ac[0] == name {
			return false
		}
	}
	for _, ds := range dss {
		if ds[0] == name {
			return false
		}
	}
	return true
}

func SaveResp(data []byte) {
	os.WriteFile(RelativeToAbsolute("temp", "resp.json"), data, 0644)
}
