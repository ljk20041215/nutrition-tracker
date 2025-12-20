package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ljk20041215/nutrition-tracker/internal/auth"
	"github.com/ljk20041215/nutrition-tracker/internal/handler"
	"github.com/ljk20041215/nutrition-tracker/internal/model"
	"github.com/ljk20041215/nutrition-tracker/internal/repository"
	"github.com/ljk20041215/nutrition-tracker/internal/service"
	"github.com/ljk20041215/nutrition-tracker/pkg/database"
)

func main() {
	// 1. 数据库配置（确保与调试脚本一致）
	host := "localhost"
	port := "5432"
	username := "postgres"
	password := "ljk071311" // ⚠️ 确保这里和调试脚本用相同的密码
	dbname := "nutrition_tracker"

	log.Println("🔍 主程序启动 - 开始初始化")

	// 2. 初始化数据库
	log.Printf("🔌 连接数据库: %s@%s:%s/%s", username, host, port, dbname)

	err := database.Init(host, port, username, password, dbname)
	if err != nil {
		log.Fatalf("❌ 数据库初始化失败: %v", err)
	}
	log.Println("✅ 数据库初始化成功")

	// 3. 获取数据库实例并验证
	db := database.GetDB()
	if db == nil {
		log.Fatal("❌ 致命错误: database.GetDB() 返回 nil")
	}
	log.Println("✅ 获取到有效的数据库实例")

	// 4. 测试数据库查询（确认连接可用）
	var userCount int64
	if err := db.Raw("SELECT COUNT(*) FROM users").Scan(&userCount).Error; err != nil {
		log.Printf("⚠️ 数据库查询测试失败（可能表不存在）: %v", err)
	} else {
		log.Printf("📊 当前用户数: %d", userCount)
	}

	// 5. 初始化 Repository
	log.Println("🔄 初始化 UserRepository...")
	userRepo := repository.NewUserRepository(db)
	if userRepo == nil {
		log.Fatal("❌ UserRepository 初始化失败")
	}
	log.Println("✅ UserRepository 初始化成功")

	// 初始化 NutritionGoalRepository
	log.Println("🔄 初始化 NutritionGoalRepository...")
	goalRepo := repository.NewNutritionGoalRepository(db)
	if goalRepo == nil {
		log.Fatal("❌ NutritionGoalRepository 初始化失败")
	}
	log.Println("✅ NutritionGoalRepository 初始化成功")

	// 6. 初始化 Service
	log.Println("🔄 初始化 UserService...")
	userService := service.NewUserService(userRepo)
	if userService == nil {
		log.Fatal("❌ UserService 初始化失败")
	}
	log.Println("✅ UserService 初始化成功")

	// 初始化 NutritionGoalService
	log.Println("🔄 初始化 NutritionGoalService...")
	goalService := service.NewNutritionGoalService(goalRepo, userRepo)
	if goalService == nil {
		log.Fatal("❌ NutritionGoalService 初始化失败")
	}
	log.Println("✅ NutritionGoalService 初始化成功")

	// 7. 初始化 Handler
	log.Println("🔄 初始化 AuthHandler...")
	authHandler := handler.NewAuthHandler(userService)
	if authHandler == nil {
		log.Fatal("❌ AuthHandler 初始化失败")
	}
	log.Println("✅ AuthHandler 初始化成功")

	// 8. 初始化 UserHandler
	log.Println("🔄 初始化 UserHandler...")
	userHandler := handler.NewUserHandler(userService)
	if userHandler == nil {
		log.Fatal("❌ UserHandler 初始化失败")
	}
	log.Println("✅ UserHandler 初始化成功")

	// 初始化 NutritionGoalHandler
	log.Println("🔄 初始化 NutritionGoalHandler...")
	goalHandler := handler.NewNutritionGoalHandler(goalService)
	if goalHandler == nil {
		log.Fatal("❌ NutritionGoalHandler 初始化失败")
	}
	log.Println("✅ NutritionGoalHandler 初始化成功")

	// 9. 创建Gin引擎
	log.Println("🔄 创建Gin引擎...")
	r := gin.Default()

	// 10. 添加日志中间件
	r.Use(func(c *gin.Context) {
		log.Printf("🌐 %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	// 11. 注册路由
	log.Println("🔄 注册路由...")

	// 公开路由（无需认证）
	public := r.Group("/api/v1")
	{
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
		public.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"service": "Nutrition Tracker",
				"version": "1.0.0",
			})
		})

		// 添加一个简单的数据库测试端点
		public.GET("/test/db", func(c *gin.Context) {
			var count int64
			if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "数据库查询失败: " + err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message":    "数据库连接正常",
				"user_count": count,
			})
		})
	}

	// 受保护路由（需要认证）
	protected := r.Group("/api/v1")
	protected.Use(auth.AuthMiddleware())
	{
		// 用户相关路由
		protected.GET("/users/profile", userHandler.GetProfile)
		protected.PUT("/users/profile", userHandler.UpdateProfile)

		// 营养目标相关路由
		protected.GET("/goals", goalHandler.GetNutritionGoal)
		protected.POST("/goals", goalHandler.SetNutritionGoal)
		protected.POST("/goals/calculate", goalHandler.CalculateNutritionGoal)
	}

	// 12. 启动服务器
	log.Println("🚀 服务器启动完成，开始监听 :8080")
	log.Println("📝 注册接口: POST http://localhost:8080/api/v1/auth/register")
	log.Println("🔑 登录接口: POST http://localhost:8080/api/v1/auth/login")
	log.Println("🧪 数据库测试: GET http://localhost:8080/api/v1/test/db")
	log.Println("🎯 营养目标接口: GET/POST http://localhost:8080/api/v1/goals")
	log.Println("⚡ 计算营养目标接口: POST http://localhost:8080/api/v1/goals/calculate")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}