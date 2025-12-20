package service

import (
	"context"
	"errors"
	"time"

	"github.com/ljk20041215/nutrition-tracker/internal/model"
	"github.com/ljk20041215/nutrition-tracker/internal/repository"
)

// MealRecordService 餐次记录服务接口
type MealRecordService interface {
	CreateMealRecord(ctx context.Context, userID string, req *CreateMealRecordRequest) (*model.MealRecord, error)
	GetMealRecord(ctx context.Context, userID string, mealID string) (*model.MealRecord, error)
	GetMealRecordsByDate(ctx context.Context, userID string, date time.Time) ([]*model.MealRecord, error)
	DeleteMealRecord(ctx context.Context, userID string, mealID string) error
}

// mealRecordService 餐次记录服务实现
type mealRecordService struct {
	mealRepo repository.MealRecordRepository
	userRepo repository.UserRepository
}

// NewMealRecordService 创建餐次记录服务实例
func NewMealRecordService(
	mealRepo repository.MealRecordRepository,
	userRepo repository.UserRepository,
) MealRecordService {
	return &mealRecordService{
		mealRepo: mealRepo,
		userRepo: userRepo,
	}
}

// CreateMealRecordRequest 创建餐次记录请求
type CreateMealRecordRequest struct {
	Date     string       `json:"record_date" binding:"required,datetime=2006-01-02"` // 日期格式：YYYY-MM-DD
	MealType model.MealType `json:"meal_type" binding:"required"` // 餐次类型：breakfast/lunch/dinner/snack 或 1/2/3/4
}

// CreateMealRecord 创建餐次记录
func (s *mealRecordService) CreateMealRecord(ctx context.Context, userID string, req *CreateMealRecordRequest) (*model.MealRecord, error) {
	// 检查用户是否存在
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 解析日期
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("日期格式错误，应为 YYYY-MM-DD")
	}

	// 检查该用户在该日期该餐次是否已存在
	existing, _ := s.mealRepo.FindByUserIDDateAndType(ctx, userID, date, req.MealType)
	if existing != nil {
		return nil, errors.New("该餐次记录已存在")
	}

	// 创建餐次记录
	mealRecord := &model.MealRecord{
		UserID:   userID,
		Date:     date,
		MealType: req.MealType,
	}

	if err := s.mealRepo.Create(ctx, mealRecord); err != nil {
		return nil, errors.New("创建餐次记录失败")
	}

	return mealRecord, nil
}

// GetMealRecord 获取餐次记录
func (s *mealRecordService) GetMealRecord(ctx context.Context, userID string, mealID string) (*model.MealRecord, error) {
	// 获取餐次记录
	mealRecord, err := s.mealRepo.FindByID(ctx, mealID)
	if err != nil {
		return nil, errors.New("餐次记录不存在")
	}

	// 检查权限
	if mealRecord.UserID != userID {
		return nil, errors.New("无权限访问该餐次记录")
	}

	return mealRecord, nil
}

// GetMealRecordsByDate 获取指定日期的所有餐次记录
func (s *mealRecordService) GetMealRecordsByDate(ctx context.Context, userID string, date time.Time) ([]*model.MealRecord, error) {
	// 检查用户是否存在
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 获取餐次记录
	mealRecords, err := s.mealRepo.FindByUserIDAndDate(ctx, userID, date)
	if err != nil {
		return nil, errors.New("获取餐次记录失败")
	}

	return mealRecords, nil
}

// DeleteMealRecord 删除餐次记录
func (s *mealRecordService) DeleteMealRecord(ctx context.Context, userID string, mealID string) error {
	// 获取餐次记录
	mealRecord, err := s.mealRepo.FindByID(ctx, mealID)
	if err != nil {
		return errors.New("餐次记录不存在")
	}

	// 检查权限
	if mealRecord.UserID != userID {
		return errors.New("无权限删除该餐次记录")
	}

	// 删除餐次记录
	if err := s.mealRepo.Delete(ctx, mealID); err != nil {
		return errors.New("删除餐次记录失败")
	}

	return nil
}

