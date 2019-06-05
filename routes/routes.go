package routes

import (
	"office365/config"
	"office365/controllers/v1/authorize"
	"office365/controllers/v1/skus"

	"github.com/gin-gonic/gin"
)

// InitRoutes init routes
func InitRoutes() *gin.Engine {
	r := gin.New()

	// global middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// authorize
	r.GET("/", authorize.Authorize)
	r.GET("/oauth/callback", authorize.Callback)

	apiv1 := r.Group("/api/v1", gin.BasicAuth(gin.Accounts{
		config.APIConfig.AccessKey: config.APIConfig.AccessSecret,
	}))
	{
		apiv1.GET("/", skus.SubscribedSkus)
	}

	return r
}
