package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

type rootController struct{}

var rootCtl rootController

func (r rootController) wrap(code int, obj ...interface{}) (int, response) {
	var msg string
	var data interface{}
	n := len(obj)
	if n > 0 {
		if n == 1 {
			switch obj[0].(type) {
			case string:
				msg = obj[0].(string)
			default:
				data = obj[0]
			}
		} else {
			data = obj[0]
			msg = obj[1].(string)
		}
	}

	resp := response{
		Message:    msg,
		StatusCode: code,
		Data:       data,
	}

	return code, resp
}

func checkErr(c *gin.Context, err error) error {
	if err != nil {
		c.JSON(rootCtl.wrap(http.StatusInternalServerError, err.Error()))
	}
	return err
}
