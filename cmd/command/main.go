package main

import (
	"fmt"
	"os"

	"github.com/Binaretech/classroom-auth/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Classroom Auth",
	Short: "Authentication service",
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	config.Initialize()
	execute()
}
