package database

import (
	"context"
	"time"

	_ "github.com/Binaretech/classroom-auth/internal/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	Connect()
}

var client *mongo.Client
var database *mongo.Database

// Connect to the database
func Connect() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err = mongo.Connect(ctx,
		options.Client().
			SetHosts([]string{viper.GetString("DB_HOST")}).
			SetAuth(
				options.Credential{
					Username: viper.GetString("DB_USERNAME"),
					Password: viper.GetString("DB_PASSWORD"),
				},
			),
	)

	if err != nil {
		return
	}

	database = client.Database(viper.GetString("DB_NAME"), &options.DatabaseOptions{})
	return
}

// Collection opens a collection in the database
func Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return database.Collection(name, opts...)
}

// Users opens users collection
func Users() *mongo.Collection {
	return database.Collection("users", &options.CollectionOptions{})
}

// Close database instance
func Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		logrus.Println(err.Error())
	}
}
