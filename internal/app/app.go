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
		authorizeCtrl := new(controller.AuthorizeController)
		r.GET("/", authorizeCtrl.Redirect)
		r.GET("/oauth/callback", authorizeCtrl.Callback)
	}

	return r
}
