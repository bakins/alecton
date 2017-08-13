package main

import (
	"fmt"
	"log"

	"github.com/bakins/darrell/api"
	"github.com/spf13/cobra"
)

var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "work with artifacts",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var artifactListCmd = &cobra.Command{
	Use:   "list",
	Short: "list artifacts",
	Run:   runArtifactList,
}

func runArtifactList(cmd *cobra.Command, args []string) {
	c, ctx := newClient()
	fmt.Println("artifact list")
	l, err := c.ListArtifacts(ctx, &api.ListArtifactsRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(l)
}

func init() {
	addClientFlags(artifactCmd)
	artifactCmd.AddCommand(artifactListCmd)
	rootCmd.AddCommand(artifactCmd)
}
