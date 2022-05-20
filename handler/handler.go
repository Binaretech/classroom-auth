package handler

import (
	"context"
	"net/http"

	"github.com/Binaretech/classroom-auth/auth"
	"github.com/Binaretech/classroom-auth/database"
	"github.com/Binaretech/classroom-auth/database/schema"
	"github.com/Binaretech/classroom-auth/errors"
	"github.com/Binaretech/classroom-auth/hash"
	"github.com/Binaretech/classroom-auth/lang"
	"github.com/Binaretech/classroom-auth/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	db *mongo.Database
}

func NewHandler(db *mongo.Database) *AuthHandler {
	return &AuthHandler{
		db: db,
	}
}

// Register a new user and returns the login tokens
func (h *AuthHandler) Register(c echo.Context) error {
	req := registerRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	collection := database.Users(h.db)

	user := schema.User{
		Email:    req.Email,
		Password: hash.Bcrypt(req.Password),
	}

	if err := user.Create(collection); err != nil {
		return err
	}

	if token, err := user.Authenticate(); err != nil {
		return err
	} else {
		return c.JSON(http.StatusCreated, echo.Map{
			"user":  user,
			"token": token,
		})
	}

}

// Login authenticate the user and returns token data
func (h *AuthHandler) Login(c echo.Context) error {
	req := loginRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	users := database.Users(h.db)

	var user schema.User

	if err := users.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user); err != nil {
		return err
	}

	if !hash.CompareHash(user.Password, req.Password) {
		return utils.ResponseError(c, http.StatusUnauthorized, lang.Trans("login error"))
	}

	token, err := user.Authenticate()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{
		"user":  user,
		"token": token,
	})
}

// Verify the auth status
func (h *AuthHandler) Verify(c echo.Context) error {
	details, valid := auth.Verify(c)

	if !valid {
		return errors.NewUnauthenticatedError()
	}

	c.Response().Header().Add("X-User", details.UserID)
	return c.NoContent(http.StatusNoContent)
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// RefreshToken refresh the token and returns the new token
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	req := refreshRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	userID, valid := auth.VerifyRefreshToken(req.RefreshToken)
	if !valid {
		return c.NoContent(http.StatusUnauthorized)
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	users := database.Users(h.db)
	var user schema.User
	if err := users.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user); err != nil {
		return err
	}

	token, err := user.Authenticate()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":  user,
		"token": token,
	})

}

func (h *AuthHandler) GoogleAuth(c echo.Context) error {
	req := googleAuthRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	var email string
	var id string

	if info, err := auth.GoogleAuth(req.IdToken); err != nil {
		return err
	} else {
		email = info.Email
		id = info.UserId
	}

	users := database.Users(h.db)
	user := new(schema.User)

	if err := user.FindByEmail(users, email); err == mongo.ErrNoDocuments {
		user = &schema.User{
			Email:    email,
			Password: hash.Bcrypt(id),
		}

		if err := user.Create(users); err != nil {
			return err
		}

	} else if err != nil {
		return err
	}

	if token, err := user.Authenticate(); err != nil {
		return err
	} else {
		return c.JSON(http.StatusOK, map[string]any{
			"user":  user,
			"token": token,
		})
	}

}

// Logout the user and invalidate the token
func (h *AuthHandler) Logout(c echo.Context) error {
	details, valid := auth.Verify(c)

	if !valid {
		return c.NoContent(http.StatusUnauthorized)
	}

	if err := auth.DeleteAuth(details.AccessUUUID, details.RefreshUUID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
