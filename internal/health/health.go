// package health contains the HTTP handlers.
package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping is a simple liveness check.
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
