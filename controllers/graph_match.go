package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// 调用 C++ /match 接口
func SubgraphMatching(c *gin.Context) {
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

	// 解析 POST 请求的 JSON 参数
	var params map[string]any
	paramsStr := c.PostForm("params")
	if err := json.Unmarshal([]byte(paramsStr), &params); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON parameters"})
		return
	}

	// 把所有非字符串的参数转换成字符串类型
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

	// 保存上传的文件
	if err := c.SaveUploadedFile(streamGraph, streamFilePath); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := c.SaveUploadedFile(queryGraph, queryFilePath); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 拼接参数
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
func GetSubgraphMatchingProgress(c *gin.Context) {
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

	// 返回 C++ API 的响应
	c.JSON(http.StatusOK, gin.H{"progress": string(body)})
}


