package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var alectonVersion = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "display the version",
	Run:   runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Println("alecton version", alectonVersion)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
