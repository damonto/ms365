package controller

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/damonto/office365/internal/pkg/config"
	"github.com/damonto/office365/internal/pkg/microsoft"
	"github.com/gin-gonic/gin"
)

// AuthorizeController struct
type AuthorizeController struct{}

// microsoft graph api instance.
var graphAPI = microsoft.NewGraphAPI()

// Redirect to microsoft graph api authorize page.
func (ctl AuthorizeController) Redirect(c *gin.Context) {
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

// Callback recieve authorize result and use it get an access token
func (ctl AuthorizeController) Callback(c *gin.Context) {
	if c.Query("error") != "" {
		c.JSON(rootCtl.wrap(http.StatusBadRequest, c.Query("error_description")))
		return
	}

	err := graphAPI.GetAccessToken(c.Query("code"))
	if err != nil {
		c.JSON(rootCtl.wrap(http.StatusInternalServerError, err))
		return
	}

	c.JSON(rootCtl.wrap(http.StatusOK))
}
