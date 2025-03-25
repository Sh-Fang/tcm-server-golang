package utils

import (
	"fmt"
	"time"

	"github.com/sony/sonyflake"
)

var flake *sonyflake.Sonyflake

// 包被导入时自动执行初始化 Sonyflake
func init() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: time.Now(), // 可调整为业务上线时间，避免 ID 过大
	})
	if flake == nil {
		panic("Sonyflake 初始化失败")
	}
}

// 生成唯一 ID
func GenerateSonyflakeID() string {
	id, err := flake.NextID()
	if err != nil {
		fmt.Println("Sonyflake 生成 ID 失败:", err)
		return ""
	}
	return fmt.Sprintf("%d", id) // 转成字符串
}