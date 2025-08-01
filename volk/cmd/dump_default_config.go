package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/awaisamjad/volk/config"
)

var dumpDefaultConfigCmd = &cobra.Command{
	Use:   "dumpDefaultConfig",
	Short: "Dump the default configuration to stdout",
	Long:  "This command dumps the default configuration, which includes sensible defaults and helpful comments, to standard output and exits.",
	Run:   dumpDefaultConfig,
}

func dumpDefaultConfig(cmd *cobra.Command, args []string) {
	fmt.Println(config.DefaultConfig().String())
}
