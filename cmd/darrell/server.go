package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bakins/alecton"
	"github.com/bakins/alecton/memory"
	"github.com/spf13/cobra"
)

var serverAddress = "127.0.0.1:8080"

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start darrel server",
	Run:   runServer,
}

func runServer(cmd *cobra.Command, args []string) {
	m := memory.New()

	s, err := alecton.NewServer(
		alecton.SetStorageProvider(m),
		alecton.SetAddress(serverAddress),
	)

	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		s.Stop()
	}()

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	serverCmd.PersistentFlags().StringVarP(&serverAddress, "address", "", serverAddress, "listening address for server")
	rootCmd.AddCommand(serverCmd)
}
