package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/damonto/ms365/internal/app/controller"
	"github.com/damonto/ms365/internal/app/middleware"
	"github.com/damonto/ms365/internal/pkg/config"
)

// Handler returns the Gin engine
func Handler() http.Handler {
	r := gin.Default()
	r.Use(middleware.CORS())
	{
		authorizeCtl := new(controller.AuthorizeController)
		r.GET("/oauth/authorize", authorizeCtl.Redirect)
		r.GET("/oauth/callback", authorizeCtl.Callback)
	}

	api := r.Group("/api/v1", gin.BasicAuth(gin.Accounts{
		config.Cfg.App.AccessKey: config.Cfg.App.AccessSecret,
	}))
	{
		{
			ctl := new(controller.AccountController)
			api.GET("/accounts", ctl.Accounts)
			api.DELETE("/accounts/:id", ctl.Delete)
		}
		{
			ctl := new(controller.SubscribedController)
			api.GET("/accounts/:id/skus", ctl.Skus)
		}
		{
			ctl := new(controller.UserController)
			api.GET("/accounts/:id/users", ctl.Users)
			api.POST("/accounts/:id/users", ctl.Create)
			api.GET("/accounts/:id/users/:uid", ctl.User)
			api.PATCH("/accounts/:id/users/:uid", ctl.Update)
			api.DELETE("/accounts/:id/users/:uid", ctl.Delete)
		}
	}

	return r
}
