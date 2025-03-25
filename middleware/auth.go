package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Authorization 令牌
		token := c.GetHeader("Authorization")

		// 这里假设 Token 是 "Bearer <token_value>" 形式
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			c.Abort() // 终止请求
			return
		}

		// 提取 Token 值（去掉 "Bearer " 前缀）
		token = strings.TrimPrefix(token, "Bearer ")

		// 这里应该校验 Token 是否有效（如调用 JWT 库解析 Token）
		if !validateToken(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 Token"})
			c.Abort()
			return
		}

		// 认证通过，继续执行请求
		c.Next()
	}
}

// 这里是一个假设的 Token 验证函数，你需要替换成真正的 JWT 验证逻辑
func validateToken(token string) bool {
	// 在这里解析 JWT Token，检查有效性
	// 这里只是模拟，实际应用中应该使用 JWT 解析库，如 github.com/golang-jwt/jwt/v4
	return token == "valid-token" // 假设 "valid-token" 是一个合法的 Token
}