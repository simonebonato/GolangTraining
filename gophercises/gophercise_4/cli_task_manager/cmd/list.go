/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all the available TODOs that have not been completed.",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Now it should show all the TODOs that have not been accomplished.")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
