package app

import (
	"net/http"

	"github.com/damonto/msonline-webapi/internal/app/controller"
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

	api := r.Group("/api/v1")
	{
		{
			accountCtl := new(controller.AccountController)
			api.GET("/accounts", accountCtl.Accounts)
			api.DELETE("/accounts/:id", accountCtl.Delete)
		}
	}

	return r
}
