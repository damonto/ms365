package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/damonto/office365/internal/pkg/config"
	"github.com/gin-gonic/gin"
)

// AuthorizeController struct
type AuthorizeController struct{}

// Redirect to microsoft graph api authorize page.
func (ctrl AuthorizeController) Redirect(c *gin.Context) {
	q := url.Values{
		"client_id":     {config.Cfg.Microsoft.ClientID},
		"scope":         {config.Cfg.Microsoft.Scope},
		"redirect_uri":  {fmt.Sprintf("%v/oauth/callback", config.Cfg.Microsoft.Domain)},
		"grant_type":    {"query"},
		"response_type": {"code"},
		"state":         {"state_code_unused"},
	}

	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?%s", config.Cfg.Microsoft.AuthorizeURL, q.Encode()))
}
