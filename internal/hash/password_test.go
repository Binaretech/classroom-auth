package hash_test

import (
	"testing"

	"github.com/Binaretech/classroom-auth/internal/hash"
	"github.com/stretchr/testify/assert"
)

func TestBcrypt(t *testing.T) {
	password := "secret"

	hashed := hash.Bcrypt(password)

	assert.NotEmpty(t, hashed)
	assert.True(t, hash.CompareHash(hashed, password))
}
