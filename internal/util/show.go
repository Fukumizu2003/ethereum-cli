package util

import (
	"fmt"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

func ShowAllAddress(accounts [][]string) {
	fmt.Println("")
	names := []string{}
	addresses := []string{}
	for _, ac := range accounts {
		names = append(names, ac[0])
		addresses = append(addresses, ac[1])
	}
	for i, name := range names {
		fmt.Println(name + ": " + addresses[i])
	}
}

func QRCodeString(text string) (string, error) {
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("QRコードにする文字列を入力してください")
	}

	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		return "", err
	}
	return qr.ToSmallString(false), nil
}

func ShowQRCode(text string) error {
	qr, err := QRCodeString(text)
	if err != nil {
		return err
	}
	fmt.Print(qr)
	return nil
}
