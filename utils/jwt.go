package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret 用于签名和验证 JWT，确保与生成 token 时使用的 Secret Key 相同
var jwtSecret = []byte("your_secret_key") // 替换为你的密钥

// 创建 JWT
func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // 72 小时有效期
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseJWT 解析并验证 JWT，返回用户名（或其他需要的字段）
func ParseJWT(tokenString string) (string, error) {
	// 解析 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// 返回验证 token 的密钥
		return jwtSecret, nil
	})

	// 错误处理：如果 token 无效或解析失败
	if err != nil {
		return "", fmt.Errorf("could not parse token: %w", err)
	}

	// 确保 token 是有效的
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 获取用户名
		username, ok := claims["username"].(string)
		if !ok {
			return "", fmt.Errorf("could not extract username from token")
		}
		// 返回用户名
		return username, nil
	}

	// 如果 token 无效
	return "", fmt.Errorf("invalid token")
}
