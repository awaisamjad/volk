package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "volk",
	Short: "Volk is a lightweight HTTP server",
	Long:  `Volk is a lightweight HTTP server written in Go, designed to serve static files with minimal configuration.`,
}

func init() {
	rootCmd.AddCommand(versionCmd, serveCmd, dumpDefaultConfigCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}
