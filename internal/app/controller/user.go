package controller

import (
	"github.com/damonto/msonline-webapi/internal/pkg/microsoft"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserController struct
type UserController struct{}

// Users retrieve a list of users
func (ctl UserController) Users(c *gin.Context) {
	users, err := microsoft.NewUser().Users(c.Param("id"), c.Query("next"))
	if err != nil {
		c.JSON(rootCtl.wrap(http.StatusInternalServerError, err))
	}

	c.JSON(rootCtl.wrap(http.StatusOK, users))
}

// Delete user
func (ctl UserController) Delete(c *gin.Context) {
	err := microsoft.NewUser().Delete(c.Query("id"), c.Query("uid"))
	if err != nil {
		c.JSON(rootCtl.wrap(http.StatusInternalServerError, err))
	}

	c.JSON(rootCtl.wrap(http.StatusNoContent))
}
