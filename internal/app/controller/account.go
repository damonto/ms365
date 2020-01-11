package controller

import (
	"github.com/damonto/msonline-webapi/internal/pkg/microsoft"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AccountController authorized microsoft account controller
type AccountController struct{}

// Accounts returns all authorized microsoft accounts
func (ctl AccountController) Accounts(c *gin.Context) {
	accounts, err := microsoft.NewStore().All()
	if err != nil {
		c.JSON(rootCtl.wrap(http.StatusInternalServerError, err.Error()))
		return
	}

	var resp = make([]map[string]string, 0)
	for _, v := range accounts {
		resp = append(resp, map[string]string{
			"id":    v.ID,
			"email": v.Email,
		})
	}

	c.JSON(rootCtl.wrap(http.StatusOK, resp))
}

// Delete an item from leveldb
func (ctl AccountController) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(rootCtl.wrap(http.StatusUnprocessableEntity, "id can not be null"))
		return
	}

	err := microsoft.NewStore().Delete(id)
	if err != nil {
		c.JSON(rootCtl.wrap(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(rootCtl.wrap(http.StatusOK))
}

// Skus Get the list of commercial subscriptions that an organization has acquired.
func (ctl AccountController) Skus(c *gin.Context) {
	skus, err := microsoft.NewSubscribed().ListSubscribedSkus(c.Param("id"))
	if err != nil {
		c.JSON(rootCtl.wrap(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(rootCtl.wrap(http.StatusOK, skus))
}
