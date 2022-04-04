package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Binaretech/classroom-auth/cache"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// TokenDetails Access and refresh token information
type TokenDetails struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`

	// Unique identifier for access token
	AccessUUID string `json:"accessUUID"`

	// Unique identifier for refresh token
	RefreshUUID string `json:"refreshUUID"`

	// AccessExpires Access token expiration
	AccessExpires int64 `json:"accessExpires"`

	// RefreshExpires Refresh token expiration
	RefreshExpires int64 `json:"refreshExpires"`
}

type AccessDetails struct {
	AccessUUUID string
	RefreshUUID string
	UserID      string
}

// Authenticate a user creating authentication tokens and storing it in cache
func Authenticate(userID string) (token *TokenDetails, err error) {
	if token, err = createToken(userID); err != nil {
		return
	}

	err = createAuth(userID, token)
	return
}

// VerifyToken Verify a token and return the auth state
func VerifyToken(tokenString string) (*AccessDetails, bool) {
	token, err := extractTokenMetadata(tokenString, viper.GetString("secret"))
	if err != nil {
		return nil, false
	}

	_, err = fetchAuth(token)
	return token, err == nil
}

// Verify if the token is valid
func Verify(c *fiber.Ctx) (*AccessDetails, bool) {
	return VerifyToken(extractToken(c))
}

func VerifyRefreshToken(tokenString string) (string, bool) {
	token, err := verifyToken(tokenString, viper.GetString("secret"))
	if err != nil {
		return "", false
	}

	if !token.Valid {
		return "", false
	}

	claims := token.Claims.(jwt.MapClaims)
	refreshUUID := claims["refresh_uuid"].(string)
	userID := claims["user_id"].(string)

	accessUUID, err := fetchTokenData(refreshUUID)

	if err != nil {
		return "", false
	}

	return userID, DeleteAuth(refreshUUID, accessUUID) == nil
}

// createToken Create a new token
func createToken(userID string) (td *TokenDetails, err error) {

	td = &TokenDetails{
		AccessExpires: time.Now().Add(time.Hour * 24 * 15).Unix(),
		AccessUUID:    uuid.New().String(),

		RefreshExpires: time.Now().Add(time.Hour * 24 * 7).Unix(),
		RefreshUUID:    uuid.New().String(),
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"refresh_uuid": td.RefreshUUID,
		"user_id":      userID,
		"exp":          td.RefreshExpires,
	})

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized":   true,
		"access_uuid":  td.AccessUUID,
		"refresh_uuid": td.RefreshUUID,
		"user_id":      userID,
		"exp":          td.AccessExpires,
	})

	if td.AccessToken, err = at.SignedString([]byte(viper.GetString("secret"))); err != nil {
		return
	}

	td.RefreshToken, err = rt.SignedString([]byte(viper.GetString("secret")))
	return
}

// createAuth Store the authentication state in cache
func createAuth(userID string, td *TokenDetails) error {
	at := time.Unix(td.AccessExpires, 0)
	rt := time.Unix(td.RefreshExpires, 0)
	now := time.Now()

	if _, err := cache.Set(context.Background(), td.AccessUUID, userID, at.Sub(now)); err != nil {
		return err
	}

	if _, err := cache.Set(context.Background(), td.RefreshUUID, td.AccessUUID, rt.Sub(now)); err != nil {
		return err
	}

	return nil
}

// extractToken Extract the token from the request
func extractToken(c *fiber.Ctx) string {
	token := strings.Split(c.Get("Authorization"), " ")

	if len(token) == 2 {
		return token[1]
	}

	return ""
}

// verifyToken Verify the token string
func verifyToken(token, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})
}

// extractTokenMetadata Extract the token metadata from the token string
func extractTokenMetadata(tokenString, secret string) (td *AccessDetails, err error) {
	var token *jwt.Token
	if token, err = verifyToken(tokenString, secret); err != nil {
		return nil, err
	}

	if !token.Valid {
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	accessUuid := claims["access_uuid"].(string)
	refreshUuid := claims["refresh_uuid"].(string)

	return &AccessDetails{
		AccessUUUID: accessUuid,
		RefreshUUID: refreshUuid,
		UserID:      claims["user_id"].(string),
	}, nil

}

// fetchAuth Fetch the authentication state from cache
func fetchAuth(authD *AccessDetails) (string, error) {
	return fetchTokenData(authD.AccessUUUID)
}

// fetchTokenData Fetch the token data from cache
func fetchTokenData(tokenUUID string) (string, error) {
	return cache.Get(context.Background(), tokenUUID)
}

// DeleteAuth Delete the authentication state from cache
func DeleteAuth(tokenUUID ...string) error {

	_, err := cache.Delete(context.Background(), tokenUUID...)
	return err
}
