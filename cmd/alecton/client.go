package main

import (
	"log"
	"time"

	"github.com/bakins/alecton/api"
	"github.com/spf13/cobra"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

var clientAddress = "127.0.0.1:8080"

func addClientFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&clientAddress, "server", "", clientAddress, "address of alecton server")
}

func newClient() (api.DeployServiceClient, context.Context) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := grpc.DialContext(ctx, clientAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	return api.NewDeployServiceClient(conn), ctx
}
