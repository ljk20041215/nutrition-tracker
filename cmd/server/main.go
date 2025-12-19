package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ljk20041215/nutrition-tracker/internal/model"
	"github.com/ljk20041215/nutrition-tracker/pkg/database"
)

func main() {
	// 1. 初始化数据库
	// 暂时硬编码数据库连接信息，后续可以从配置文件读取
	host := "localhost"
	port := "5432"
	username := "postgres"
	password := "ljk071311" // 请替换为你的 PostgreSQL 密码
	dbname := "nutrition_tracker"

	if err := database.Init(host, port, username, password, dbname); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}

	// 2. 创建Gin引擎
	r := gin.Default()

	// 3. 注册中间件（可选，这里添加一个简单的日志中间件）
	r.Use(gin.Logger())

	// 4. 注册路由
	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Nutrition Tracker API is running",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// 数据库测试端点：查询用户总数
	r.GET("/test-db", func(c *gin.Context) {
		db := database.GetDB()
		var count int64
		if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to query database: " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":     "Database connection is healthy",
			"total_users": count,
		})
	})

	// 5. 启动服务器
	portStr := ":8080"
	log.Printf("🚀 Server starting on http://localhost%s", portStr)
	if err := r.Run(portStr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
