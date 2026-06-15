/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"ethereum-cli/internal/config"
	"ethereum-cli/internal/util"
	"strings"

	"github.com/spf13/cobra"
)

var showAddress bool
var showBalance bool
var showAll bool
var showGwei bool
var showWei bool

var showToken string

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if (showAddress && showBalance) || (!showAddress && !showBalance) {
			return
		}
		if showAll {
			if showAddress {
				accounts := util.LoadAccounts()
				accounts = append(accounts, util.LoadDestinations()...)
				util.ShowAllAddress(accounts)
			}
		} else {
			if showAddress {
				ac := config.GetMainAccount()
				fmt.Println((ac.Address))
				util.ShowQRCode(ac.Address)
			} else if showBalance {
				ac := config.GetMainAccount()
				chain := ac.Chain
				if showToken == "" {
					balance, err := util.GetBalance(ac.Address, chain)
					if err != nil {
						fmt.Println("情報取得失敗")
						return
					}
					weistr := balance.String()
					if showWei {
						fmt.Println("Balance: " + weistr + " wei")
						return
					}
					if showGwei {
						fmt.Println("Balance: " + util.GweiToEth(weistr) + " Gwei")
						return
					}
					ethstr := util.WeiToEth(weistr)
					if len(ethstr) > 11 {
						ethstr = string([]rune(ethstr)[:11])
					}
					fmt.Println("Balance: " + ethstr + " " + strings.ToUpper(chain))
				} else {
					balance, err := util.GetTokenBalance(ac.Address, chain, showToken)
					if err != nil {
						fmt.Println("情報取得失敗")
						fmt.Println(balance)
						return
					}
					fmt.Println("Balance: " + balance + " " + strings.ToUpper(showToken))
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.Flags().BoolVarP(&showAddress, "address", "a", false, "")
	showCmd.Flags().BoolVarP(&showBalance, "balance", "b", false, "")
	showCmd.Flags().BoolVar(&showAll, "all", false, "")
	showCmd.Flags().BoolVar(&showGwei, "gwei", false, "")
	showCmd.Flags().BoolVar(&showWei, "wei", false, "")
	showCmd.Flags().StringVar(&showToken, "token", "", "")
}
