package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ljk20041215/nutrition-tracker/internal/service"
)

// FoodRecordHandler 食物记录处理器
type FoodRecordHandler struct {
	foodService service.FoodRecordService
}

// NewFoodRecordHandler 创建食物记录处理器实例
func NewFoodRecordHandler(foodService service.FoodRecordService) *FoodRecordHandler {
	return &FoodRecordHandler{foodService: foodService}
}

// CreateFoodRecord 创建食物记录
// @Summary 创建食物记录
// @Description 创建新的食物记录
// @Tags 食物记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateFoodRecordRequest true "创建食物记录请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/food-records [post]
func (h *FoodRecordHandler) CreateFoodRecord(c *gin.Context) {
	// 从认证中间件设置的上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 绑定请求参数
	var req service.CreateFoodRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 创建食物记录
	foodRecord, err := h.foodService.CreateFoodRecord(c.Request.Context(), userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建食物记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    foodRecord,
	})
}

// GetFoodRecord 获取食物记录
// @Summary 获取食物记录
// @Description 根据ID获取食物记录
// @Tags 食物记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "食物记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/food-records/{id} [get]
func (h *FoodRecordHandler) GetFoodRecord(c *gin.Context) {
	// 从认证中间件设置的上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 获取路径参数
	foodID := c.Param("id")
	if foodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "食物记录ID不能为空"})
		return
	}

	foodRecord, err := h.foodService.GetFoodRecord(c.Request.Context(), userID.(string), foodID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "获取食物记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    foodRecord,
	})
}

// GetFoodRecordsByMeal 获取餐次下的食物记录
// @Summary 获取餐次下的食物记录
// @Description 根据餐次ID获取所有食物记录
// @Tags 食物记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param meal_id query string true "餐次记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/food-records/meal [get]
func (h *FoodRecordHandler) GetFoodRecordsByMeal(c *gin.Context) {
	// 从认证中间件设置的上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 获取查询参数
	mealID := c.Query("meal_id")
	if mealID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "餐次记录ID不能为空"})
		return
	}

	foodRecords, err := h.foodService.GetFoodRecordsByMeal(c.Request.Context(), userID.(string), mealID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取食物记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    foodRecords,
	})
}

// GetFoodRecordsByDate 获取指定日期的食物记录
// @Summary 获取指定日期的食物记录
// @Description 获取用户在指定日期的所有食物记录
// @Tags 食物记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date query string true "日期，格式：YYYY-MM-DD"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/food-records [get]
func (h *FoodRecordHandler) GetFoodRecordsByDate(c *gin.Context) {
	// 从认证中间件设置的上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 获取查询参数
	dateStr := c.Query("date")
	if dateStr == "" {
		// 如果没有提供日期，默认使用今天
		dateStr = time.Now().Format("2006-01-02")
	}

	// 解析日期
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式错误，应为 YYYY-MM-DD"})
		return
	}

	foodRecords, err := h.foodService.GetFoodRecordsByDate(c.Request.Context(), userID.(string), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取食物记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    foodRecords,
	})
}

// UpdateFoodRecord 更新食物记录
// @Summary 更新食物记录
// @Description 更新食物记录信息
// @Tags 食物记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "食物记录ID"
// @Param request body service.UpdateFoodRecordRequest true "更新食物记录信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/food-records/{id} [put]
func (h *FoodRecordHandler) UpdateFoodRecord(c *gin.Context) {
	// 从认证中间件设置的上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 获取路径参数
	foodID := c.Param("id")
	if foodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "食物记录ID不能为空"})
		return
	}

	var req service.UpdateFoodRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	foodRecord, err := h.foodService.UpdateFoodRecord(c.Request.Context(), userID.(string), foodID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新食物记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    foodRecord,
	})
}

// DeleteFoodRecord 删除食物记录
// @Summary 删除食物记录
// @Description 根据ID删除食物记录
// @Tags 食物记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "食物记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/food-records/{id} [delete]
func (h *FoodRecordHandler) DeleteFoodRecord(c *gin.Context) {
	// 从认证中间件设置的上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 获取路径参数
	foodID := c.Param("id")
	if foodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "食物记录ID不能为空"})
		return
	}

	if err := h.foodService.DeleteFoodRecord(c.Request.Context(), userID.(string), foodID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除食物记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}


