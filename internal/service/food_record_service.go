package service

import (
	"context"
	"errors"
	"time"

	"github.com/ljk20041215/nutrition-tracker/internal/model"
	"github.com/ljk20041215/nutrition-tracker/internal/repository"
)

// FoodRecordService 食物记录服务接口
type FoodRecordService interface {
	CreateFoodRecord(ctx context.Context, userID string, req *CreateFoodRecordRequest) (*model.FoodRecord, error)
	GetFoodRecord(ctx context.Context, userID string, foodID string) (*model.FoodRecord, error)
	GetFoodRecordsByMeal(ctx context.Context, userID string, mealID string) ([]*model.FoodRecord, error)
	GetFoodRecordsByDate(ctx context.Context, userID string, date time.Time) ([]*model.FoodRecord, error)
	UpdateFoodRecord(ctx context.Context, userID string, foodID string, req *UpdateFoodRecordRequest) (*model.FoodRecord, error)
	DeleteFoodRecord(ctx context.Context, userID string, foodID string) error
}

// foodRecordService 食物记录服务实现
type foodRecordService struct {
	foodRecordRepo repository.FoodRecordRepository
	mealRepo       repository.MealRecordRepository
	userRepo       repository.UserRepository
	foodRepo       repository.FoodRepository
}

// NewFoodRecordService 创建食物记录服务实例
func NewFoodRecordService(
	foodRecordRepo repository.FoodRecordRepository,
	mealRepo repository.MealRecordRepository,
	userRepo repository.UserRepository,
	foodRepo repository.FoodRepository,
) FoodRecordService {
	return &foodRecordService{
		foodRecordRepo: foodRecordRepo,
		mealRepo:       mealRepo,
		userRepo:       userRepo,
		foodRepo:       foodRepo,
	}
}

// CreateFoodRecordRequest 创建食物记录请求
type CreateFoodRecordRequest struct {
	MealRecordID string  `json:"meal_record_id" binding:"required"` // 餐次记录ID
	FoodID       string  `json:"food_id" binding:"required"`        // 食物ID
	Quantity     float64 `json:"quantity" binding:"required,gt=0"`   // 份量
	Unit         string  `json:"unit" binding:"required"`           // 单位（g, kg, ml, 个等）
}

// UpdateFoodRecordRequest 更新食物记录请求
type UpdateFoodRecordRequest struct {
	Quantity float64 `json:"quantity" binding:"required,gt=0"` // 份量
}

// CreateFoodRecord 创建食物记录
func (s *foodRecordService) CreateFoodRecord(ctx context.Context, userID string, req *CreateFoodRecordRequest) (*model.FoodRecord, error) {
	// 检查用户是否存在
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查餐次记录是否存在且属于该用户
	mealRecord, err := s.mealRepo.FindByID(ctx, req.MealRecordID)
	if err != nil {
		return nil, errors.New("餐次记录不存在")
	}

	if mealRecord.UserID != userID {
		return nil, errors.New("无权限访问该餐次记录")
	}

	// 获取食物信息
	food, err := s.foodRepo.FindByID(ctx, req.FoodID)
	if err != nil {
		return nil, errors.New("食物不存在")
	}

	// 计算实际摄入的营养成分（假设基础数据是每100g的含量）
	calories := (req.Quantity / 100) * food.Calories
	protein := (req.Quantity / 100) * food.Protein
	carbohydrates := (req.Quantity / 100) * food.Carbohydrates
	fat := (req.Quantity / 100) * food.Fat

	// 创建食物记录
	foodRecord := &model.FoodRecord{
		MealRecordID:  req.MealRecordID,
		FoodID:        req.FoodID,
		FoodName:      food.Name,
		Quantity:      req.Quantity,
		Unit:          req.Unit,
		Calories:      calories,
		Protein:       protein,
		Carbohydrates: carbohydrates,
		Fat:           fat,
	}

	if err := s.foodRecordRepo.Create(ctx, foodRecord); err != nil {
		return nil, errors.New("创建食物记录失败")
	}

	return foodRecord, nil
}

