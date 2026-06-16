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
		if !proc {
			return
		}
		var glInt int
		// -------------------------------------------------------------------
		ac := config.GetMainAccount()
		chainPay := ac.Chain
		address := ac.Address
		nonce, _ := util.GetNonce(address, chainPay)
		basefeeWei, err := util.GetBaseFee(chainPay)
		if err != nil {
			fmt.Println(err)
			return
		}
		feeWeiToCalc := basefeeWei
		if basefeeWei == 0 {
			feeWeiToCalc, err = util.GetGasPrice(chainPay)
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		to := util.NameToAddress(sendTo)
		toBytes, _ := hex.DecodeString(util.PureHex(strings.ToLower(to)))
		tx := util.NewTx()
		if token == "" {
			util.InitNativeTx(&tx, chainPay)
			glInt = 21000
			if gasLimit != "" {
				gl, _ := strconv.Atoi(gasLimit)
				glInt = gl
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
					if uint64(maxfeewei) < feeWeiToCalc*uint64(glInt) {
						fmt.Print("指定の最大手数料では不足です。\n現在の最低手数料: ")
						fmt.Print(feeWeiToCalc * uint64(glInt) / 1000000000)
						fmt.Println(" Gwei")
						return
					}
					maxfeepergas := maxfeewei / glInt
					maxpriorityfeepergas := basefeeWei / 5
					if basefeeWei == 0 {
						maxpriorityfeepergas = uint64(maxfeepergas)
					}
					fmt.Print("Estimated fee: ")
					fmt.Print((basefeeWei + maxpriorityfeepergas) * uint64(glInt) / 1000000000)
					fmt.Println(" Gwei")
					fmt.Print("Max fee:       ")
					fmt.Print(maxfeepergas * glInt / 1000000000)
					fmt.Println(" Gwei")
					tx.MaxFeePerGas = util.IntToBytes(uint64(maxfeepergas))
					tx.MaxPriorityFeePerGas = util.IntToBytes(maxpriorityfeepergas)
				} else if feeAbout != "" {
					feewei := util.EthToGwei(feeAbout)
					if uint64(feewei) < feeWeiToCalc*uint64(glInt) {
						fmt.Print("指定の最大手数料では不足です。\n現在の最低手数料: ")
						fmt.Print(feeWeiToCalc * uint64(glInt) / 1000000000)
						fmt.Println(" Gwei")
						return
					}
					feepergas := feewei / glInt
					fmt.Print("Fee: ")
					fmt.Print(feepergas * glInt / 1000000000)
					fmt.Println(" Gwei")
					tx.MaxFeePerGas = util.IntToBytes(uint64(feepergas))
					tx.MaxPriorityFeePerGas = util.IntToBytes(uint64(feepergas))
				} else {
					feepergas := feeWeiToCalc
					maxpriorityfeepergas := feepergas / 10
					maxfeepergas := feepergas * 2
					if basefeeWei == 0 {
						maxpriorityfeepergas = uint64(maxfeepergas)
					}
					fmt.Print("Estimated fee: ")
					fmt.Print((basefeeWei + uint64(maxpriorityfeepergas)) * uint64(glInt) / 1000000000)
					fmt.Println(" Gwei")
					fmt.Print("Max fee:       ")
					fmt.Print(maxfeepergas * uint64(glInt) / 1000000000)
					fmt.Println(" Gwei")
					tx.MaxFeePerGas = util.IntToBytes(uint64(maxfeepergas))
					tx.MaxPriorityFeePerGas = util.IntToBytes(uint64(maxpriorityfeepergas))
				}
			} else {
				feewei := util.EthToGwei(feeAbout)
				if uint64(feewei) < feeWeiToCalc*uint64(glInt) {
					fmt.Print("指定の手数料では不足です。\n現在の最低手数料: ")
					fmt.Print(feeWeiToCalc * uint64(glInt) / 1000000000)
					fmt.Println(" Gwei")
					return
				}
				feepergas := feewei / glInt
				fmt.Print("Fee: ")
				fmt.Print(feepergas * glInt / 1000000000)
				fmt.Println(" Gwei")
				bal, err := util.GetBalance(ac.Address, chainPay)
				if err != nil {
					fmt.Println("残高取得失敗 インターネット接続を確認してください")
					return
				}
				val := new(big.Int)
				feeBigInt := big.NewInt(int64(feepergas * glInt))
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
				gl, _ := strconv.Atoi(gasLimit)
				glInt = gl
				tx.GasLimit = util.IntToBytes(uint64(glInt))
			}
			maxfeewei := feeWeiToCalc * uint64(glInt*2)
			if maxfeeGweiStr != "" {
				maxfeewei = uint64(util.EthToGwei(maxfeeGweiStr))
				if maxfeewei < feeWeiToCalc*uint64(glInt) {
					fmt.Print("指定の最大手数料では不足です。\n現在の最低手数料: ")
					fmt.Print(feeWeiToCalc * uint64(glInt) / 1000000000)
					fmt.Println(" Gwei")
					return
				}
			}
			maxfeepergas := maxfeewei / uint64(glInt)
			maxpriorityfeepergas := basefeeWei / 5
			if basefeeWei == 0 {
				maxpriorityfeepergas = maxfeepergas
			}
			fmt.Print("Estimated fee: ")
			fmt.Print((basefeeWei + maxpriorityfeepergas) * uint64(glInt) / 1000000000)
			fmt.Println(" Gwei")
			fmt.Print("Max fee:       ")
			fmt.Print(maxfeepergas * uint64(glInt) / 1000000000)
			fmt.Println(" Gwei")
			tx.MaxFeePerGas = util.IntToBytes(uint64(maxfeepergas))
			tx.MaxPriorityFeePerGas = util.IntToBytes(maxpriorityfeepergas)
			decimal, id := util.TokenInfo(chainPay, token)
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
			tx.Data = dataField
		}
		util.SaveTx(tx)
	},
}

func init() {
	rootCmd.AddCommand(paymentCmd)

	paymentCmd.Flags().StringVar(&gasLimit, "gaslimit", "", "")
	paymentCmd.Flags().StringVarP(&amountStr, "amount", "a", "", "")
	paymentCmd.Flags().StringVarP(&sendTo, "destination", "d", "", "")
	paymentCmd.Flags().BoolVar(&sendAll, "all", false, "")
	paymentCmd.Flags().StringVar(&maxfeeGweiStr, "maxfee", "", "")
	paymentCmd.Flags().StringVar(&feeAbout, "aboutfee", "", "")
	paymentCmd.Flags().StringVar(&token, "token", "", "")
}
