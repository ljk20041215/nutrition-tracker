package repository

import (
	"context"
	"errors"
	"log"

	"github.com/ljk20041215/nutrition-tracker/internal/model"
	"gorm.io/gorm"
)

type NutritionGoalRepository interface {
	Create(ctx context.Context, goal *model.NutritionGoal) error
	FindByUserID(ctx context.Context, userID string) (*model.NutritionGoal, error)
	Update(ctx context.Context, goal *model.NutritionGoal) error
	Delete(ctx context.Context, id string) error
}

type nutritionGoalRepository struct {
	db *gorm.DB
}

func NewNutritionGoalRepository(db *gorm.DB) NutritionGoalRepository {
	if db == nil {
		log.Fatal("❌ NewNutritionGoalRepository: db 参数为 nil")
	}
	return &nutritionGoalRepository{db: db}
}

func (r *nutritionGoalRepository) Create(ctx context.Context, goal *model.NutritionGoal) error {
	if r == nil || r.db == nil {
		return errors.New("repository 未初始化")
	}
	return r.db.WithContext(ctx).Create(goal).Error
}

func (r *nutritionGoalRepository) FindByUserID(ctx context.Context, userID string) (*model.NutritionGoal, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("repository 未初始化")
	}

	var goal model.NutritionGoal
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&goal).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("营养目标不存在")
		}
		return nil, err
	}

	return &goal, nil
}

func (r *nutritionGoalRepository) Update(ctx context.Context, goal *model.NutritionGoal) error {
	// 使用 GORM 的 Save 方法，它会根据 ID 更新所有字段
	result := r.db.WithContext(ctx).Save(goal)
	if result.Error != nil {
		return result.Error
	}

	// 检查是否真的更新了记录
	if result.RowsAffected == 0 {
		return errors.New("没有找到要更新的营养目标")
	}

	return nil
}

func (r *nutritionGoalRepository) Delete(ctx context.Context, id string) error {
	// 使用软删除（如果模型有 DeletedAt 字段）
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.NutritionGoal{})
	if result.Error != nil {
		return result.Error
	}

	// 检查是否真的删除了记录
	if result.RowsAffected == 0 {
		return errors.New("没有找到要删除的营养目标")
	}

	return nil
}
