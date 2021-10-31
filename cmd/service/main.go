package main

import (
	"fmt"

	"github.com/Binaretech/classroom-auth/internal/database"
	"github.com/Binaretech/classroom-auth/internal/server"
	"github.com/spf13/viper"
)

func main() {
	defer database.Close()

	server.App().Listen(fmt.Sprintf(":%s", viper.GetString("port")))

}
