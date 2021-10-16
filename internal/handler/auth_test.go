package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/Binaretech/classroom-auth/internal/auth"
	"github.com/Binaretech/classroom-auth/internal/database/schema"
	"github.com/Binaretech/classroom-auth/internal/server"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	response := login(t)
	assert.Equal(t, response.User.Email, "user")

	_, authenticated := auth.VerifyToken(response.Token)
	assert.True(t, authenticated)
}

func TestVerify(t *testing.T) {
	response := login(t)

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Header.Set(fiber.HeaderAuthorization, fmt.Sprintf("Bearer: %s", response.Token))

	resp, _ := server.App().Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

type loginResponse struct {
	Token string
	User  schema.User
}

func login(t *testing.T) loginResponse {
	body, _ := json.Marshal(map[string]string{
		"username": "user",
		"password": "secret",
	})

	req := httptest.NewRequest(fiber.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	resp, _ := server.App().Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	raw, _ := ioutil.ReadAll(resp.Body)

	response := loginResponse{}
	json.Unmarshal(raw, &response)
	return response
}
