/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"
	boltDb "todo/boltDB"

	"github.com/asdine/storm"
	"github.com/spf13/cobra"
)

type DbVars struct {
	load_folder string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		DbConfig := DbVars{
			load_folder: "boltDb",
		}

		db := boltDb.CreateDb(DbConfig.load_folder)
		cmd.SetContext(context.WithValue(cmd.Context(), "db", db))
		boltDb.ResetTaskKeys(db)
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		db := cmd.Context().Value("db").(*storm.DB)
		db.Close()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.todo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
