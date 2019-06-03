package routes

import (
	v1 "office365/controllers/v1"

	"github.com/gin-gonic/gin"
)

// InitRoutes init routes
func InitRoutes() *gin.Engine {
	r := gin.New()

	// global middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", v1.Authorize)

	return r
}
