package service

import (
	"context"
	"errors"

	"github.com/ljk20041215/nutrition-tracker/internal/auth"
	"github.com/ljk20041215/nutrition-tracker/internal/model"
	"github.com/ljk20041215/nutrition-tracker/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req *RegisterRequest) (*model.User, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	GetProfile(ctx context.Context, userID string) (*model.User, error)
	UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname" binding:"required"`
}

func (s *userService) Register(ctx context.Context, req *RegisterRequest) (*model.User, error) {
	// 1. 检查邮箱是否已存在
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 2. 加密密码
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 3. 创建用户
	user := &model.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		Nickname:     req.Nickname,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}

func (s *userService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 1. 查找用户
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("用户不存在或密码错误")
	}

	// 2. 验证密码
	if !checkPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("用户不存在或密码错误")
	}

	// 3. 生成JWT令牌
	token, err := auth.GenerateJWT(user.ID, user.Email, user.Nickname)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	// 4. 返回响应
	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

// UpdateProfileRequest 更新资料请求
type UpdateProfileRequest struct {
	Nickname      string  `json:"nickname"`
	Gender        int     `json:"gender"`
	Age           int     `json:"age"`
	Height        float64 `json:"height"`
	Weight        float64 `json:"weight"`
	ActivityLevel int     `json:"activity_level"`
}

// internal/service/user_service.go 中的相关方法
func (s *userService) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 隐藏敏感信息
	user.PasswordHash = ""
	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) error {
	// 1. 获取现有用户
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// 2. 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	// 更新性别
	if req.Gender != 0 {
		user.Gender = req.Gender
	}
	// 更新年龄 (年龄应为正数)
	if req.Age > 0 {
		user.Age = req.Age
	}
	// 更新身高 (身高应为正数)
	if req.Height > 0 {
		user.Height = req.Height
	}
	// 更新体重 (体重应为正数)
	if req.Weight > 0 {
		user.Weight = req.Weight
	}
	// 更新活动水平 (活动水平应为1-5)
	if req.ActivityLevel >= 1 && req.ActivityLevel <= 5 {
		user.ActivityLevel = req.ActivityLevel
	}

	// 3. 保存更新
	return s.userRepo.Update(ctx, user)
}

// 密码加密和验证辅助函数
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
