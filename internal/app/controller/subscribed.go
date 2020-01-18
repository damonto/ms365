package controller

import (
	"net/http"

	"github.com/damonto/msonline/internal/pkg/microsoft"
	"github.com/gin-gonic/gin"
)

// SubscribedController struct
type SubscribedController struct{}

//Skus Get the list of commercial subscriptions that an organization has acquired.
func (ctl SubscribedController) Skus(c *gin.Context) {
	skus, err := microsoft.NewSubscribed().ListSubscribedSkus(c.Param("id"))
	if err != nil {
		c.JSON(rootCtl.wrap(http.StatusBadRequest, err.Error()))
		return
	}

	c.JSON(rootCtl.wrap(http.StatusOK, skus))
}
