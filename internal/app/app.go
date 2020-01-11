package app

import (
	"net/http"

	"github.com/damonto/office365/internal/app/controller"
	"github.com/gin-gonic/gin"
)

// Handler returns the Gin engine
func Handler() http.Handler {
	r := gin.Default()
	{
		authorizeCtl := new(controller.AuthorizeController)
		r.GET("/oauth/authorize", authorizeCtl.Redirect)
		r.GET("/oauth/callback", authorizeCtl.Callback)
	}

	return r
}
