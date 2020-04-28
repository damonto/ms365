package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/damonto/ms365/internal/pkg/microsoft"
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

// Create a new user
func (ctl UserController) Create(c *gin.Context) {
	var createRequest microsoft.CreateUserRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		c.JSON(rootCtl.wrap(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	uid, err := microsoft.NewUser().Create(c.Param("id"), createRequest)
	if err != nil {
		c.JSON(rootCtl.wrap(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(rootCtl.wrap(http.StatusCreated, gin.H{
		"id": uid,
	}))
}

// User retrieve the properties and relationships of user object.
func (ctl UserController) User(c *gin.Context) {
	user, err := microsoft.NewUser().Retrieve(c.Param("id"), c.Param("uid"))
	if err != nil {
		c.JSON(rootCtl.wrap(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(rootCtl.wrap(http.StatusOK, user))
}

// Update user object, now only supported update password.
func (ctl UserController) Update(c *gin.Context) {
	var body struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(rootCtl.wrap(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	if err := microsoft.NewUser().UpdatePassword(c.Param("id"), c.Param("uid"), body.Password); err != nil {
		c.JSON(rootCtl.wrap(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(rootCtl.wrap(http.StatusOK))
}

// Delete user
func (ctl UserController) Delete(c *gin.Context) {
	if err := microsoft.NewUser().Delete(c.Param("id"), c.Param("uid")); err != nil {
		c.JSON(rootCtl.wrap(http.StatusInternalServerError, err.Error()))
		return
	}

	// 204 or 200?
	c.JSON(rootCtl.wrap(http.StatusOK))
}
