package main

import (
	"fmt"
	"log"

	"github.com/bakins/alecton/api"
	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "work with images",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var imageListCmd = &cobra.Command{
	Use:   "list NAME",
	Short: "list images",
	Run:   runImageList,
}

func runImageList(cmd *cobra.Command, args []string) {
	c, ctx := newClient()
	r := &api.ListImagesRequest{}
	if len(args) > 0 && args[0] != "" {
		r.Name = args[0]
	}
	l, err := c.ListImages(ctx, r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(l)
}

var imageCreateCmd = &cobra.Command{
	Use:   "create NAME VERSION BRANCH IMAGE",
	Short: "create image",
	Run:   runImageCreate,
}

func runImageCreate(cmd *cobra.Command, args []string) {
	if len(args) != 3 {
		log.Fatal("need NAME, VERSION, and IMAGE")
	}

	c, ctx := newClient()
	i := &api.Image{
		Name:    args[0],
		Version: args[1],
		Image:   args[2],
	}

	i, err := c.CreateImage(ctx, i)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(i)
}

func init() {
	addClientFlags(imageCmd)
	imageCmd.AddCommand(imageListCmd)
	imageCmd.AddCommand(imageCreateCmd)
	rootCmd.AddCommand(imageCmd)
}
