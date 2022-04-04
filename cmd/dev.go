package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Binaretech/classroom-auth/config"
	"github.com/Binaretech/classroom-auth/database"
	"github.com/Binaretech/classroom-auth/database/schema"
	"github.com/Binaretech/classroom-auth/hash"
	"github.com/Binaretech/classroom-auth/server"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"

	_ "github.com/Binaretech/classroom-auth/config"
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

func Setup() {
	config.Initialize()
	execute()
}

var testUsers = &cobra.Command{
	Use:   "create:users",
	Short: "Insert users in database",
	Run: func(cmd *cobra.Command, args []string) {

		id, _ := primitive.ObjectIDFromHex("61a406ea18f8a0bdf663e144")

		users := []interface{}{
			schema.User{
				ID:       id,
				Email:    "test@classroom.com",
				Password: hash.Bcrypt("secret"),
			},
		}

		for i := 0; i < 10; i++ {
			users = append(users, schema.User{
				Email:    gofakeit.Email(),
				Password: hash.Bcrypt("secret"),
			})
		}

		collection := database.Users()
		collection.InsertMany(context.Background(), users)
	},
}

var serve = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		defer database.Close()

		server.App().Listen(fmt.Sprintf(":%s", viper.GetString("port")))
	},
}

func init() {
	rootCmd.AddCommand(testUsers)
	rootCmd.AddCommand(serve)
}
