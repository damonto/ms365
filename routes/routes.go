package routes

import (
	v1 "office365/controllers/v1"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// InitRoutes init routes
func InitRoutes() *gin.Engine {
	r := gin.New()

	// global middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// sessions
	r.Use(sessions.Sessions("sess", cookie.NewStore([]byte("secret"))))

	// authorize
	r.GET("/", v1.Authorize)
	r.GET("/oauth/callback", v1.Callback)

	// Dashboard
	dashboard := r.Group("/dashboard")
	{
		dashboard.GET("/", v1.Dashboard)
	}

	return r
}
