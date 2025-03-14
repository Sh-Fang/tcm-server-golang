package controllers

import (
	"encoding/json"
	"net/http"
	"tcm-server-go/services" // 引入你定义的服务

	"github.com/gin-gonic/gin"
)

// 调用 C++ /match 接口
func CallMatch(c *gin.Context) {
	// 从请求中获取数据（假设是 JSON 格式）
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 将 map 数据序列化为 JSON 格式的字节数组
	dataBytes, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}

	// 调用 C++ API
	response, err := services.CallMatchAPI(dataBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 解析 C++ API 返回的 JSON 字符串
	var responseData map[string]interface{}
	err = json.Unmarshal([]byte(response), &responseData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// 返回 C++ API 的响应
	c.JSON(http.StatusOK, responseData)
}

// 调用 C++ /progress 接口
func CallProgress(c *gin.Context) {
	// 调用 C++ API
	response, err := services.CallProgressAPI()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 解析 C++ API 返回的 JSON 字符串
	var responseData map[string]interface{}
	err = json.Unmarshal([]byte(response), &responseData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// 返回 C++ API 的响应
	c.JSON(http.StatusOK, responseData)
}
