package utils

import "github.com/gofiber/fiber/v2"

func ResponseBadRequest(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(map[string]string{
		"message": msg,
	})
}
