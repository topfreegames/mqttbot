package cmd

import (
	"github.com/spf13/cobra"
	"github.com/topfreegames/mqttbot/app"
)

var host string
var port int
var debug bool

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts mqttbot server",
	Long:  `Starts mqtt server with the specified arguments. You can use environment variables to override configuration keys.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := app.GetApp(
			host,
			port,
			debug,
		)
		app.Start()
	},
}

func init() {
	RootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&host, "bind", "b", "0.0.0.0", "Host to bind mqttbot to")
	startCmd.Flags().IntVarP(&port, "port", "p", 8888, "Port to bind mqttbot to")
	startCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Debug mode")
}
