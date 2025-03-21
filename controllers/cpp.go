package controllers

import (
	"bufio"
	"encoding/json"
	"net/http"
	"strings"
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

func AnalyzeStreamGraph(c *gin.Context) {
	// 获取上传的文件
	streamGraph, err := c.FormFile("streamGraph") // "streamGraph" 是前端上传的字段名
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 读取文件内容
	file, err := streamGraph.Open()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	nodeDegrees := make(map[string]int)
	edgeCount := 0

	// 流式读取文件内容
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) != 6 {
			continue // 跳过格式不对的行
		}

		src, dst := parts[0], parts[1]
		nodeDegrees[src]++
		nodeDegrees[dst]++
		edgeCount++
	}

	if err := scanner.Err(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 计算最大度数 & 平均度数
	maxDegree, totalDegree := 0, 0
	for _, degree := range nodeDegrees {
		totalDegree += degree
		if degree > maxDegree {
			maxDegree = degree
		}
	}

	nodeCount := len(nodeDegrees)
	avgDegree := float64(totalDegree) / float64(nodeCount)
	density := float64(edgeCount) / float64(nodeCount*(nodeCount-1))

	// 获取文件大小
	fileSize := streamGraph.Size

	// 文件路径
	filePath := streamGraph.Filename

	// 返回响应
	c.JSON(200, gin.H{
		"fileSize":            fileSize,
		"nodeCount":           nodeCount,
		"edgeCount":           edgeCount,
		"maxDegree":           maxDegree,
		"avgDegree":           avgDegree,
		"density":             density,
		"connectedComponents": 3,
		"filePath":            filePath,
	})
}

func AnalyzeQueryGraph(c *gin.Context) {
	// 获取上传的文件
	queryGraph, err := c.FormFile("queryGraph")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 读取文件内容
	file, err := queryGraph.Open()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	nodeDegrees := make(map[string]int)
	edgeCount := 0

	// 流式读取文件内容
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 || parts[0] != "e" {
			continue // 跳过格式不对的行
		}

		src, dst := parts[2], parts[3]
		nodeDegrees[src]++
		nodeDegrees[dst]++
		edgeCount++
	}

	if err := scanner.Err(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 计算最大度数 & 平均度数
	maxDegree, totalDegree := 0, 0
	for _, degree := range nodeDegrees {
		totalDegree += degree
		if degree > maxDegree {
			maxDegree = degree
		}
	}

	nodeCount := len(nodeDegrees)
	avgDegree := float64(totalDegree) / float64(nodeCount)
	density := float64(edgeCount) / float64(nodeCount*(nodeCount-1))

	// 获取文件大小
	fileSize := queryGraph.Size

	// 文件路径
	filePath := queryGraph.Filename

	// 返回响应
	c.JSON(200, gin.H{
		"fileSize":            fileSize,
		"nodeCount":           nodeCount,
		"edgeCount":           edgeCount,
		"maxDegree":           maxDegree,
		"avgDegree":           avgDegree,
		"density":             density,
		"connectedComponents": 3,
		"filePath":            filePath,
	})
}
