package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"volk/config"
)

var dumpConfigCmd = &cobra.Command{
	Use:   "dumpconfig",
	Short: "Dump the current configuration to stdout",
	Long:  "This command dumps the current configuration, including defaults, to standard output and exits.",
	Run:   dumpConfig,
}

func dumpConfig(cmd *cobra.Command, args []string) {
	fmt.Println(config.DefaultConfig().String())
}
