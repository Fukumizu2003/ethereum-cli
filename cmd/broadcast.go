/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"ethereum-cli/internal/util"
	"fmt"

	"github.com/spf13/cobra"
)

var curBroadcast string

// broadcastCmd represents the broadcast command
var broadcastCmd = &cobra.Command{
	Use:   "broadcast",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		tx := util.LoadTx()
		raw := util.SignedRLP(tx)
		resp, err := util.Broadcast(raw, curBroadcast)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(resp))
	},
}

func init() {
	rootCmd.AddCommand(broadcastCmd)
	broadcastCmd.Flags().StringVarP(&curBroadcast, "chain", "c", "ETH", "")
}
