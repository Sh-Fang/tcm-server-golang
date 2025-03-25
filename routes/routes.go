package routes

import (
	"tcm-server-go/controllers"

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
	apiRoutes := r.Group("/match")
	{
		apiRoutes.POST("/subgraphMatch", controllers.SubgraphMatching)
		apiRoutes.GET("/getSubgraphMatchProgress", controllers.GetSubgraphMatchingProgress)
	}

	// 分析图数据
	graphAnalysisRoutes := r.Group("/analyze")
	{
		graphAnalysisRoutes.POST("/streamGraph", controllers.AnalyzeStreamGraph)
		graphAnalysisRoutes.POST("/queryGraph", controllers.AnalyzeQueryGraph)
	}

	// 图可视化
	visualizationRoutes := r.Group("/visualize")
	{
		visualizationRoutes.POST("/streamGraph", controllers.VisualizeStreamGraph)
		visualizationRoutes.POST("/queryGraph", controllers.VisualizeQueryGraph)
	}
}