/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"ethereum-cli/internal/config"
	"ethereum-cli/internal/util"

	"github.com/spf13/cobra"
)

var mainAcName string

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if mainAcName == "" {
			return fmt.Errorf("アカウント名を -n で指定してください")
		}
		address, err := util.GetAddressFromName(mainAcName)
		if err != nil {
			return fmt.Errorf("このアカウント名は存在しません")
		}
		config.ChangeMainAccount(address)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&mainAcName, "name", "n", "", "")
}
