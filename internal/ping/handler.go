package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandlePing(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{
		"message": "Pong",
	})
}
