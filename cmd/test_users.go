//go:build !production

package cmd

import (
	"context"

	"github.com/Binaretech/classroom-auth/database"
	"github.com/Binaretech/classroom-auth/database/schema"
	"github.com/Binaretech/classroom-auth/hash"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var testUsers = &cobra.Command{
	Use:   "create:users",
	Short: "Insert users in database",
	Run: func(cmd *cobra.Command, args []string) {

		client, _ := database.Connect()
		db := database.GetDatabase(client)

		defer database.Close(client)

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

		collection := database.Users(db)
		collection.InsertMany(context.Background(), users)
	},
}

func init() {
	rootCmd.AddCommand(testUsers)
}
