package main

import (
	"tcm-server-go/config"
	"tcm-server-go/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 初始化数据库连接
	config.InitDB()

	// 启用 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"}, // 允许前端访问的地址
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true, // 允许跨域请求携带 Cookie
	}))

	// auth 路由
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/updateUserProfile", controllers.UpdateUserProfile)
	r.POST("/updateUserPassword", controllers.UpdateUserPassword)

	// C++ API 接口路由
	r.POST("/subgraphMatching", controllers.SubgraphMatching)
	r.GET("/getSubgraphMatchingProgress", controllers.GetSubgraphMatchingProgress)

	// 分析图数据
	r.POST("/analyzeStreamGraph", controllers.AnalyzeStreamGraph)
	r.POST("/analyzeQueryGraph", controllers.AnalyzeQueryGraph)

	// 图可视化
	r.POST("/visualizeStreamGraph", controllers.VisualizeStreamGraph)
	r.POST("/visualizeQueryGraph", controllers.VisualizeQueryGraph)

	// 启动服务
	r.Run(":8082")
}
