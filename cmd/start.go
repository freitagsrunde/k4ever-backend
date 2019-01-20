package cmd

import (
	"fmt"

	"github.com/freitagsrunde/k4ever-backend/internal/context"
	"github.com/freitagsrunde/k4ever-backend/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the server on port 8080",
	Long: `This starts the k4ever webserver on the 
specified port
The default port is port 8080`,
	Run: startServer,
}
var port int

func init() {
	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to start the server on")
	rootCmd.AddCommand(startCmd)
	viper.BindPFlag("port", startCmd.Flags().Lookup("port"))
	viper.SetDefault("port", 8080)
}

func startServer(cmd *cobra.Command, args []string) {
	fmt.Printf("Starting web server on port %d...\n", port)

	config := context.NewConfig()

	config.MigrateDB()

	server.Start(config)

	fmt.Println("Done.")
}
