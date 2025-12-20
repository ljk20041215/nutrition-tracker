package repository

import (
	"context"
	"errors"
	"log"

	"github.com/ljk20041215/nutrition-tracker/internal/model"
	"gorm.io/gorm"
)

// FoodRepository 食物仓库接口
type FoodRepository interface {
	Create(ctx context.Context, food *model.Food) error
	FindByID(ctx context.Context, id string) (*model.Food, error)
	FindByName(ctx context.Context, name string) (*model.Food, error)
	FindAll(ctx context.Context) ([]*model.Food, error)
	Update(ctx context.Context, food *model.Food) error
	Delete(ctx context.Context, id string) error
}

// foodRepository 食物仓库实现
type foodRepository struct {
	db *gorm.DB
}

// NewFoodRepository 创建食物仓库实例
func NewFoodRepository(db *gorm.DB) FoodRepository {
	if db == nil {
		log.Fatal("❌ NewFoodRepository: db 参数为 nil")
	}
	return &foodRepository{db: db}
}

// Create 创建食物
func (r *foodRepository) Create(ctx context.Context, food *model.Food) error {
	if r == nil || r.db == nil {
		return errors.New("repository 未初始化")
	}
	return r.db.WithContext(ctx).Create(food).Error
}

// FindByID 根据ID查找食物
func (r *foodRepository) FindByID(ctx context.Context, id string) (*model.Food, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("repository 未初始化")
	}

	var food model.Food
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&food).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("食物不存在")
		}
		return nil, err
	}

	return &food, nil
}

// FindByName 根据名称查找食物
func (r *foodRepository) FindByName(ctx context.Context, name string) (*model.Food, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("repository 未初始化")
	}

	var food model.Food
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&food).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("食物不存在")
		}
		return nil, err
	}

	return &food, nil
}

// FindAll 查找所有食物
func (r *foodRepository) FindAll(ctx context.Context) ([]*model.Food, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("repository 未初始化")
	}

	var foods []*model.Food
	err := r.db.WithContext(ctx).Find(&foods).Error
	if err != nil {
		return nil, err
	}

	return foods, nil
}

// Update 更新食物
func (r *foodRepository) Update(ctx context.Context, food *model.Food) error {
	if r == nil || r.db == nil {
		return errors.New("repository 未初始化")
	}

	result := r.db.WithContext(ctx).Save(food)
	if result.Error != nil {
		return result.Error
	}

	// 检查是否真的更新了记录
	if result.RowsAffected == 0 {
		return errors.New("没有找到要更新的食物")
	}

	return nil
}

// Delete 删除食物
func (r *foodRepository) Delete(ctx context.Context, id string) error {
	if r == nil || r.db == nil {
		return errors.New("repository 未初始化")
	}

	// 使用软删除（如果模型有 DeletedAt 字段）
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Food{})
	if result.Error != nil {
		return result.Error
	}

	// 检查是否真的删除了记录
	if result.RowsAffected == 0 {
		return errors.New("没有找到要删除的食物")
	}

	return nil
}