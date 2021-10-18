package database

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var database *mongo.Database

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

func Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return database.Collection(name, opts...)
}

func Users() *mongo.Collection {
	return database.Collection("users", &options.CollectionOptions{})
}

func Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		logrus.Println(err.Error())
	}
}
