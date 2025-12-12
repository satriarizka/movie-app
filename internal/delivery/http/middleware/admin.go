package middleware

import (
	"movie-app/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			utils.ErrorResponse(c, http.StatusForbidden, "Access denied: admins only", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
