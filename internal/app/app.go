package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler returns the Gin engine
func Handler() http.Handler {
	r := gin.Default()

	return r
}
