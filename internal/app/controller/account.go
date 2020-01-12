package controller

import (
	"github.com/damonto/msonline/internal/pkg/microsoft"
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
	if err := microsoft.NewStore().Delete(c.Param("id")); err != nil {
		c.JSON(rootCtl.wrap(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(rootCtl.wrap(http.StatusOK))
}
