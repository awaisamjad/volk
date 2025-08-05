package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/awaisamjad/volk/config"
)

var dumpDefaultConfigCmd = &cobra.Command{
	Use:   "dump-config",
	Short: "Dump the default configuration to stdout",
	Long:  "This command dumps the default configuration, which includes sensible defaults and helpful comments, to standard output and exits.",
	Run:   dumpConfig,
}

func dumpConfig(cmd *cobra.Command, args []string) {
	fmt.Println(config.DefaultConfig().String())
}
