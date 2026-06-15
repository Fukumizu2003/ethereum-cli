/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"ethereum-cli/internal/util"
	"fmt"

	"github.com/spf13/cobra"
)

var test string

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		info, _ := util.GetChainInfo("ETH")
		data, _ := util.ReadBaseFee(info)
		fmt.Println(data)
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
	debugCmd.Flags().StringVarP(&test, "test", "t", "", "")
}
