package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	// 引入你定义的服务
	"github.com/gin-gonic/gin"
)

// 调用 C++ /match 接口
func CallMatch(c *gin.Context) {
	// 获取上传的文件
	streamGraph, err := c.FormFile("streamGraph") // "streamGraph" 是前端上传的字段名
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	queryGraph, err := c.FormFile("queryGraph")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var params map[string]any
	paramsStr := c.PostForm("params")
	if err := json.Unmarshal([]byte(paramsStr), &params); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON parameters"})
		return
	}

	// 把所有非字符串的参数转换成字符串类型，如果本身字符串是"true"或"false"，则转换成"y"或"n"
	for key, value := range params {
		if _, ok := value.(string); !ok {
			if value == true {
				params[key] = "y"
			} else if value == false {
				params[key] = "n"
			} else {
				params[key] = fmt.Sprintf("%v", value)
			}
		}
	}

	// 目标路径
	streamFilePath := "/tmp/tcm_download/" + streamGraph.Filename
	queryFilePath := "/tmp/tcm_download/" + queryGraph.Filename

	if err := c.SaveUploadedFile(streamGraph, streamFilePath); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := c.SaveUploadedFile(queryGraph, queryFilePath); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	params["stream_path"] = streamFilePath
	params["query_path"] = queryFilePath

	// 将params封装为 JSON
	jsonParams, _ := json.Marshal(params)

	// 发送 POST 请求到 C++ 后端
	response, err := http.Post("http://localhost:8081/match", "application/json", strings.NewReader(string(jsonParams)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer response.Body.Close()

	// 解析 C++ API 返回的 JSON 字符串
	var responseData map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// 删除临时文件
	os.Remove(streamFilePath)
	os.Remove(queryFilePath)

	// 返回 C++ API 的响应
	c.JSON(http.StatusOK, responseData)

}

// 调用 C++ /progress 接口
func CallProgress(c *gin.Context) {
	url := "http://localhost:8081/progress"
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	fmt.Println(string(body))

	// 返回 C++ API 的响应
	c.JSON(http.StatusOK, gin.H{"progress": string(body)})
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
