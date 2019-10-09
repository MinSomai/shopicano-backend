package cmd

import (
	"fmt"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/config"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/spf13/cobra"
	"os"
)

var (
	// RootCmd is the root command of nur service
	RootCmd = &cobra.Command{
		Use:   "fin",
		Short: "A http service",
		Long:  `A HTTP JSON API backend service`,
	}
)

func init() {
	RootCmd.AddCommand(migrationCmd)
	RootCmd.AddCommand(serveCmd)
}

// Execute executes the root command
func Execute() {
	if err := config.LoadConfig(); err != nil {
		fmt.Println("Failed to read config : ", err)
		os.Exit(1)
	}
	log.SetupLog()

	if err := app.ConnectSQLDB(); err != nil {
		log.Log().Fatalf("Failed to connect to postgres : %v\n", err)
		os.Exit(-1)
	}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
