package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// 调用 C++ API 的 /match 接口（POST 请求）
func CallMatchAPI(data []byte) (string, error) {
	// C++ /match 接口 URL
	url := "http://localhost:8081/match"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("error calling C++ /match API: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response from C++ /match API: %v", err)
	}

	// 返回响应结果
	return string(body), nil
}

// 调用 C++ API 的 /progress 接口（GET 请求）
func CallProgressAPI() (string, error) {
	// C++ /progress 接口 URL
	url := "http://localhost:8081/progress"
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error calling C++ /progress API: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response from C++ /progress API: %v", err)
	}

	// 返回响应结果
	return string(body), nil
}
