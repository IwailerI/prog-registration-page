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

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: "Creates database",
		Long:  "Creates database, used for server. If database already exists, does nothing.",
		Run: func(cmd *cobra.Command, args []string) {
			server.CreateDB()
		},
	}
	cmdCreate.PersistentFlags().VarP(&server.DBType, "database", "t", "type of database to use")

	var filename string
	cmdExport := &cobra.Command{
		Use:   "export",
		Short: "Exports database data",
		Long:  "Exports database data as csv (comma seperated value) file. Safely escaped using backslashes (\\)",
		Run: func(cmd *cobra.Command, args []string) {
			server.ExportDB(filename)
		},
	}
	cmdExport.PersistentFlags().StringVarP(&filename, "output", "o", "export.csv", "name of the generated file")
	cmdExport.PersistentFlags().VarP(&server.DBType, "database", "t", "type of database to use")
	rootCmd := &cobra.Command{Use: "server"}
	rootCmd.AddCommand(cmdRun, cmdCreate, cmdExport)
	rootCmd.Execute()
}
