package user

import (
	"net/http"
	"office365/model"
	"office365/request"
	"strings"

	"github.com/gin-gonic/gin"
)

// NewUser 创建用户
type NewUser struct {
	AccountID     string `form:"account_id" json:"account_id" binding:"required"`
	Enabled       bool   `form:"enabled" json:"enabled" binding:"required"`
	Nickname      string `form:"nickname" json:"nickname" binding:"required"`
	Email         string `form:"email" json:"email" binding:"required"`
	Password      string `form:"password" json:"password" binding:"required"`
	AssignLicense bool   `form:"assign_license" json:"assign_license" binding:"required"`
	SkuID         string `form:"sku_id" json:"sku_id" binding:"required"`
}

// CreateUser 创建新的 Office 365 用户
func CreateUser(c *gin.Context) {
	var json NewUser
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "create user failed",
			"description": err.Error(),
		})
	}

	var account = model.Account{}
	model.DB.Where("user_id = ?", json.AccountID).Find(&account)
	s := strings.Split(account.Email, "@")

	user, err := request.CreateUser(json.AccountID, json.Enabled, json.Nickname, json.Email, json.Password, s[1])
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":       "create user failed",
			"description": err.Error(),
		})
		return
	}

	if json.AssignLicense {
		err = request.AssignLicense(json.AccountID, json.SkuID, user["id"].(string))
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":       "assign license failed",
				"description": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"error":       "",
		"description": "",
	})
}
