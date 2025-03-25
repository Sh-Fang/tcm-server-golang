package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 定义 JWT 密钥（生产环境中应使用环境变量存储）
var jwtSecret = []byte("your-secret-key") // 请替换为安全的密钥

// 生成 Token
func GenerateToken(userID uint) (string, error) {
	// 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24小时过期
	})

	// 生成 Token 字符串
	return token.SignedString(jwtSecret)
}

// 验证 Token
func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保使用的是正确的签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
}


// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization 头
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			c.Abort()
			return
		}

		// 提取 Token
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// 验证 Token
		token, err := ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 Token"})
			c.Abort()
			return
		}

		// 解析 Token Claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 解析失败"})
			c.Abort()
			return
		}

		// 将 user_id 存入上下文
		c.Set("user_id", uint(claims["user_id"].(float64)))

		// 继续请求
		c.Next()
	}
}