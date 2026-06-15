/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"ethereum-cli/internal/config"
	"ethereum-cli/internal/util"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var amountStr string
var gasLimit string
var sendAll bool
var sendTo string
var maxfeeGweiStr string
var feeAbout string
var chainPay string
var token string

var paymentCmd = &cobra.Command{
	Use:   "payment",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		proc := true
		if !sendAll && amountStr == "" {
			fmt.Println("-aで金額を指定してください。")
			proc = false
		}
		if sendAll && amountStr != "" {
			fmt.Println("送金額と--allを同時に指定しないでください。")
			proc = false
		}
		if sendTo == "" {
			fmt.Println("-dで送金先を指定してください。")
			proc = false
		}
		if maxfeeGweiStr != "" && feeAbout != "" {
			fmt.Println("--maxfeeもしくは--aboutfeeで合計手数料を指定してください。")
			proc = false
		}
		if !proc {
			return
		}
		// -------------------------------------------------------------------
		ac := config.GetMainAccount()
		address := ac.Address
		nonce, _ := util.GetNonce(address, chainPay)
		chaininfo, err := util.GetChainInfo(chainPay)
		basefeeWei, _ := util.ReadBaseFee(chaininfo)
		to := util.NameToAddress(sendTo)
		toBytes, _ := hex.DecodeString(util.PureHex(strings.ToLower(to)))
		tx := util.NewTx()
		if token == "" {
			util.InitNativeTx(&tx, chainPay)
			if gasLimit != "" {
				glInt, _ := strconv.Atoi(gasLimit)
				tx.GasLimit = util.IntToBytes(uint64(glInt))
			}
			tx.Nonce = nonce
			tx.To = toBytes
			if err != nil {
				fmt.Println(err)
				return
			}
			if !sendAll {
				tx.Value = util.EthToWei(amountStr)
				if maxfeeGweiStr != "" {
					maxfeewei := util.EthToGwei(maxfeeGweiStr)
					if uint64(maxfeewei) < basefeeWei*21000 {
						fmt.Print("指定の最大手数料では不足です。\n現在の最低手数料: ")
						fmt.Print(basefeeWei * 21 / 1000000)
						fmt.Println(" Gwei")
						return
					}
					maxfeepergas := maxfeewei / 21000
					maxpriorityfeepargas := basefeeWei / 5
					fmt.Print("Estimated fee: ")
					fmt.Print((basefeeWei + maxpriorityfeepargas) * 21 / 1000000)
					fmt.Println(" Gwei")
					fmt.Print("Max fee:       ")
					fmt.Print(maxfeepergas * 21 / 1000000)
					fmt.Println(" Gwei")
					tx.MaxFeePerGas = util.IntToBytes(uint64(maxfeepergas))
					tx.MaxPriorityFeePerGas = util.IntToBytes(maxpriorityfeepargas)
				} else if feeAbout != "" {
					feewei := util.EthToGwei(feeAbout)
					if uint64(feewei) < basefeeWei*21000 {
						fmt.Print("指定の最大手数料では不足です。\n現在の最低手数料: ")
						fmt.Print(basefeeWei * 21 / 1000000)
						fmt.Println(" Gwei")
						return
					}
					feepergas := feewei / 21000
					fmt.Print("Fee: ")
					fmt.Print(feepergas * 21 / 1000000)
					fmt.Println(" Gwei")
					tx.MaxFeePerGas = util.IntToBytes(uint64(feepergas))
					tx.MaxPriorityFeePerGas = util.IntToBytes(uint64(feepergas))
				}
			} else {
				if feeAbout == "" {
					fmt.Println("--aboutfeeによりおよその合計手数料(Gwei)を設定してください。")
					return
				}
				feewei := util.EthToGwei(feeAbout)
				if uint64(feewei) < basefeeWei*21000 {
					fmt.Print("指定の手数料では不足です。\n現在の最低手数料: ")
					fmt.Print(basefeeWei * 21 / 1000000)
					fmt.Println(" Gwei")
					return
				}
				feepergas := feewei / 21000
				fmt.Print("Fee: ")
				fmt.Print(feepergas * 21 / 1000000)
				fmt.Println(" Gwei")
				bal, err := util.GetBalance(ac.Address, chainPay)
				if err != nil {
					fmt.Println("残高取得失敗 インターネット接続を確認してください")
					return
				}
				val := new(big.Int)
				feeBigInt := big.NewInt(int64(feepergas * 21000))
				val.Sub(bal, feeBigInt)
				tx.Value = val.Bytes()
				tx.MaxFeePerGas = util.IntToBytes(uint64(feepergas))
				tx.MaxPriorityFeePerGas = util.IntToBytes(uint64(feepergas))
			}
		} else {
			util.InitTokenTx(&tx, chainPay)
			tx.Nonce = nonce
			glInt := 100000
			if gasLimit != "" {
				glInt, _ = strconv.Atoi(gasLimit)
				tx.GasLimit = util.IntToBytes(uint64(glInt))
			}
			maxfeewei := util.EthToGwei(maxfeeGweiStr)
			if uint64(maxfeewei) < basefeeWei*uint64(glInt) {
				fmt.Print("指定の最大手数料では不足です。\n現在の最低手数料: ")
				fmt.Print(basefeeWei * uint64(glInt) / 1000000000)
				fmt.Println(" Gwei")
				return
			}
			maxfeepergas := maxfeewei / glInt
			maxpriorityfeepargas := basefeeWei / 5
			fmt.Print("Estimated fee: ")
			fmt.Print((basefeeWei + maxpriorityfeepargas) * uint64(glInt) / 1000000000)
			fmt.Println(" Gwei")
			fmt.Print("Max fee:       ")
			fmt.Print(maxfeepergas * glInt / 1000000000)
			fmt.Println(" Gwei")
			tx.MaxFeePerGas = util.IntToBytes(uint64(maxfeepergas))
			tx.MaxPriorityFeePerGas = util.IntToBytes(maxpriorityfeepargas)
			decimal, id := util.TokenInfo(token)
			tx.To = id
			dataField := []byte{}
			dataField = append(dataField, []byte{0xa9, 0x05, 0x9c, 0xbb}...)
			dataField = append(dataField, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}...)
			dataField = append(dataField, toBytes...)
			sendamount := util.DecstrToBigint(amountStr, decimal).Bytes()
			for i := 0; i < 32-len(sendamount); i++ {
				dataField = append(dataField, byte(0x00))
			}
			dataField = append(dataField, sendamount...)
		}
		util.SaveTx(tx)
	},
}

func init() {
	rootCmd.AddCommand(paymentCmd)

	paymentCmd.Flags().StringVarP(&chainPay, "chain", "c", "ETH", "")
	paymentCmd.Flags().StringVar(&gasLimit, "gaslimit", "", "")
	paymentCmd.Flags().StringVarP(&amountStr, "amount", "a", "", "")
	paymentCmd.Flags().StringVarP(&sendTo, "destination", "d", "", "")
	paymentCmd.Flags().BoolVar(&sendAll, "all", false, "")
	paymentCmd.Flags().StringVar(&maxfeeGweiStr, "maxfee", "", "")
	paymentCmd.Flags().StringVar(&feeAbout, "aboutfee", "", "")
	paymentCmd.Flags().StringVar(&token, "token", "", "")
}
