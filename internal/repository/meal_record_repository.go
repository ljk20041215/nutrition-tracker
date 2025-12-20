package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ljk20041215/nutrition-tracker/internal/model"
	"gorm.io/gorm"
)

// MealRecordRepository 餐次记录仓库接口
type MealRecordRepository interface {
	Create(ctx context.Context, mealRecord *model.MealRecord) error
	FindByID(ctx context.Context, id string) (*model.MealRecord, error)
	FindByUserIDAndDate(ctx context.Context, userID string, date time.Time) ([]*model.MealRecord, error)
	FindByUserIDDateAndType(ctx context.Context, userID string, date time.Time, mealType model.MealType) (*model.MealRecord, error)
	Update(ctx context.Context, mealRecord *model.MealRecord) error
	Delete(ctx context.Context, id string) error
}

// mealRecordRepository 餐次记录仓库实现
type mealRecordRepository struct {
	db *gorm.DB
}

// NewMealRecordRepository 创建餐次记录仓库实例
func NewMealRecordRepository(db *gorm.DB) MealRecordRepository {
	if db == nil {
		log.Fatal("❌ NewMealRecordRepository: db 参数为 nil")
	}
	return &mealRecordRepository{db: db}
}

// Create 创建餐次记录
func (r *mealRecordRepository) Create(ctx context.Context, mealRecord *model.MealRecord) error {
	if r == nil || r.db == nil {
		return errors.New("repository 未初始化")
	}
	return r.db.WithContext(ctx).Create(mealRecord).Error
}

// FindByID 根据ID查找餐次记录
func (r *mealRecordRepository) FindByID(ctx context.Context, id string) (*model.MealRecord, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("repository 未初始化")
	}

	var mealRecord model.MealRecord
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&mealRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("餐次记录不存在")
		}
		return nil, err
	}

	return &mealRecord, nil
}

// FindByUserIDAndDate 根据用户ID和日期查找餐次记录
func (r *mealRecordRepository) FindByUserIDAndDate(ctx context.Context, userID string, date time.Time) ([]*model.MealRecord, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("repository 未初始化")
	}

	// 格式化日期为 YYYY-MM-DD 格式
	dateStr := date.Format("2006-01-02")

	var mealRecords []*model.MealRecord
	err := r.db.WithContext(ctx).Where("user_id = ? AND DATE(date) = ?", userID, dateStr).Order("meal_type").Find(&mealRecords).Error
	if err != nil {
		return nil, err
	}

	return mealRecords, nil
}

// FindByUserIDDateAndType 根据用户ID、日期和餐次类型查找餐次记录
func (r *mealRecordRepository) FindByUserIDDateAndType(ctx context.Context, userID string, date time.Time, mealType model.MealType) (*model.MealRecord, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("repository 未初始化")
	}

	// 格式化日期为 YYYY-MM-DD 格式
	dateStr := date.Format("2006-01-02")

	var mealRecord model.MealRecord
	err := r.db.WithContext(ctx).Where("user_id = ? AND DATE(date) = ? AND meal_type = ?", userID, dateStr, mealType).First(&mealRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("餐次记录不存在")
		}
		return nil, err
	}

	return &mealRecord, nil
}

// Update 更新餐次记录
func (r *mealRecordRepository) Update(ctx context.Context, mealRecord *model.MealRecord) error {
	if r == nil || r.db == nil {
		return errors.New("repository 未初始化")
	}

	result := r.db.WithContext(ctx).Save(mealRecord)
	if result.Error != nil {
		return result.Error
	}

	// 检查是否真的更新了记录
	if result.RowsAffected == 0 {
		return errors.New("没有找到要更新的餐次记录")
	}

	return nil
}

// Delete 删除餐次记录
func (r *mealRecordRepository) Delete(ctx context.Context, id string) error {
	if r == nil || r.db == nil {
		return errors.New("repository 未初始化")
	}

	// 使用软删除（如果模型有 DeletedAt 字段）
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.MealRecord{})
	if result.Error != nil {
		return result.Error
	}

	// 检查是否真的删除了记录
	if result.RowsAffected == 0 {
		return errors.New("没有找到要删除的餐次记录")
	}

	return nil
}
