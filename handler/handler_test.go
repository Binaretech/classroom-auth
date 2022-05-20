package handler_test

import (
	"bytes"
	"context"
	"encoding/json"

	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Binaretech/classroom-auth/auth"
	"github.com/Binaretech/classroom-auth/database"
	"github.com/Binaretech/classroom-auth/database/schema"
	"github.com/Binaretech/classroom-auth/handler"
	"github.com/Binaretech/classroom-auth/hash"
	"github.com/Binaretech/classroom-auth/validation"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createContext(req *http.Request, rec *httptest.ResponseRecorder, db *mongo.Database) echo.Context {
	app := echo.New()
	app.Validator = validation.SetUpValidator(db)

	return app.NewContext(req, rec)
}

func TestRegister(t *testing.T) {
	client, _ := database.Connect()
	defer database.Close(client)

	db := database.GetDatabase(client)

	email := gofakeit.Email()

	body, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": "secret",
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := createContext(req, rec, db)

	h := handler.NewHandler(db)

	if assert.NoError(t, h.Register(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestLogin(t *testing.T) {
	email := createTestUser(t)
	response := login(t, email)
	assert.Equal(t, response.User.Email, email)

	_, authenticated := auth.VerifyToken(response.Token.AccessToken)
	assert.True(t, authenticated)
}

func TestVerify(t *testing.T) {
	client, _ := database.Connect()
	defer database.Close(client)

	db := database.GetDatabase(client)

	email := createTestUser(t)
	response := login(t, email)

	req := httptest.NewRequest(http.MethodGet, "/auth", nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer: %s", response.Token.AccessToken))

	rec := httptest.NewRecorder()

	c := createContext(req, rec, db)

	h := handler.NewHandler(db)

	if assert.NoError(t, h.Verify(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}

func TestLogout(t *testing.T) {
	client, _ := database.Connect()
	defer database.Close(client)

	db := database.GetDatabase(client)

	email := createTestUser(t)
	response := login(t, email)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer: %s", response.Token.AccessToken))

	rec := httptest.NewRecorder()

	c := createContext(req, rec, db)

	h := handler.NewHandler(db)

	if assert.NoError(t, h.Logout(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}

func TestRefreshToken(t *testing.T) {
	client, _ := database.Connect()
	defer database.Close(client)

	db := database.GetDatabase(client)

	email := createTestUser(t)
	response := login(t, email)

	body, _ := json.Marshal(map[string]string{
		"refreshToken": response.Token.RefreshToken,
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	h := handler.NewHandler(db)

	c := createContext(req, rec, db)

	if assert.NoError(t, h.RefreshToken(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		response = loginResponse{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, response.User.Email, email)
	}

}

type loginResponse struct {
	Token auth.TokenDetails `json:"token"`
	User  schema.User       `json:"user"`
}

func createTestUser(t *testing.T) string {
	client, _ := database.Connect()
	defer database.Close(client)

	db := database.GetDatabase(client)

	email := gofakeit.Email()

	upsert := true

	if _, err := database.Users(db).ReplaceOne(
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
	client, _ := database.Connect()
	defer database.Close(client)

	db := database.GetDatabase(client)

	body, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": "secret",
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := createContext(req, rec, db)

	h := handler.NewHandler(db)
	assert.NoError(t, h.Login(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	response := loginResponse{}
	json.Unmarshal(rec.Body.Bytes(), &response)

	assert.Equal(t, response.User.Email, email)
	_, authenticated := auth.VerifyToken(response.Token.AccessToken)
	assert.True(t, authenticated)

	return response
}
