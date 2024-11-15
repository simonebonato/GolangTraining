/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"
	boltDb "todo/boltDB"

	"github.com/asdine/storm"
	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "To indicate that a task had been accomplished.",
	Long:  `Takes in one or more numbers to related tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db := cmd.Context().Value("db").(*storm.DB)

		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to parse: ", arg)
			} else {
				ids = append(ids, id)
			}
		}
		
		for _, id := range ids {
			taskToDelete := boltDb.Task{Key: id}
			err := db.DeleteStruct(&taskToDelete)
			if err != nil {
				fmt.Printf("Failed to delete task with id: %v", id)
			}
		}


		return nil
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
