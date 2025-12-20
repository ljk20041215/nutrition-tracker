package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ljk20041215/nutrition-tracker/internal/service"
)

// MealRecordHandler 餐次记录处理器
type MealRecordHandler struct {
	mealService service.MealRecordService
}

// NewMealRecordHandler 创建餐次记录处理器实例
func NewMealRecordHandler(mealService service.MealRecordService) *MealRecordHandler {
	return &MealRecordHandler{mealService: mealService}
}

// CreateMealRecord 创建餐次记录
// @Summary 创建餐次记录
// @Description 创建新的餐次记录
// @Tags 餐次记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateMealRecordRequest true "创建餐次记录信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/meals [post]
func (h *MealRecordHandler) CreateMealRecord(c *gin.Context) {
	// 从认证中间件设置的上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var req service.CreateMealRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	mealRecord, err := h.mealService.CreateMealRecord(c.Request.Context(), userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建餐次记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    mealRecord,
	})
}

// GetMealRecord 获取餐次记录
// @Summary 获取餐次记录
// @Description 根据ID获取餐次记录
// @Tags 餐次记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "餐次记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/meals/{id} [get]
func (h *MealRecordHandler) GetMealRecord(c *gin.Context) {
	// 从认证中间件设置的上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 获取路径参数
	mealID := c.Param("id")
	if mealID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "餐次记录ID不能为空"})
		return
	}

	mealRecord, err := h.mealService.GetMealRecord(c.Request.Context(), userID.(string), mealID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "获取餐次记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    mealRecord,
	})
}

// GetMealRecordsByDate 获取指定日期的餐次记录
// @Summary 获取指定日期的餐次记录
// @Description 获取用户在指定日期的所有餐次记录
// @Tags 餐次记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date query string true "日期，格式：YYYY-MM-DD"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/meals [get]
func (h *MealRecordHandler) GetMealRecordsByDate(c *gin.Context) {
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

	mealRecords, err := h.mealService.GetMealRecordsByDate(c.Request.Context(), userID.(string), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取餐次记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    mealRecords,
	})
}

// DeleteMealRecord 删除餐次记录
// @Summary 删除餐次记录
// @Description 根据ID删除餐次记录
// @Tags 餐次记录
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "餐次记录ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/meals/{id} [delete]
func (h *MealRecordHandler) DeleteMealRecord(c *gin.Context) {
	// 从认证中间件设置的上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 获取路径参数
	mealID := c.Param("id")
	if mealID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "餐次记录ID不能为空"})
		return
	}

	if err := h.mealService.DeleteMealRecord(c.Request.Context(), userID.(string), mealID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除餐次记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
