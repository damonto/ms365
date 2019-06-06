package routes

import (
	"office365/config"
	"office365/controllers/v1/account"
	"office365/controllers/v1/authorize"
	"office365/controllers/v1/skus"
	"office365/controllers/v1/user"

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
		apiv1.GET("/accounts", account.ListAccount)
		apiv1.GET("/skus/:userID", skus.SubscribedSkus)
		apiv1.POST("/users", user.CreateUser)
	}

	return r
}
