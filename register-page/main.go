package main

import (
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
	cmdRun.PersistentFlags().BoolVarP(&server.DebugPrint, "debug", "d", false, "whether debug information will be outputed")
	cmdRun.PersistentFlags().VarP(&server.DBType, "database", "t", "type of database to use")
	rootCmd := &cobra.Command{Use: "server"}
	rootCmd.AddCommand(cmdRun)
	rootCmd.Execute()
}
