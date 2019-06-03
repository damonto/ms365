package v1

import (
	"fmt"
	"net/http"
	"net/url"
	"office365/config"

	"github.com/gin-gonic/gin"
)

// Authorize 检查跳转到授权页面
func Authorize(c *gin.Context) {
	q := url.Values{}
	q.Add("client_id", config.AppConfig.ClientID)
	q.Add("scope", config.AppConfig.Scope)
	q.Add("redirect_uri", fmt.Sprintf("%v/oauth/callback", config.AppConfig.Domain))
	q.Add("grant_type", "query")
	q.Add("state", config.AppConfig.ClientID)
	query := q.Encode()

	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?%s", config.AppConfig.AuthorizeURL, query))
}
