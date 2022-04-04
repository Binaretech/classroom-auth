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

	"github.com/Binaretech/classroom-auth/auth"
	"github.com/Binaretech/classroom-auth/database"
	"github.com/Binaretech/classroom-auth/database/schema"
	"github.com/Binaretech/classroom-auth/hash"
	"github.com/Binaretech/classroom-auth/server"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMain(m *testing.M) {
	defer database.Close()

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

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode, string(raw))
}

func TestLogin(t *testing.T) {
	email := createTestUser(t)
	response := login(t, email)
	assert.Equal(t, response.User.Email, email)

	_, authenticated := auth.VerifyToken(response.Token.AccessToken)
	assert.True(t, authenticated)
}

func TestVerify(t *testing.T) {
	email := createTestUser(t)
	response := login(t, email)

	req := httptest.NewRequest(fiber.MethodGet, "/auth", nil)
	req.Header.Set(fiber.HeaderAuthorization, fmt.Sprintf("Bearer: %s", response.Token.AccessToken))

	resp, _ := server.App().Test(req)

	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

func TestLogout(t *testing.T) {
	email := createTestUser(t)
	response := login(t, email)

	req := httptest.NewRequest(fiber.MethodPost, "/auth/logout", nil)
	req.Header.Set(fiber.HeaderAuthorization, fmt.Sprintf("Bearer: %s", response.Token.AccessToken))

	resp, _ := server.App().Test(req)

	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

func TestRefreshToken(t *testing.T) {
	email := createTestUser(t)
	response := login(t, email)

	body, _ := json.Marshal(map[string]string{
		"refreshToken": response.Token.RefreshToken,
	})

	req := httptest.NewRequest(fiber.MethodPost, "/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	resp, _ := server.App().Test(req)
	raw, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode, string(raw))

	response = loginResponse{}
	json.Unmarshal(raw, &response)

	assert.Equal(t, response.User.Email, email)
}

type loginResponse struct {
	Token auth.TokenDetails `json:"token"`
	User  schema.User       `json:"user"`
}

func createTestUser(t *testing.T) string {
	email := gofakeit.Email()

	upsert := true

	if _, err := database.Users().ReplaceOne(
		context.Background(),
		bson.M{
			"email": email,
		},
		bson.M{
			"email":    email,
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
	raw, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode, string(raw))

	response := loginResponse{}
	json.Unmarshal(raw, &response)

	assert.Equal(t, response.User.Email, email)

	_, authenticated := auth.VerifyToken(response.Token.AccessToken)
	assert.True(t, authenticated)

	return response
}
