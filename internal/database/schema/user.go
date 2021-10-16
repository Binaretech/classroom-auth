package schema

import (
	"github.com/Binaretech/classroom-auth/internal/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

func (u *User) Authenticate() (*auth.TokenDetails, error) {
	return auth.Authenticate(u.ID.String())
}
