package main

import (
	"github.com/spf13/cobra"
	"github.com/cuongvng/blockchain-bar/database"
	"fmt"
	"os"
)

func balancesCmd() *cobra.Command {
	var balanceCmd = &cobra.Command{
		Use:   "balances",
		Short: "Interact with balances (list...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}
	balanceCmd.AddCommand(listBalancesCmd())
	return balanceCmd
}

func listBalancesCmd() *cobra.Command {
	var result = &cobra.Command{
		Use:   "list",
		Short: "Lists all balances.",
		Run: func(cmd *cobra.Command, args []string) {
			state, err := database.GetStateFromDisk()
			if err != nil {
				fmt.Println(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()

			fmt.Println("Accounts balances:")
			fmt.Println("__________________")
			fmt.Println("")
			for account, balance := range state.Balances {
				fmt.Println(fmt.Sprintf("%s: %d", account, balance))
			}
		},
	}
	return result
}