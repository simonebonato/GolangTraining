/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "To indicate that a task had been accomplished.",
	Long:  `Takes in one or more numbers to related tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db := cmd.Context().Value("db").(*bolt.DB)

		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to parse: ", arg)
			} else {
				ids = append(ids, id)
			}
		}
		fmt.Println(ids)
		fmt.Println(db)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
