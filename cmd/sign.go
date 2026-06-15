/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"ethereum-cli/internal/config"
	"ethereum-cli/internal/util"
	"fmt"

	"github.com/spf13/cobra"
)

var password string

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ac := config.GetMainAccount()
		tx := util.LoadTx()
		priv, err := util.AesDecrypt(util.B64Decode(ac.Key), []byte(password))
		if err != nil {
			fmt.Println("パスワードが違います。")
			return
		}
		util.Sign(&tx, priv)
		raw := util.SignedRLP(tx)
		fmt.Println("\nRaw Transaction (HEX):")
		fmt.Println(hex.EncodeToString(raw))
		util.SaveTx(tx)
	},
}

func init() {
	rootCmd.AddCommand(signCmd)

	signCmd.Flags().StringVarP(&password, "password", "p", "", "")
}
