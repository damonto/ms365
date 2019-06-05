package skus

import (
	"net/http"
	"office365/request"

	"github.com/gin-gonic/gin"
)

// SubscribedSkus office 365 management dashboard
func SubscribedSkus(c *gin.Context) {
	skus, err := request.GetSubscribedSkus(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Get subscribed skus failed",
			"description": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"error": "",
			"data":  skus,
		})
	}
}
