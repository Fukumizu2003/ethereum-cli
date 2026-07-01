/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/csv"
	"ethereum-cli/internal/util"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var deleteName string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if deleteName == "" {
			return fmt.Errorf("削除するアカウント名を -n で指定してください")
		}
		accounts := util.LoadAccounts()
		destinations := util.LoadDestinations()
		newAccounts := [][]string{}
		newDestinations := [][]string{}
		acflag := false
		deflag := false
		for _, ac := range accounts {
			if ac[0] != deleteName {
				newAccounts = append(newAccounts, ac)
			} else {
				acflag = true
			}
		}
		for _, ds := range destinations {
			if ds[0] != deleteName {
				newDestinations = append(newDestinations, ds)
			} else {
				deflag = true
			}
		}
		if !acflag && !deflag {
			return fmt.Errorf("アカウント名が存在しません")
		}
		var buf bytes.Buffer
		writer := csv.NewWriter(&buf)
		if acflag {
			writer.WriteAll(newAccounts)
			return os.WriteFile(util.RelativeToAbsolute("ref", "ETH_keypair.csv"), buf.Bytes(), 0644)
		} else if deflag {
			writer.WriteAll(newDestinations)
			return os.WriteFile(util.RelativeToAbsolute("ref", "ETH_destinations.csv"), buf.Bytes(), 0644)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&deleteName, "name", "n", "", "削除するアカウントの名前を指定")
}
