package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bakins/alecton"
	"github.com/spf13/cobra"

	// load providers
	_ "github.com/bakins/alecton/providers/chart/chartdir"
	_ "github.com/bakins/alecton/providers/deploy/helm"
	_ "github.com/bakins/alecton/providers/deploy/mock"
	_ "github.com/bakins/alecton/providers/storage/memory"
)

var serverAddress = "127.0.0.1:8080"

var serverCmd = &cobra.Command{
	Use:   "server CONFIG",
	Short: "start alecton server",
	Run:   runServer,
}

func runServer(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "config file is required")
		os.Exit(1)
	}

	s, err := alecton.NewServerFromConfigFile(args[0])

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		s.Stop()
	}()

	if err := s.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}

}

func init() {
	rootCmd.AddCommand(serverCmd)
}
