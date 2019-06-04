package v1

import (
	"fmt"
	"office365/request"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Dashboard office 365 management dashboard
func Dashboard(c *gin.Context) {
	session := sessions.Default(c)
	request.GetSubscribedSkus(fmt.Sprintf("%v", session.Get("access_token")))
}
