package handler

import (
	"context"
	"time"

	"github.com/Binaretech/classroom-auth/internal/auth"
	"github.com/Binaretech/classroom-auth/internal/database"
	"github.com/Binaretech/classroom-auth/internal/database/schema"
	"github.com/Binaretech/classroom-auth/internal/hash"
	"github.com/Binaretech/classroom-auth/internal/validation"
	"github.com/gofiber/fiber/v2"
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
func Register(c *fiber.Ctx) error {
	req := registerRequest{}
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := validation.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err)
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

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user":  user,
		"token": token,
	})
}

// Login authenticate the user and returns token data
func Login(c *fiber.Ctx) error {
	req := loginRequest{}

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := validation.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err)
	}

	users := database.Users()

	var user schema.User

	if err := users.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user); err != nil {
		return err
	}

	if !hash.CompareHash(user.Password, req.Password) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	token, err := user.Authenticate()
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

// Verify the auth status
func Verify(c *fiber.Ctx) error {
	details, valid := auth.Verify(c)

	if !valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Append("X-User", details.UserID)

	return c.SendStatus(fiber.StatusNoContent)
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// RefreshToken refresh the token and returns the new token
func RefreshToken(c *fiber.Ctx) error {
	req := refreshRequest{}

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := validation.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err)
	}

	userID, valid := auth.VerifyRefreshToken(req.RefreshToken)
	if !valid {
		return c.SendStatus(fiber.StatusUnauthorized)
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

	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"user":  user,
		"token": token,
	})

}

// Logout the user and invalidate the token
func Logout(c *fiber.Ctx) error {
	details, valid := auth.Verify(c)

	if !valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if err := auth.DeleteAuth(details.AccessUUUID, details.RefreshUUID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
