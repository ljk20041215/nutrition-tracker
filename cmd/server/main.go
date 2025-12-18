// cmd/server/main.go
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化 Gin 路由引擎
	r := gin.Default()

	// 注册一个最基础的健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Nutrition Tracker is running!",
		})
	})

	// 用户相关路由组 (后续扩展)
	userRoutes := r.Group("/api/users")
	{
		userRoutes.POST("/register", func(c *gin.Context) {
			// TODO: 实现注册逻辑
			c.JSON(http.StatusOK, gin.H{"message": "Register endpoint (TODO)"})
		})
		userRoutes.POST("/login", func(c *gin.Context) {
			// TODO: 实现登录逻辑
			c.JSON(http.StatusOK, gin.H{"message": "Login endpoint (TODO)"})
		})
	}

	// 食物记录路由组 (后续扩展)
	r.POST("/api/foods", func(c *gin.Context) {
		// TODO: 使用 channel 异步记录食物 (体现 Go 特色)
		c.JSON(http.StatusOK, gin.H{"message": "Food recorded (TODO)"})
	})

	// 启动服务器，监听 8080 端口
	log.Println("🚀 Server starting on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
