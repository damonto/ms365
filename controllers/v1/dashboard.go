package v1

import (
	"fmt"
	"net/http"
	"office365/request"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Dashboard office 365 management dashboard
func Dashboard(c *gin.Context) {
	session := sessions.Default(c)
	skus, err := request.GetSubscribedSkus(fmt.Sprintf("%v", session.Get("access_token")))
	if err != nil {
		// due to lazy
		if err.Error() == "InvalidAuthenticationToken" {
			session := sessions.Default(c)
			session.Delete("access_token")
			session.Save()

			c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Get subscribed skus failed",
			"description": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"skus": skus,
		})
	}
}
