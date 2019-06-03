package v1

import (
	"net/http"
	"office365/config"

	"github.com/gin-gonic/gin"
)

// Authorize 检查跳转到授权页面
func Authorize(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":      "Hello",
		"redirect_uri": config.AppConfig.AuthorizeURL,
	})
}
