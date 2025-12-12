package middleware

import (
	"movie-app/internal/config"
	"movie-app/internal/constants"
	"movie-app/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required", nil)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization format", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", nil)
			c.Abort()
			return
		}

		// Simpan claims ke context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set(constants.CtxUserID, claims["user_id"])
			c.Set(constants.CtxName, claims["name"])
			c.Set(constants.CtxRole, claims["role"])
		}

		c.Next()
	}
}
