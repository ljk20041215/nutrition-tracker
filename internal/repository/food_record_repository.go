package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ljk20041215/nutrition-tracker/internal/model"
	"gorm.io/gorm"
)

// FoodRecordRepository 食物记录仓库接口
type FoodRecordRepository interface {
	Create(ctx context.Context, foodRecord *model.FoodRecord) error
	FindByID(ctx context.Context, id string) (*model.FoodRecord, error)
	FindByMealRecordID(ctx context.Context, mealRecordID string) ([]*model.FoodRecord, error)
	FindByUserIDAndDate(ctx context.Context, userID string, date time.Time) ([]*model.FoodRecord, error)
	Update(ctx context.Context, foodRecord *model.FoodRecord) error
	Delete(ctx context.Context, id string) error
	DeleteByMealRecordID(ctx context.Context, mealRecordID string) error
}

// foodRecordRepository 食物记录仓库实现
type foodRecordRepository struct {
	db *gorm.DB
}

// NewFoodRecordRepository 创建食物记录仓库实例
func NewFoodRecordRepository(db *gorm.DB) FoodRecordRepository {
	if db == nil {
		log.Fatal("❌ NewFoodRecordRepository: db 参数为 nil")
	}
	return &foodRecordRepository{db: db}
}

// Create 创建食物记录
func (r *foodRecordRepository) Create(ctx context.Context, foodRecord *model.FoodRecord) error {
	if r == nil || r.db == nil {
		return errors.New("repository 未初始化")
	}
	return r.db.WithContext(ctx).Create(foodRecord).Error
}

// FindByID 根据ID查找食物记录
func (r *foodRecordRepository) FindByID(ctx context.Context, id string) (*model.FoodRecord, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("repository 未初始化")
	}

	var foodRecord model.FoodRecord
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&foodRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("食物记录不存在")
		}
		return nil, err
	}

	return &foodRecord, nil
}

// FindByMealRecordID 根据餐次记录ID查找所有食物记录
func (r *foodRecordRepository) FindByMealRecordID(ctx context.Context, mealRecordID string) ([]*model.FoodRecord, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("repository 未初始化")
	}

	var foodRecords []*model.FoodRecord
	err := r.db.WithContext(ctx).Where("meal_record_id = ?", mealRecordID).Find(&foodRecords).Error
	if err != nil {
		return nil, err
	}

	return foodRecords, nil
}

// FindByUserIDAndDate 根据用户ID和日期查找所有食物记录
func (r *foodRecordRepository) FindByUserIDAndDate(ctx context.Context, userID string, date time.Time) ([]*model.FoodRecord, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("repository 未初始化")
	}

	// 格式化日期为 YYYY-MM-DD 格式
	dateStr := date.Format("2006-01-02")

	var foodRecords []*model.FoodRecord
	err := r.db.WithContext(ctx).Joins("JOIN meal_records ON meal_records.id = food_records.meal_record_id").Where("meal_records.user_id = ? AND DATE(meal_records.date) = ?", userID, dateStr).Find(&foodRecords).Error
	if err != nil {
		return nil, err
	}

	return foodRecords, nil
}

// Update 更新食物记录
func (r *foodRecordRepository) Update(ctx context.Context, foodRecord *model.FoodRecord) error {
	if r == nil || r.db == nil {
		return errors.New("repository 未初始化")
	}

	result := r.db.WithContext(ctx).Save(foodRecord)
	if result.Error != nil {
		return result.Error
	}

	// 检查是否真的更新了记录
	if result.RowsAffected == 0 {
		return errors.New("没有找到要更新的食物记录")
	}

	return nil
}

// Delete 删除食物记录
func (r *foodRecordRepository) Delete(ctx context.Context, id string) error {
	if r == nil || r.db == nil {
		return errors.New("repository 未初始化")
	}

	// 使用软删除（如果模型有 DeletedAt 字段）
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.FoodRecord{})
	if result.Error != nil {
		return result.Error
	}

	// 检查是否真的删除了记录
	if result.RowsAffected == 0 {
		return errors.New("没有找到要删除的食物记录")
	}

	return nil
}

// DeleteByMealRecordID 根据餐次记录ID删除所有食物记录
func (r *foodRecordRepository) DeleteByMealRecordID(ctx context.Context, mealRecordID string) error {
	if r == nil || r.db == nil {
		return errors.New("repository 未初始化")
	}

	// 使用软删除（如果模型有 DeletedAt 字段）
	result := r.db.WithContext(ctx).Where("meal_record_id = ?", mealRecordID).Delete(&model.FoodRecord{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
