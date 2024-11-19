/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"
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

		ids := argsToIds(args)

		var taskToDelete boltDb.Task
		for _, id := range ids {
			err := db.One("Key", id, &taskToDelete)
			if err != nil {
				panic(err)
			}

			// add the task to the Done db
			taskToAdd := boltDb.DoneTask{
				DoneTask: taskToDelete,
				DoneAt:   time.Now(),
			}
			if err := taskToAdd.Validate(); err != nil {
				panic(err)
			}
			boltDb.AddTask(db, &taskToAdd)

			err = db.DeleteStruct(&taskToDelete)
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
