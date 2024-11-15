/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	boltDb "todo/boltDB"

	"github.com/asdine/storm"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all the available TODOs that have not been completed.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		
		db := cmd.Context().Value("db").(*storm.DB)

		var tasks boltDb.TaskSlice
		err := db.All(&tasks)
		if err != nil {
			log.Fatal(err)
		}
		
		fmt.Println(tasks)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
