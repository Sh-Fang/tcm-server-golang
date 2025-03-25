package routes

import (
	"tcm-server-go/controllers"
	"tcm-server-go/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置所有路由
func SetupRoutes(r *gin.Engine) {
	// 启用 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"}, // 允许前端访问的地址
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true, // 允许跨域请求携带 Cookie
	}))

	// auth 路由
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/login", controllers.Login)
		authRoutes.POST("/updateUserProfile", controllers.UpdateUserProfile)
		authRoutes.POST("/updateUserPassword", controllers.UpdateUserPassword)
	}

	// C++ API 接口路由
	matchRoutes := r.Group("/match")
	matchRoutes.Use(middleware.AuthMiddleware())
	{
		matchRoutes.POST("/subgraphMatch", controllers.SubgraphMatching)
		matchRoutes.GET("/getSubgraphMatchProgress", controllers.GetSubgraphMatchingProgress)
	}

	// 分析图数据
	analyzeRoutes := r.Group("/analyze")
	analyzeRoutes.Use(middleware.AuthMiddleware())
	{
		analyzeRoutes.POST("/streamGraph", controllers.AnalyzeStreamGraph)
		analyzeRoutes.POST("/queryGraph", controllers.AnalyzeQueryGraph)
	}

	// 图可视化
	visualizationRoutes := r.Group("/visualize")
	visualizationRoutes.Use(middleware.AuthMiddleware())
	{
		visualizationRoutes.POST("/streamGraph", controllers.VisualizeStreamGraph)
		visualizationRoutes.POST("/queryGraph", controllers.VisualizeQueryGraph)
	}
}