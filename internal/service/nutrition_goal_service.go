package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ljk20041215/nutrition-tracker/internal/model"
	"github.com/ljk20041215/nutrition-tracker/internal/repository"
)

type NutritionGoalService interface {
	GetNutritionGoal(ctx context.Context, userID string) (*model.NutritionGoal, error)
	SetNutritionGoal(ctx context.Context, userID string, req *SetGoalRequest) (*model.NutritionGoal, error)
	CalculateNutritionGoal(ctx context.Context, userID string, req *CalculateGoalRequest) (*model.NutritionGoal, error)
}

type nutritionGoalService struct {
	goalRepo repository.NutritionGoalRepository
	userRepo repository.UserRepository
}

func NewNutritionGoalService(
	goalRepo repository.NutritionGoalRepository,
	userRepo repository.UserRepository,
) NutritionGoalService {
	return &nutritionGoalService{
		goalRepo: goalRepo,
		userRepo: userRepo,
	}
}

// SetGoalRequest 手动设置营养目标请求
type SetGoalRequest struct {
	Calories      float64 `json:"calories" binding:"required,gt=0"`
	Protein       float64 `json:"protein" binding:"required,gt=0"`
	Carbohydrates float64 `json:"carbohydrates" binding:"required,gt=0"`
	Fat           float64 `json:"fat" binding:"required,gt=0"`
}

// CalculateGoalRequest 自动计算营养目标请求
type CalculateGoalRequest struct {
	// 可选参数，如果提供则覆盖用户当前资料
	Gender        int     `json:"gender"`
	Age           int     `json:"age"`
	Height        float64 `json:"height"`
	Weight        float64 `json:"weight"`
	ActivityLevel int     `json:"activity_level"`
	GoalType      string  `json:"goal_type" binding:"required,oneof=maintain lose gain"` // maintain: 维持, lose: 减脂, gain: 增肌
}

func (s *nutritionGoalService) GetNutritionGoal(ctx context.Context, userID string) (*model.NutritionGoal, error) {
	goal, err := s.goalRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("获取营养目标失败")
	}
	return goal, nil
}

func (s *nutritionGoalService) SetNutritionGoal(ctx context.Context, userID string, req *SetGoalRequest) (*model.NutritionGoal, error) {
	// 检查用户是否存在
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 查找用户是否已有营养目标
	goal, err := s.goalRepo.FindByUserID(ctx, userID)
	if err != nil {
		// 如果不存在，则创建新目标
		goal = &model.NutritionGoal{
			UserID:        userID,
			Calories:      req.Calories,
			Protein:       req.Protein,
			Carbohydrates: req.Carbohydrates,
			Fat:           req.Fat,
		}
		if err := s.goalRepo.Create(ctx, goal); err != nil {
			return nil, errors.New("创建营养目标失败")
		}
	} else {
		// 如果存在，则更新目标
		goal.Calories = req.Calories
		goal.Protein = req.Protein
		goal.Carbohydrates = req.Carbohydrates
		goal.Fat = req.Fat
		if err := s.goalRepo.Update(ctx, goal); err != nil {
			return nil, errors.New("更新营养目标失败")
		}
	}

	return goal, nil
}

