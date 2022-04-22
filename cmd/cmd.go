package cmd

import (
	"fmt"
	"os"

	"github.com/Binaretech/classroom-auth/config"
	"github.com/Binaretech/classroom-auth/database"
	"github.com/Binaretech/classroom-auth/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/Binaretech/classroom-auth/config"
)

var rootCmd = &cobra.Command{
	Use:   "Classroom Auth",
	Short: "Authentication service",
	Run: func(cmd *cobra.Command, args []string) {
		config.Initialize()
		defer database.Close()

		logrus.Fatalln(server.App().Start(fmt.Sprintf(":%s", viper.GetString("port"))))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