// UpdateFoodRecord 更新食物记录
func (s *foodRecordService) UpdateFoodRecord(ctx context.Context, userID string, foodID string, req *UpdateFoodRecordRequest) (*model.FoodRecord, error) {
	// 获取食物记录
	foodRecord, err := s.foodRecordRepo.FindByID(ctx, foodID)
	if err != nil {
		return nil, errors.New("食物记录不存在")
	}

	// 检查餐次记录是否属于该用户
	mealRecord, err := s.mealRepo.FindByID(ctx, foodRecord.MealRecordID)
	if err != nil {
		return nil, errors.New("餐次记录不存在")
	}

	if mealRecord.UserID != userID {
		return nil, errors.New("无权限访问该食物记录")
	}

	// 获取食物信息
	food, err := s.foodRepo.FindByID(ctx, foodRecord.FoodID)
	if err != nil {
		return nil, errors.New("食物不存在")
	}

	// 更新份量
	foodRecord.Quantity = req.Quantity

	// 重新计算营养成分（基于食物的基础营养数据和新的份量）
	foodRecord.Calories = (req.Quantity / 100) * food.Calories
	foodRecord.Protein = (req.Quantity / 100) * food.Protein
	foodRecord.Carbohydrates = (req.Quantity / 100) * food.Carbohydrates
	foodRecord.Fat = (req.Quantity / 100) * food.Fat

	// 更新记录
	if err := s.foodRecordRepo.Update(ctx, foodRecord); err != nil {
		return nil, errors.New("更新食物记录失败")
	}

	return foodRecord, nil
}

// GetFoodRecord 获取食物记录
func (s *foodRecordService) GetFoodRecord(ctx context.Context, userID string, foodID string) (*model.FoodRecord, error) {
	// 获取食物记录
	foodRecord, err := s.foodRecordRepo.FindByID(ctx, foodID)
	if err != nil {
		return nil, errors.New("食物记录不存在")
	}

	// 检查餐次记录是否属于该用户
	mealRecord, err := s.mealRepo.FindByID(ctx, foodRecord.MealRecordID)
	if err != nil {
		return nil, errors.New("餐次记录不存在")
	}

	if mealRecord.UserID != userID {
		return nil, errors.New("无权限访问该食物记录")
	}

	return foodRecord, nil
}

// GetFoodRecordsByMeal 获取餐次下的所有食物记录
func (s *foodRecordService) GetFoodRecordsByMeal(ctx context.Context, userID string, mealID string) ([]*model.FoodRecord, error) {
	// 检查餐次记录是否存在且属于该用户
	mealRecord, err := s.mealRepo.FindByID(ctx, mealID)
	if err != nil {
		return nil, errors.New("餐次记录不存在")
	}

	if mealRecord.UserID != userID {
		return nil, errors.New("无权限访问该餐次记录")
	}

	// 获取食物记录
	foodRecords, err := s.foodRecordRepo.FindByMealRecordID(ctx, mealID)
	if err != nil {
		return nil, errors.New("获取食物记录失败")
	}

	return foodRecords, nil
}

// GetFoodRecordsByDate 获取指定日期的所有食物记录
func (s *foodRecordService) GetFoodRecordsByDate(ctx context.Context, userID string, date time.Time) ([]*model.FoodRecord, error) {
	// 检查用户是否存在
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 获取食物记录
	foodRecords, err := s.foodRecordRepo.FindByUserIDAndDate(ctx, userID, date)
	if err != nil {
		return nil, errors.New("获取食物记录失败")
	}

	return foodRecords, nil
}

// DeleteFoodRecord 删除食物记录
func (s *foodRecordService) DeleteFoodRecord(ctx context.Context, userID string, foodID string) error {
	// 获取食物记录
	foodRecord, err := s.foodRecordRepo.FindByID(ctx, foodID)
	if err != nil {
		return errors.New("食物记录不存在")
	}

	// 检查餐次记录是否属于该用户
	mealRecord, err := s.mealRepo.FindByID(ctx, foodRecord.MealRecordID)
	if err != nil {
		return errors.New("餐次记录不存在")
	}

	if mealRecord.UserID != userID {
		return errors.New("无权限删除该食物记录")
	}

	// 删除食物记录
	if err := s.foodRecordRepo.Delete(ctx, foodID); err != nil {
		return errors.New("删除食物记录失败")
	}

	return nil
}

