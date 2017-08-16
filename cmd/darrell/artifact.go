package main

import (
	"fmt"
	"log"

	"github.com/bakins/alecton/api"
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
	Use:   "list NAME",
	Short: "list artifacts",
	Run:   runArtifactList,
}

func runArtifactList(cmd *cobra.Command, args []string) {
	c, ctx := newClient()
	r := &api.ListArtifactsRequest{}
	if len(args) > 0 && args[0] != "" {
		r.Name = args[0]
	}
	l, err := c.ListArtifacts(ctx, r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(l)
}

var artifactGetCmd = &cobra.Command{
	Use:   "get NAME VERSION",
	Short: "get an artifact",
	Run:   runArtifactGet,
}

func runArtifactGet(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.Fatal("need NAME and VERSION")
	}
	c, ctx := newClient()
	a, err := c.GetArtifact(ctx, &api.GetArtifactRequest{Name: args[0], Version: args[1]})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(a)
}

var artifactCreateCmd = &cobra.Command{
	Use:   "create NAME VERSION BRANCH IMAGE",
	Short: "create artifact",
	Run:   runArtifactCreate,
}

func runArtifactCreate(cmd *cobra.Command, args []string) {
	if len(args) != 4 {
		log.Fatal("need NAME, VERSION, BRANCH, and IMAGE")
	}
	c, ctx := newClient()
	r := &api.CreateArtifactRequest{
		Artifact: &api.Artifact{
			Name:    args[0],
			Version: args[1],
			Branch:  args[2],
			Image:   args[3],
		},
		Overwrite: getOverwriteObject(),
	}
	a, err := c.CreateArtifact(ctx, r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(a)
}

func init() {
	addClientFlags(artifactCmd)
	addOverwriteFlags(artifactCmd)
	artifactCmd.AddCommand(artifactListCmd)
	artifactCmd.AddCommand(artifactGetCmd)
	artifactCmd.AddCommand(artifactCreateCmd)
	rootCmd.AddCommand(artifactCmd)
}
