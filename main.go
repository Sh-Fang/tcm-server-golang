package main

import (
	"tcm-server-go/config"
	"tcm-server-go/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 初始化数据库连接
	config.InitDB()

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务
	r.Run(":8082")
}
