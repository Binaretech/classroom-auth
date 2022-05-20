package database

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect to the database
func Connect() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return mongo.Connect(ctx,
		options.Client().
			SetHosts([]string{viper.GetString("DB_HOST")}).
			SetAuth(
				options.Credential{
					Username: viper.GetString("DB_USERNAME"),
					Password: viper.GetString("DB_PASSWORD"),
				},
			),
	)
}

func GetDatabase(client *mongo.Client) *mongo.Database {
	return client.Database(viper.GetString("DB_NAME"), &options.DatabaseOptions{})

}

// Collection opens a collection in the database
func Collection(database *mongo.Database, name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return database.Collection(name, opts...)
}

// Users opens users collection
func Users(database *mongo.Database) *mongo.Collection {
	return database.Collection("users", &options.CollectionOptions{})
}

// Close database instance
func Close(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		logrus.Println(err.Error())
	}
}
