package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/Binaretech/classroom-auth/auth"
	"github.com/Binaretech/classroom-auth/database"
	"github.com/Binaretech/classroom-auth/database/schema"
	"github.com/Binaretech/classroom-auth/errors"
	"github.com/Binaretech/classroom-auth/hash"
	"github.com/Binaretech/classroom-auth/lang"
	"github.com/Binaretech/classroom-auth/utils"
	"github.com/Binaretech/classroom-auth/validation"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// loginRequest is the request body for login endpoint
type loginRequest struct {
	Email    string `json:"email" validate:"required,exists=users"`
	Password string `json:"password" validate:"required"`
}

// registerRequest is the request body for register endpoint
type registerRequest struct {
	Email    string `json:"email" validate:"required,email,unique=users"`
	Password string `json:"password" validate:"required,min=6"`
}

// Register a new user and returns the login tokens
func Register(c echo.Context) error {
	req := registerRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := validation.Struct(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	collection := database.Users()

	user := schema.User{
		Email:    req.Email,
		Password: hash.Bcrypt(req.Password),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if result, err := collection.InsertOne(ctx, &user); err != nil {
		return err
	} else {
		user.ID = result.InsertedID.(primitive.ObjectID)
	}

	token, err := user.Authenticate()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"user":  user,
		"token": token,
	})
}

// Login authenticate the user and returns token data
func Login(c echo.Context) error {
	req := loginRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := validation.Struct(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	users := database.Users()

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
func Verify(c echo.Context) error {
	details, valid := auth.Verify(c)

	if !valid {
		return errors.NewUnauthenticatedError()
	}

	c.Request().Header.Add("X-User", details.UserID)
	return c.NoContent(http.StatusNoContent)
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// RefreshToken refresh the token and returns the new token
func RefreshToken(c echo.Context) error {
	req := refreshRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := validation.Struct(req); err != nil {
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

	users := database.Users()
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

// Logout the user and invalidate the token
func Logout(c echo.Context) error {
	details, valid := auth.Verify(c)

	if !valid {
		return c.NoContent(http.StatusUnauthorized)
	}

	if err := auth.DeleteAuth(details.AccessUUUID, details.RefreshUUID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
