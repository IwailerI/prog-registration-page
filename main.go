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
	cmdRun.PersistentFlags().BoolVarP(&server.DebugPrint, "debug", "d", false, "whether debug information will be outputed")

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

	var filename string
	cmdExport := &cobra.Command{
		Use:   "export",
		Short: "Exports database data",
		Long:  "Exports database data as csv (comma seperated value) file. Safely escaped using backslashes (\\)",
		Run: func(cmd *cobra.Command, args []string) {
			err := database.Open()
			if err != nil {
				log.Fatal(err)
			}
			err = database.Export(filename)
			if err != nil {
				log.Fatal(err)
			}
			database.Close()
		},
	}
	cmdExport.PersistentFlags().StringVarP(&filename, "output", "o", "export.csv", "name of the generated file")

	rootCmd := &cobra.Command{Use: "server"}
	rootCmd.AddCommand(cmdRun, cmdCreate, cmdExport)
	rootCmd.Execute()
}
