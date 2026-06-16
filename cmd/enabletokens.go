/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"ethereum-cli/internal/config"
	"ethereum-cli/internal/util"
	"fmt"

	"github.com/spf13/cobra"
)

// enabletokensCmd represents the enabletokens command
var enabletokensCmd = &cobra.Command{
	Use:   "enabletokens",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		apikey := config.AnkrAPIKey()
		if apikey == "" {
			fmt.Println("トークン機能の利用にはAnkrのAPI key（無料）が必要です。以下の手順に従ってください。\n① Ankr API keyの取得\n　　https://www.ankr.com/rpc/にアクセスし、アカウント作成、プロジェクト作成を行う。\n\n② API keyの設定\n　　以下のコマンドを実行してください。\n\n　　ethereum-cli set --apikey <取得したAPI key>")
			return
		}
		addr := config.GetMainAccount().Address
		info, err := util.GetAddressTokensInfo(apikey, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = util.AddTokensInfo(info)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(enabletokensCmd)
}
