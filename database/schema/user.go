package schema

import (
	"context"
	"time"

	"github.com/Binaretech/classroom-auth/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User schema
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"-"`
}

// Authenticate user
func (u *User) Authenticate() (*auth.TokenDetails, error) {
	return auth.Authenticate(u.ID.Hex())
}

func (u *User) Create(collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	if result, err := collection.InsertOne(ctx, u); err != nil {
		return err
	} else {
		u.ID = result.InsertedID.(primitive.ObjectID)
	}

	return nil
}

func (u *User) FindByEmail(collection *mongo.Collection, email string) error {
	ctx, abort := context.WithTimeout(context.Background(), 1*time.Minute)
	defer abort()

	filter := bson.M{
		"email": email,
	}

	return collection.FindOne(ctx, filter).Decode(u)
}
