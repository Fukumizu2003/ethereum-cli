/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"ethereum-cli/internal/util"

	"github.com/spf13/cobra"
)

var getBasefee bool
var curGetInfo string

// getinfoCmd represents the getinfo command
var getinfoCmd = &cobra.Command{
	Use:   "getinfo",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if !getBasefee {
			fmt.Println("取得する情報の種類をフラグで指定してください。")
			return
		}
		if getBasefee {
			chaininfo, _ := util.GetChainInfo(curGetInfo)
			basefee, _ := util.ReadBaseFee(chaininfo)
			about := basefee - basefee%1000000
			gwei := util.GweiToEth(strconv.Itoa(int(about)))
			fmt.Println(gwei + " Gwei")
		}
	},
}

func init() {
	rootCmd.AddCommand(getinfoCmd)

	getinfoCmd.Flags().BoolVar(&getBasefee, "basefee", false, "")
	getinfoCmd.Flags().StringVarP(&curGetInfo, "chain", "c", "ETH", "")
}
