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

var mainAcName string
var setChain string

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if mainAcName == "" && setChain == "" {
			return fmt.Errorf("設定内容を指定してください。\nアカウント変更: -n <アカウント名>\nチェーン変更: -c <チェーン名>\n対応チェーン: ETH, BNB, POL")
		}
		var newstate *config.State
		if mainAcName != "" {
			address, err := util.GetAddressFromName(mainAcName)
			if err != nil {
				return fmt.Errorf("このアカウント名は存在しません")
			}
			newstate, _ = config.ChangeMainAccount(address)
		}
		if setChain != "" {
			ns, err := config.SetMainChain(setChain)
			if err != nil {
				return err
			}
			newstate = ns
		}
		config.SaveConfig(*newstate)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&mainAcName, "name", "n", "", "")
	setCmd.Flags().StringVarP(&setChain, "chain", "c", "", "")
}
