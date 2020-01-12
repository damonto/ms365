package app

import (
	"net/http"

	"github.com/damonto/msonline/internal/app/controller"
	"github.com/damonto/msonline/internal/pkg/config"
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

	api := r.Group("/api/v1", gin.BasicAuth(gin.Accounts{
		config.Cfg.App.AccessKey: config.Cfg.App.AccessSecret,
	}))
	{
		{
			accountCtl := new(controller.AccountController)
			api.GET("/accounts", accountCtl.Accounts)
			api.DELETE("/accounts/:id", accountCtl.Delete)
		}
		{
			skuCtl := new(controller.SkuController)
			api.GET("/accounts/:id/skus", skuCtl.Skus)
		}
		{
			userCtl := new(controller.UserController)
			api.GET("/accounts/:id/users", userCtl.Users)
			api.POST("/accounts/:id/users", userCtl.Create)
			api.DELETE("/accounts/:id/users/:uid", userCtl.Delete)
		}
	}

	return r
}
