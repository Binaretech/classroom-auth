package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"

	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Binaretech/classroom-auth/internal/auth"
	"github.com/Binaretech/classroom-auth/internal/cache"
	"github.com/Binaretech/classroom-auth/internal/config"
	"github.com/Binaretech/classroom-auth/internal/database"
	"github.com/Binaretech/classroom-auth/internal/database/schema"
	"github.com/Binaretech/classroom-auth/internal/hash"
	"github.com/Binaretech/classroom-auth/internal/server"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMain(m *testing.M) {
	config.Initialize()
	cache.Initialize()
	database.Connect()

	os.Exit(m.Run())
}

func TestRegister(t *testing.T) {
	collection := database.Users()

	email := gofakeit.Email()

	upsert := true
	if _, err := collection.ReplaceOne(
		context.Background(),
		bson.M{
			"email": email,
		},
		bson.M{
			"password": hash.Bcrypt("secret"),
		},
		&options.ReplaceOptions{
			Upsert: &upsert,
		},
	); err != nil {
		t.Error(err.Error())
	}

	body, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": "secret",
	})

	req := httptest.NewRequest(fiber.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	resp, _ := server.App().Test(req)

	raw, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode, string(raw))
}

func TestLogin(t *testing.T) {
	email := createTestUser(t)
	response := login(t, email)
	assert.Equal(t, response.User.Email, email)

	_, authenticated := auth.VerifyToken(response.Token)
	assert.True(t, authenticated)
}

func TestVerify(t *testing.T) {
	email := createTestUser(t)
	response := login(t, email)

	req := httptest.NewRequest(fiber.MethodGet, "/auth", nil)
	req.Header.Set(fiber.HeaderAuthorization, fmt.Sprintf("Bearer: %s", response.Token))

	resp, _ := server.App().Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

type loginResponse struct {
	Token string      `json:"token"`
	User  schema.User `json:"user"`
}

func createTestUser(t *testing.T) string {
	email := gofakeit.Email()

	upsert := true

	collection := database.Users()
	fmt.Println(collection.Database().Name())

	if _, err := database.Users().ReplaceOne(
		context.Background(),
		bson.M{
			"email": email,
		},
		bson.M{
			"password": hash.Bcrypt("secret"),
		},
		&options.ReplaceOptions{
			Upsert: &upsert,
		},
	); err != nil {
		t.Fatal(err.Error())
	}

	return email
}

func login(t *testing.T, email string) loginResponse {
	body, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": "secret",
	})

	req := httptest.NewRequest(fiber.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	server.App().Test(req)
	resp, _ := server.App().Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	raw, _ := ioutil.ReadAll(resp.Body)

	response := loginResponse{}
	json.Unmarshal(raw, &response)

	assert.Equal(t, response.User.Email, email)

	_, authenticated := auth.VerifyToken(response.Token)
	assert.True(t, authenticated)

	return response
}