func (s *nutritionGoalService) CalculateNutritionGoal(ctx context.Context, userID string, req *CalculateGoalRequest) (*model.NutritionGoal, error) {
	// 获取用户当前资料
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 使用请求参数覆盖用户当前资料（如果提供）
	gender := user.Gender
	age := user.Age
	height := user.Height
	weight := user.Weight
	activityLevel := user.ActivityLevel

	if req.Gender != 0 {
		gender = req.Gender
	}
	if req.Age > 0 {
		age = req.Age
	}
	if req.Height > 0 {
		height = req.Height
	}
	if req.Weight > 0 {
		weight = req.Weight
	}
	if req.ActivityLevel >= 1 && req.ActivityLevel <= 5 {
		activityLevel = req.ActivityLevel
	}

	// 验证必要参数
	if gender == 0 || age <= 0 || height <= 0 || weight <= 0 || activityLevel == 0 {
		return nil, errors.New("缺少必要的用户信息，请先完善个人资料")
	}

	// 计算BMR（基础代谢率）使用Mifflin-St Jeor公式
	var bmr float64
	if gender == 1 { // 男性
		bmr = 10*weight + 6.25*height - 5*float64(age) + 5
	} else { // 女性
		bmr = 10*weight + 6.25*height - 5*float64(age) - 161
	}

	// 根据活动水平计算TDEE（总能量消耗）
	var activityFactor float64
	switch activityLevel {
	case 1: // 久坐不动
		activityFactor = 1.2
	case 2: // 轻度活跃（每周运动1-3次）
		activityFactor = 1.375
	case 3: // 中度活跃（每周运动3-5次）
		activityFactor = 1.55
	case 4: // 高度活跃（每周运动6-7次）
		activityFactor = 1.725
	case 5: // 非常活跃（从事体力劳动或每天高强度训练）
		activityFactor = 1.9
	default:
		activityFactor = 1.2
	}

	tdee := bmr * activityFactor

	// 根据目标类型调整热量
	var calories float64
	switch req.GoalType {
	case "maintain":
		calories = tdee
	case "lose":
		calories = tdee - 500 // 每天减少500卡路里，每周预计减重0.5kg
	case "gain":
		calories = tdee + 500 // 每天增加500卡路里，每周预计增重0.5kg
	default:
		return nil, errors.New("无效的目标类型")
	}

	// 计算宏量营养素目标（基于热量百分比）
	// 默认比例：蛋白质20%，碳水化合物50%，脂肪30%
	protein := calories * 0.2 / 4  // 1克蛋白质提供4卡路里
	carbs := calories * 0.5 / 4    // 1克碳水化合物提供4卡路里
	fat := calories * 0.3 / 9      // 1克脂肪提供9卡路里

	// 根据目标类型调整宏量营养素比例
	switch req.GoalType {
	case "lose":
		// 减脂期间适当增加蛋白质比例
		protein = calories * 0.25 / 4
		carbs = calories * 0.45 / 4
		fat = calories * 0.3 / 9
	case "gain":
		// 增肌期间适当增加碳水化合物比例
		protein = calories * 0.2 / 4
		carbs = calories * 0.55 / 4
		fat = calories * 0.25 / 9
	}

	// 查找用户是否已有营养目标
	goal, err := s.goalRepo.FindByUserID(ctx, userID)
	if err != nil {
		// 如果不存在，则创建新目标
		goal = &model.NutritionGoal{
			UserID:        userID,
			Calories:      calories,
			Protein:       protein,
			Carbohydrates: carbs,
			Fat:           fat,
		}
		if err := s.goalRepo.Create(ctx, goal); err != nil {
			return nil, errors.New("创建营养目标失败")
		}
	} else {
		// 如果存在，则更新目标
		goal.Calories = calories
		goal.Protein = protein
		goal.Carbohydrates = carbs
		goal.Fat = fat
		if err := s.goalRepo.Update(ctx, goal); err != nil {
			return nil, errors.New("更新营养目标失败")
		}
	}

	return goal, nil
}

// 辅助函数：格式化宏量营养素比例
func formatMacroRatio(calories, protein, carbs, fat float64) string {
	proteinCalories := protein * 4
	carbsCalories := carbs * 4
	fatCalories := fat * 9

	proteinRatio := (proteinCalories / calories) * 100
	carbsRatio := (carbsCalories / calories) * 100
	fatRatio := (fatCalories / calories) * 100

	return fmt.Sprintf("蛋白质: %.1f%%, 碳水化合物: %.1f%%, 脂肪: %.1f%%", proteinRatio, carbsRatio, fatRatio)
}
