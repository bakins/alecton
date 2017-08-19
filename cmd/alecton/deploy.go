package main

import (
	"fmt"
	"log"

	"github.com/bakins/alecton/api"
	"github.com/spf13/cobra"
)

// todo values options
var deployCmd = &cobra.Command{
	Use:   "deploy NAME TARGET VERSION",
	Short: "deploy an application",
	Run:   runDeployCmd,
}

func runDeployCmd(cmd *cobra.Command, args []string) {
	if len(args) != 3 {
		log.Fatal("need NAME, TARGET, and VERSION")
	}
	c, ctx := newClient()

	req := &api.DeployRequest{
		Application: args[0],
		Target:      args[1],
		Version:     args[2],
	}

	release, err := c.DeployApplication(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(release)
}

func init() {
	addClientFlags(deployCmd)
	rootCmd.AddCommand(deployCmd)
}
