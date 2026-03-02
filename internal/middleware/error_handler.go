package middleware

import (
	"fmt"
	"kama_chat_server/pkg/util/response"
	"kama_chat_server/pkg/zlog"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// GlobalRecovery 统一异常处理中间件
func GlobalRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取 panic 堆栈信息
				stack := debug.Stack()
				errStr := fmt.Sprintf("Panic recovered: %v\n%s", err, string(stack))

				// 记录日志
				zlog.Error(errStr)

				// 返回统一规范的 500 错误 JSON
				response.FailWithMessage(c, "服务器内部发生异常，请稍后再试")

				// 中断后续的处理函数
				c.Abort()
			}
		}()
		c.Next()
	}
}
