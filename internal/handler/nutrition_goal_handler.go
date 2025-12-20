package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ljk20041215/nutrition-tracker/internal/service"
)

type NutritionGoalHandler struct {
	goalService service.NutritionGoalService
}

func NewNutritionGoalHandler(goalService service.NutritionGoalService) *NutritionGoalHandler {
	return &NutritionGoalHandler{goalService: goalService}
}

// GetNutritionGoal 获取营养目标
// @Summary 获取营养目标
// @Description 获取当前登录用户的营养目标
// @Tags 营养目标
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/goals [get]
func (h *NutritionGoalHandler) GetNutritionGoal(c *gin.Context) {
	// 从认证中间件设置的上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	goal, err := h.goalService.GetNutritionGoal(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "获取营养目标失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    goal,
	})
}

// SetNutritionGoal 设置营养目标
// @Summary 设置营养目标
// @Description 手动设置当前登录用户的营养目标
// @Tags 营养目标
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.SetGoalRequest true "营养目标信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/goals [post]
func (h *NutritionGoalHandler) SetNutritionGoal(c *gin.Context) {
	// 从认证中间件设置的上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var req service.SetGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	goal, err := h.goalService.SetNutritionGoal(c.Request.Context(), userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "设置营养目标失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "设置成功",
		"data":    goal,
	})
}

// CalculateNutritionGoal 自动计算营养目标
// @Summary 自动计算营养目标
// @Description 根据用户信息自动计算营养目标
// @Tags 营养目标
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CalculateGoalRequest true "计算参数"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/goals/calculate [post]
func (h *NutritionGoalHandler) CalculateNutritionGoal(c *gin.Context) {
	// 从认证中间件设置的上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var req service.CalculateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	goal, err := h.goalService.CalculateNutritionGoal(c.Request.Context(), userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "计算营养目标失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "计算成功",
		"data":    goal,
	})
}
