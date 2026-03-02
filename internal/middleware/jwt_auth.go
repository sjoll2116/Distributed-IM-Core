package middleware

import (
	"kama_chat_server/pkg/util/jwt"
	"kama_chat_server/pkg/zlog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于 JWT 识别身份验证信息的 Gin 中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 客户端需要按格式携带：Authorization: Bearer <token>
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			zlog.Error("缺少授权头 (Authorization Header)")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未登录或登录已过期",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			zlog.Error("授权头格式错误: " + authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "非法的 Authorization 格式",
			})
			c.Abort()
			return
		}

		// 解析 Token
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			zlog.Error("解析 Token 失败: " + err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "登录已失效，请重新登录",
			})
			c.Abort()
			return
		}

		// 验证通过，将 uuid 保存到上下文中
		c.Set("uuid", claims.Uuid)
		c.Next()
	}
}
