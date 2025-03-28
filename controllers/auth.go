package controllers

import (
	"net/http"
	"tcm-server-go/config"
	"tcm-server-go/models"
	"tcm-server-go/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Test",
	})
}

func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not encrypt password"})
		return
	}

	// 创建用户
	user := models.User{
		UserID: utils.GenerateSonyflakeID(),
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Phone:    "",
		Bio:      "",
	}

	// 写入数据库
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully!",
		"user": gin.H{
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func Login(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 查找用户
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid email"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful!",
		"user": gin.H{
			"name":  user.Name,
			"email": user.Email,
			"token": user.UserID,
		},
	})
}

func UpdateUserProfile(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 获取请求体中的 email
	userEmail := input.Email

	// 查找用户
	var userToUpdate models.User
	if err := config.DB.Where("email = ?", userEmail).First(&userToUpdate).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 更新用户信息
	if err := config.DB.Model(&userToUpdate).Updates(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "User updated successfully!",
		"data": gin.H{
			"name":  userToUpdate.Name,
			"email": userToUpdate.Email,
			"phone": userToUpdate.Phone,
			"bio":   userToUpdate.Bio,
		},
	})
}

func UpdateUserPassword(c *gin.Context) {
	var passwordData struct {
		Email           string `json:"email"`
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	// 绑定 JSON 数据
	if err := c.ShouldBindJSON(&passwordData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 查找用户
	var user models.User
	if err := config.DB.Where("email = ?", passwordData.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 验证旧密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordData.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect current password"})
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	// 更新数据库中的密码
	if err := config.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// 返回成功信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Password updated successfully!",
	})

}
