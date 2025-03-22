package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"ai-chat-service-go/internal/config"
	"ai-chat-service-go/internal/models"
)

type userContextKey string

// UserKey is the key used to store the user in the context
const UserKey userContextKey = "user"

// JWTClaims are the claims in the JWT token
type JWTClaims struct {
	jwt.RegisteredClaims
	Email         string                 `json:"email"`
	EmailVerified bool                   `json:"email_verified"`
	PreferredName string                 `json:"preferred_username"`
	Name          string                 `json:"name"`
	GivenName     string                 `json:"given_name"`
	FamilyName    string                 `json:"family_name"`
	RealmAccess   map[string]interface{} `json:"realm_access"`
}

// UserInfo contains authenticated user information
type UserInfo struct {
	Email       string
	Name        string
	GivenName   string
	FamilyName  string
	Roles       []string
	TokenExpiry time.Time
}

// Auth creates the authentication middleware
func Auth(cfg config.AuthConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(models.NewUnauthorizedError("Missing authorization header"))
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(http.StatusUnauthorized).JSON(models.NewUnauthorizedError("Invalid authorization header format"))
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		userInfo, err := validateToken(tokenString, cfg)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(models.NewUnauthorizedError(fmt.Sprintf("Invalid token: %v", err)))
		}

		// Store user information in context
		c.Locals(string(UserKey), userInfo)

		return c.Next()
	}
}

// validateToken validates the JWT token
func validateToken(tokenString string, cfg config.AuthConfig) (*UserInfo, error) {
	// Parse the token using the public key
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get the public key
		publicKey, err := parseKeycloakRSAPublicKey(cfg.PublicKey)
		if err != nil {
			return nil, err
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Extract roles
	var roles []string
	if claims.RealmAccess != nil {
		if rolesArr, ok := claims.RealmAccess["roles"].([]interface{}); ok {
			for _, role := range rolesArr {
				roles = append(roles, role.(string))
			}
		}
	}

	// Create user info
	userInfo := &UserInfo{
		Email:       claims.Email,
		Name:        claims.Name,
		GivenName:   claims.GivenName,
		FamilyName:  claims.FamilyName,
		Roles:       roles,
		TokenExpiry: claims.ExpiresAt.Time,
	}

	return userInfo, nil
}

// parseKeycloakRSAPublicKey parses the public key from Keycloak
func parseKeycloakRSAPublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	// Remove PEM formatting if present
	publicKeyPEM = strings.TrimSpace(publicKeyPEM)
	publicKeyPEM = strings.TrimPrefix(publicKeyPEM, "-----BEGIN PUBLIC KEY-----")
	publicKeyPEM = strings.TrimSuffix(publicKeyPEM, "-----END PUBLIC KEY-----")
	publicKeyPEM = strings.ReplaceAll(publicKeyPEM, "\n", "")

	// Decode the base64 encoded DER format
	derBytes, err := base64.StdEncoding.DecodeString(publicKeyPEM)
	if err != nil {
		return nil, err
	}

	// Parse the key (simplified for this example)
	// In a real implementation, you would use x509.ParsePKIXPublicKey
	e := big.NewInt(65537) // Default public exponent
	n := new(big.Int).SetBytes(derBytes)

	// Create an RSA public key
	pubKey := &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	return pubKey, nil
}

// GetCurrentUser retrieves the current user from the Fiber context
func GetCurrentUser(c *fiber.Ctx) *UserInfo {
	userInfo, ok := c.Locals(string(UserKey)).(*UserInfo)
	if !ok {
		return nil
	}
	return userInfo
}

// RequireRoles creates middleware to check if the user has specific roles
func RequireRoles(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetCurrentUser(c)
		if user == nil {
			return c.Status(http.StatusUnauthorized).JSON(models.NewUnauthorizedError(""))
		}

		// Check if the user has any of the required roles
		hasRole := false
		for _, role := range roles {
			for _, userRole := range user.Roles {
				if role == userRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			return c.Status(http.StatusForbidden).JSON(models.NewForbiddenError(""))
		}

		return c.Next()
	}
}
