package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "volk",
	Short: "Volk is a lightweight HTTP server",
	Long:  `Volk is a lightweight HTTP server written in Go, designed to serve static files with minimal configuration.`,
}

func init() {
	rootCmd.AddCommand(serveCmd, dumpDefaultConfigCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
