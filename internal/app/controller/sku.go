package controller

import (
	"net/http"

	"github.com/damonto/msonline/internal/pkg/microsoft"
	"github.com/gin-gonic/gin"
)

// SkuController struct
type SkuController struct{}

//Skus Get the list of commercial subscriptions that an organization has acquired.
func (ctl SkuController) Skus(c *gin.Context) {
	skus, err := microsoft.NewSubscribed().ListSubscribedSkus(c.Param("id"))
	if err != nil {
		c.JSON(rootCtl.wrap(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(rootCtl.wrap(http.StatusOK, skus))
}
