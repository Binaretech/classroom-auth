//go:build production

package cmd

import (
	"fmt"

	"github.com/Binaretech/classroom-auth/database"
	"github.com/Binaretech/classroom-auth/server"
	"github.com/spf13/viper"

	_ "github.com/Binaretech/classroom-auth/config"
)

func Setup() {
	defer database.Close()

	server.App().Listen(fmt.Sprintf(":%s", viper.GetString("port")))

}
