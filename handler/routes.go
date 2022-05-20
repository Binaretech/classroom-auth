package handler

import "github.com/labstack/echo/v4"

func (handler *AuthHandler) Routes(auth *echo.Group) {
	auth.POST("/google", handler.GoogleAuth)

	auth.POST("/login", handler.Login)
	auth.POST("/register", handler.Register)
	auth.POST("/refresh", handler.RefreshToken)
	auth.POST("/logout", handler.Logout)
	auth.GET("/", handler.Verify)

}
