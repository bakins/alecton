package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var darrelVersion = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "display the version",
	Run:   runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Println("darrel version", darrelVersion)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
