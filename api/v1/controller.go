package v1

import (
	"kama_chat_server/pkg/util/response"

	"github.com/gin-gonic/gin"
)

func JsonBack(c *gin.Context, message string, ret int, data interface{}) {
	switch ret {
	case 0:
		if data != nil {
			response.SuccessWithMessage(c, message, data)
		} else {
			response.SuccessWithMessage(c, message, nil)
		}
	case -2:
		response.Fail(c, 400, message)
	case -1:
		response.Fail(c, 500, message)
	}
}
