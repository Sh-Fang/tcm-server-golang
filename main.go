package main

import (
	"tcm-server-go/config"
	"tcm-server-go/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 初始化数据库连接
	config.InitDB()

	// 路由
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// 添加 C++ API 接口路由
	r.POST("/match", controllers.CallMatch)
	r.GET("/progress", controllers.CallProgress)

	// 启动服务
	r.Run(":8080")
}
