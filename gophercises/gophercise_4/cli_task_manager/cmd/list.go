/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
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

		done, err := cmd.Flags().GetBool("done")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		db := cmd.Context().Value("db").(*storm.DB)

		// to print the tasks to do
		if !done {
			boltDb.PrintToDos(db)
		} else {
			boltDb.PrintDoneToDos(db)
		}
	},
}

func init() {
	listCmd.Flags().BoolP("done", "d", false, "Print the tasks that need to be done (done=false) or that were already done (done=true)")
	rootCmd.AddCommand(listCmd)

}
