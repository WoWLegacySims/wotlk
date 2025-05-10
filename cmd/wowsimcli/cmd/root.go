package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wowsimcli",
	Short: "WoWLegacySims command line tool",
	Long:  "WoWLegacySims command line tool",
}

func Execute(version string) {
	rootCmd.AddCommand(newVersionCommand(version))
	rootCmd.AddCommand(simCmd)
	rootCmd.AddCommand(bulkCmd)
	rootCmd.AddCommand(decodeLinkCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
