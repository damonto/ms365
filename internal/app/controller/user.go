package controller

import (
	"github.com/damonto/msonline/internal/pkg/microsoft"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserController struct
type UserController struct{}

// Users retrieve a list of users
func (ctl UserController) Users(c *gin.Context) {
	users, err := microsoft.NewUser().Users(c.Param("id"), c.Query("next"))
	if err != nil {
		c.JSON(rootCtl.wrap(http.StatusInternalServerError, err.Error()))
	}

	c.JSON(rootCtl.wrap(http.StatusOK, users))
}

// Delete user
func (ctl UserController) Delete(c *gin.Context) {
	if err := microsoft.NewUser().Delete(c.Param("id"), c.Param("uid")); err != nil {
		c.JSON(rootCtl.wrap(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(rootCtl.wrap(http.StatusOK))
}

// Create a new user
func (ctl UserController) Create(c *gin.Context) {
	var createRequest microsoft.CreateUserRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		c.JSON(rootCtl.wrap(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	if err := microsoft.NewUser().Create(c.Param("id"), createRequest); err != nil {
		c.JSON(rootCtl.wrap(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(rootCtl.wrap(http.StatusOK))
}
