package main

import (
	"fmt"

	"github.com/Binaretech/classroom-auth/internal/cache"
	"github.com/Binaretech/classroom-auth/internal/config"
	"github.com/Binaretech/classroom-auth/internal/database"
	"github.com/Binaretech/classroom-auth/internal/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	config.Initialize()

	database.Connect()

	defer database.Close()

	if err := cache.Initialize(); err != nil {
		logrus.Info(err)
	}

	server.App().Listen(fmt.Sprintf(":%s", viper.GetString("port")))

}
