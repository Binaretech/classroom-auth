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

		client, err := database.Connect()

		if err != nil {
			logrus.Fatalln(err.Error())
		}

		defer database.Close(client)

		database := database.GetDatabase(client)

		logrus.Fatalln(server.App(database).Start(fmt.Sprintf(":%s", viper.GetString("port"))))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
