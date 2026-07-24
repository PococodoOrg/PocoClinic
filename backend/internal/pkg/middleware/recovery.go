package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	pkgerrors "github.com/PococodoOrg/PocoClinic/internal/pkg/errors"
	"github.com/PococodoOrg/PocoClinic/internal/pkg/logging"
	"github.com/gin-gonic/gin"
)

// Recovery returns middleware that catches panics, logs them, and keeps the server running.
func Recovery(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				if logger != nil {
					logger.Error("panic recovered",
						fmt.Errorf("panic: %v", recovered),
						"method", c.Request.Method,
						"path", logging.SanitizeString(c.Request.URL.Path),
					)
					logger.Logger.Error("panic stack trace", "stack", string(debug.Stack()))
				}

				if c.Writer.Written() {
					return
				}

				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    pkgerrors.ErrInternalServer,
					"error":   "An internal error occurred",
					"message": "An internal error occurred",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
