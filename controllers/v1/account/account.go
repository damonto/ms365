package account

import (
	"net/http"
	"office365/model"

	"github.com/gin-gonic/gin"
)

// ListAccount 列出所有可用的账号
func ListAccount(c *gin.Context) {
	var accounts = model.Account{}
	model.DB.Find(&accounts)

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  accounts,
	})
}
               