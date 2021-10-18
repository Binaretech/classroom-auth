package auth_test

import (
	"testing"

	"github.com/Binaretech/classroom-auth/internal/auth"
	"github.com/Binaretech/classroom-auth/internal/cache"
	"github.com/Binaretech/classroom-auth/internal/config"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	config.Initialize()
	cache.Initialize()

	id := gofakeit.DigitN(3)
	token, err := auth.Authenticate(id)

	assert.Empty(t, err)

	_, status := auth.VerifyToken(token.AccessToken)
	assert.True(t, status)
}
