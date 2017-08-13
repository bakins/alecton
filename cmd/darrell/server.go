package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bakins/darrell"
	"github.com/bakins/darrell/memory"
	"github.com/spf13/cobra"
	context "golang.org/x/net/context"
)

var serverAddress = "127.0.0.1:8080"

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start darrel server",
	Run:   runServer,
}

func runServer(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := memory.New()

	s := darrell.NewServer(ctx, m)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		s.Stop()
	}()

	if err := s.Run(serverAddress); err != nil {
		log.Fatal(err)
	}
}

func init() {
	serverCmd.PersistentFlags().StringVarP(&serverAddress, "address", "", serverAddress, "listening address for server")
	rootCmd.AddCommand(serverCmd)
}
