package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/auth4me/pkg"
)

func AuthMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")

		if c.Get("Authorization") == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized, token not found"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized, invalid token format"})
		}

		encryptedToken := parts[1]
		user, err := pkg.ValidateToken(encryptedToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized, invalid token"})
		}

		c.Locals("userID", user.UserID)
		c.Locals("provider", user.Provider)
		c.Locals("email", user.Email)
		c.Locals("role", user.Role)
		c.Locals("sessionID", user.SessionID)
		c.Locals("mfaCompleted", user.MFACompleted)

		return c.Next()
	}
}
