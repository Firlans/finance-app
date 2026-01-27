package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func AuthMiddleware(cfg *viper.Viper) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil Header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":      "Missing Authorization header",
				"request_id": c.Locals("request_id"),
			})
		}

		// 2. Format Harus "Bearer <Token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":      "Invalid Authorization Format",
				"request_id": c.Locals("request_id"),
			})
		}
		tokenString := parts[1]

		// 3. Check if token is blacklisted
		if GetBlacklist().IsBlacklisted(tokenString) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":      "Token has been revoked",
				"request_id": c.Locals("request_id"),
			})
		}

		// 4. Parse & Validasi Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			secret := cfg.GetString("jwt.secret")
			return []byte(secret), nil
		})

		// 5. Cek Error Parse
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":      "Invalid or Expired token",
				"request_id": c.Locals("request_id"),
			})
		}

		// 6. Ekstrak Claims (Data User)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":      "Invalid Token Claims",
				"request_id": c.Locals("request_id"),
			})
		}

		// 7. Simpan User ID & token ke Context
		c.Locals("user_id", claims["sub"])
		c.Locals("email", claims["email"])
		c.Locals("access_token", tokenString) // For logout

		// Store expiry for blacklist
		if exp, ok := claims["exp"].(float64); ok {
			c.Locals("token_exp", time.Unix(int64(exp), 0))
		}

		return c.Next()
	}
}
