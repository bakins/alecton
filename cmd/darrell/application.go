package main

import (
	"fmt"
	"log"

	"github.com/bakins/alecton/api"
	"github.com/spf13/cobra"
)

var applicationCmd = &cobra.Command{
	Use:   "application",
	Short: "work with applications",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var applicationListCmd = &cobra.Command{
	Use:   "list NAME",
	Short: "list applications",
	Run:   runApplicationList,
}

func runApplicationList(cmd *cobra.Command, args []string) {
	c, ctx := newClient()
	r := &api.ListApplicationsRequest{}
	l, err := c.ListApplications(ctx, r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(l)
}

var applicationGetCmd = &cobra.Command{
	Use:   "get NAME VERSION",
	Short: "get an application",
	Run:   runApplicationGet,
}

func runApplicationGet(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("need NAME")
	}
	c, ctx := newClient()
	a, err := c.GetApplication(ctx, &api.GetApplicationRequest{Name: args[0]})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(a)
}

var applicationCreateCmd = &cobra.Command{
	Use:   "create NAME CHART ARTIFACT",
	Short: "create application",
	Run:   runApplicationCreate,
}

func runApplicationCreate(cmd *cobra.Command, args []string) {
	if len(args) != 3 {
		log.Fatal("need NAME, CHART, and ARTIFACT")
	}
	c, ctx := newClient()
	r := &api.CreateApplicationRequest{
		Application: &api.Application{
			Name:     args[0],
			Chart:    args[1],
			Artifact: args[2],
		},
		Overwrite: getOverwriteObject(),
	}
	a, err := c.CreateApplication(ctx, r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(a)
}

func init() {
	addClientFlags(applicationCmd)
	addOverwriteFlags(applicationCmd)
	applicationCmd.AddCommand(applicationListCmd)
	applicationCmd.AddCommand(applicationGetCmd)
	applicationCmd.AddCommand(applicationCreateCmd)
	rootCmd.AddCommand(applicationCmd)
}
