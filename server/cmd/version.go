package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var (
    // These would be set during build time using ldflags
    version = "dev"
    commit  = "none"
    date    = "unknown"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print the version number of Volk",
    Long:  `All software has versions. This is Volk's.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Volk version %s (commit: %s, built on: %s)\n", version, commit, date)
    },
}


