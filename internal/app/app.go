package app

import (
	"net/http"

	"github.com/damonto/office365/internal/app/controllers"
	"github.com/gin-gonic/gin"
)

// Handler returns the Gin engine
func Handler() http.Handler {
	r := gin.Default()
	{
		authorizeCtrl := new(controllers.AuthorizeController)
		r.GET("/", authorizeCtrl.Redirect)
	}

	return r
}
