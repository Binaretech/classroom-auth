package schema

import (
	"github.com/Binaretech/classroom-auth/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
