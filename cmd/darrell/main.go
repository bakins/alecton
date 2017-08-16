package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "alecton",
	Short: "experimental Kubernetes deployer",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var overwriteObject = false

func addOverwriteFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&overwriteObject, "overwrite", "o", false, "overwrite existing object")
}

func getOverwriteObject() bool {
	return overwriteObject
}
