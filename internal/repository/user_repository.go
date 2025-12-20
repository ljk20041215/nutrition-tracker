package repository

import (
	"context"
	"errors"
	"log"

	"github.com/ljk20041215/nutrition-tracker/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	if db == nil {
		log.Fatal("❌ NewUserRepository: db 参数为 nil")
	}
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if r == nil || r.db == nil {
		return errors.New("repository 未初始化")
	}
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	// 使用 GORM 的 First 方法，按 ID 查找用户
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("repository 未初始化")
	}

	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	// 使用 GORM 的 Save 方法，它会根据 ID 更新所有字段
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}

	// 检查是否真的更新了记录
	if result.RowsAffected == 0 {
		return errors.New("没有找到要更新的用户")
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	// 使用软删除（如果模型有 DeletedAt 字段）
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{})
	if result.Error != nil {
		return result.Error
	}

	// 检查是否真的删除了记录
	if result.RowsAffected == 0 {
		return errors.New("没有找到要删除的用户")
	}

	return nil
}

// 可选：添加其他有用的方法
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
