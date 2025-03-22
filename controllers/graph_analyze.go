package controllers

import (
	"bufio"
	"strings"

	"github.com/gin-gonic/gin"
)



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

	// 初始化
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