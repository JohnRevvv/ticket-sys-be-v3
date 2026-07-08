package jwt

import (
	"errors"
	"ideyanale-be/pkg/config"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

// Permissions mirrors the boolean flags on model.Roles.
type Permissions struct {
	CanCreateTicket  bool
	CanEndorseTicket bool
	CanApproveTicket bool
	CanResolveTicket bool
	CanAudit         bool
}

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

// NEW: permission-based check, for ticketing actions
func RequirePermission(c fiber.Ctx, check func(Permissions) bool) error {
	perms, ok := c.Locals("permissions").(Permissions)
	if !ok {
		return errors.New("unauthorized")
	}

	if !check(perms) {
		return errors.New("forbidden")
	}

	return nil
}

func GenerateSuperAdminToken(ID int, username string) (string, error) {
	claims := jwt.MapClaims{
		"id":       ID,
		"username": username,
		"role":     "Super-Admin",
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

// UPDATED: institutionID is now uint (matches model.UserDetails), and
// permission flags are embedded alongside the role name.
func GenerateUserToken(ID int, staffID string, institutionID uint, roleID uint, roleName string, perms Permissions) (string, error) {

	claims := jwt.MapClaims{
		"id":              ID,
		"staff_id":        staffID,
		"institution_id":  institutionID,
		"role_id":         roleID,
		"role":            roleName,
		"can_create":      perms.CanCreateTicket,
		"can_endorse":     perms.CanEndorseTicket,
		"can_approve":     perms.CanApproveTicket,
		"can_resolve":     perms.CanResolveTicket,
		"can_audit":       perms.CanAudit,
		"exp":             time.Now().Add(1 * time.Hour).Unix(),
		"iat":             time.Now().Unix(),
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

		// role_id
		if roleID, ok := claims["role_id"].(float64); ok {
			c.Locals("role_id", uint(roleID))
		}

		// institution_id (safe optional field)
		if inst, ok := claims["institution_id"]; ok && inst != nil {
			if instFloat, ok := inst.(float64); ok {
				c.Locals("institution_id", uint(instFloat))
			}
		}

		// NEW: rebuild Permissions struct from individual claims
		getBool := func(key string) bool {
			v, ok := claims[key].(bool)
			return ok && v
		}
		c.Locals("permissions", Permissions{
			CanCreateTicket:  getBool("can_create"),
			CanEndorseTicket: getBool("can_endorse"),
			CanApproveTicket: getBool("can_approve"),
			CanResolveTicket: getBool("can_resolve"),
			CanAudit:         getBool("can_audit"),
		})

		return c.Next()
	}
}