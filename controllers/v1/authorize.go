package v1

import (
	"fmt"
	"net/http"
	"net/url"
	"office365/config"
	"office365/request"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

// Authorize 检查跳转到授权页面
func Authorize(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("access_token") != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
	}

	q := url.Values{
		"client_id":     {config.AppConfig.ClientID},
		"scope":         {config.AppConfig.Scope},
		"redirect_uri":  {fmt.Sprintf("%v/oauth/callback", config.AppConfig.Domain)},
		"grant_type":    {"query"},
		"response_type": {"code"},
		"state":         {"unused_now"},
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
		accessToken, err := request.GetAccessToken(c.Query("code"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":       "internal error",
				"description": err.Error(),
			})
		} else {
			session := sessions.Default(c)
			session.Set("access_token", accessToken)
			session.Save()

			c.Redirect(http.StatusMovedPermanently, "/dashboard")
		}
	}
}
