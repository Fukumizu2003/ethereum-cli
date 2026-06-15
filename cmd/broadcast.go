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

// broadcastCmd represents the broadcast command
var broadcastCmd = &cobra.Command{
	Use:   "broadcast",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		chain := config.GetMainAccount().Chain
		tx := util.LoadTx()
		raw := util.SignedRLP(tx)
		resp, err := util.Broadcast(raw, chain)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(resp))
	},
}

func init() {
	rootCmd.AddCommand(broadcastCmd)
}
