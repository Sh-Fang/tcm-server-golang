package controllers

import (
	"encoding/json"
	"net/http"
	"tcm-server-go/services" // 引入你定义的服务
	"time"

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

func AnalyzeStreamGraph(c *gin.Context) {
	// 获取上传的文件
	streamGraph, err := c.FormFile("streamGraph") // "streamGraph" 是前端上传的字段名
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 你可以读取文件内容
	file, err := streamGraph.Open()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// 如果你需要读取文件内容，可以使用 io.Reader
	fileContent := make([]byte, streamGraph.Size)
	_, err = file.Read(fileContent)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	time.Sleep(5 * time.Second)

	// 返回响应
	c.JSON(200, gin.H{
		"fileSize":            82345678,
		"nodeCount":           5000,
		"edgeCount":           10000,
		"maxDegree":           100,
		"avgDegree":           4.0,
		"density":             0.05,
		"connectedComponents": 3,
		"filePath":            "/path/to/graph/stream.txt",
	})
}

func AnalyzeQueryGraph(c *gin.Context) {
	// 获取上传的文件
	queryGraph, err := c.FormFile("queryGraph")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 你可以读取文件内容
	file, err := queryGraph.Open()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// 如果你需要读取文件内容，可以使用 io.Reader
	fileContent := make([]byte, queryGraph.Size)
	_, err = file.Read(fileContent)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	time.Sleep(5 * time.Second)

	// 返回响应
	c.JSON(200, gin.H{
		"fileSize":            12345678,
		"nodeCount":           5000,
		"edgeCount":           10000,
		"maxDegree":           100,
		"avgDegree":           4.0,
		"density":             0.05,
		"connectedComponents": 3,
		"filePath":            "/path/to/graph/query.txt",
	})
}
