package jwt

import (
	"errors"
	"ideyanale-be/pkg/config"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

// type BaseClaims struct {
// 	ID   int    `json:"id"`
// 	Role string `json:"role"`
// 	jwt.RegisteredClaims
// }

func RequireRoles(c fiber.Ctx, allowed ...string) error {
	role, ok := c.Locals("role").(string)
	if !ok {
		return errors.New("unauthorized")
	}

	for _, r := range allowed {
		if r == role {
			return nil
		}
	}

	return errors.New("forbidden")
}

func GenerateSuperAdminToken(ID int, username string) (string, error) {

	claims := jwt.MapClaims{
		"id":       ID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

func GenerateUserToken(ID int, staffID string, institutionID int, role string) (string, error) {

	claims := jwt.MapClaims{
		"id":       ID,
		"staff_id": staffID,
		"institution_id": institutionID,
		"role":         role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

func JWTProtected() fiber.Handler {
	return func(c fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing token")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid signing method")
			}
			return []byte(config.JWTSecret), nil
		})

		if err != nil || token == nil || !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token claims")
		}

		// id
		if id, ok := claims["id"].(float64); ok {
			c.Locals("id", int(id))
		}

		// staff_id
		if staffID, ok := claims["staff_id"].(string); ok {
			c.Locals("staff_id", staffID)
		}

		// role
		if role, ok := claims["role"].(string); ok {
			c.Locals("role", role)
		}

		// institution_id (safe optional field)
		if inst, ok := claims["institution_id"]; ok && inst != nil {
			if instFloat, ok := inst.(float64); ok {
				c.Locals("institution_id", int(instFloat))
			}
		}

		return c.Next()
	}
}
