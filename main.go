package main

import (
	"log"
	"web-server/database"
	"web-server/server"

	"github.com/spf13/cobra"
)

func main() {
	cmdRun := &cobra.Command{
		Use:   "run",
		Short: "Start the server",
		Long:  "Starts the server, initializes or creates underlaying database.",
		Run: func(cmd *cobra.Command, args []string) {
			server.Start()
		},
	}
	cmdRun.PersistentFlags().StringVarP(&server.Port, "port", "p", "8080", "port, on which server will run")

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: "Creates database",
		Long:  "Creates database, used for server. If database already exists, does nothing.",
		Run: func(cmd *cobra.Command, args []string) {
			err := database.Open()
			if err != nil {
				log.Fatal(err)
			}
			err = database.Create()
			if err != nil {
				log.Fatal(err)
			}
			database.Close()
		},
	}

	rootCmd := &cobra.Command{Use: "server"}
	rootCmd.AddCommand(cmdRun, cmdCreate)
	rootCmd.Execute()
}
