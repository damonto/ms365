package authorize

import (
	"fmt"
	"net/http"
	"net/url"
	"office365/config"
	"office365/request"

	"github.com/gin-gonic/gin"
)

// Authorize 检查跳转到授权页面
func Authorize(c *gin.Context) {
	q := url.Values{
		"client_id":     {config.AppConfig.ClientID},
		"scope":         {config.AppConfig.Scope},
		"redirect_uri":  {fmt.Sprintf("%v/oauth/callback", config.AppConfig.Domain)},
		"grant_type":    {"query"},
		"response_type": {"code"},
		"state":         {"nouse"},
	}

	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?%s", config.AppConfig.AuthorizeURL, q.Encode()))
}

// Callback OAuth
func Callback(c *gin.Context) {
	if c.Query("error") != "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       c.Query("error"),
			"description": c.Query("error_description"),
		})
	} else {
		err := request.GetAccessToken(c.Query("code"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":       "internal error",
				"description": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"error":       "",
				"description": "Authorized",
			})
		}
	}
}
