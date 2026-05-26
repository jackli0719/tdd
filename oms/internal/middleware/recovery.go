package middleware

import (
	"log"
	"net/http"
	"oms/pkg/errors"
	"oms/pkg/response"

	"github.com/gin-gonic/gin"
)

// Recovery returns a middleware that recovers from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic with stack trace
				log.Printf("[PANIC] Recovered from panic: %v", r)

				// Return 500 Internal Server Error
				err := errors.Internal("internal server error")
				response.ErrorWithStatus(c, http.StatusInternalServerError, err.Code, err.Message)

				// Abort the request
				c.Abort()
			}
		}()

		c.Next()
	}
}
