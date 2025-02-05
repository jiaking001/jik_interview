package service

import (
	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/repository"
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type UserService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (string, *model.User, error)
	GetLoginUser(ctx *gin.Context) error
}

func NewUserService(
	service *Service,
	userRepo repository.UserRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
		Service:  service,
	}
}

type userService struct {
	userRepo repository.UserRepository
	*Service
}

// GetLoginUser 获取当前登录用户
func (s *userService) GetLoginUser(ctx *gin.Context) error {
	// 判断是否已登录
	return nil
}

// Register 用户注册
func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	// check username
	user, err := s.userRepo.GetByAccount(ctx, req.UserAccount)
	if err != nil {
		return v1.ErrInternalServerError
	}
	if user != nil {
		return v1.ErrAccountAlreadyUse
	}
	if len(req.UserAccount) < 3 || len(req.UserAccount) > 20 {
		return v1.ErrIllegalAccount
	}
	if len(req.UserPassword) < 6 || len(req.UserPassword) > 60 {
		return v1.ErrIllegalPassword
	}

	if req.UserPassword != req.CheckPassword {
		return v1.ErrInconsistentPasswords
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user = &model.User{
		UserAccount:  req.UserAccount,
		UserPassword: string(hashedPassword),
	}
	// Transaction demo
	err = s.tm.Transaction(ctx, func(ctx context.Context) error {
		// Create a user
		if err = s.userRepo.Create(ctx, user); err != nil {
			return err
		}
		return nil
	})
	return err
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, req *v1.LoginRequest) (string, *model.User, error) {
	user, err := s.userRepo.GetByAccount(ctx, req.UserAccount)
	if err != nil || user == nil {
		return "", nil, v1.ErrPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(req.UserPassword))
	if err != nil {
		return "", nil, v1.ErrPassword
	}
	token, err := s.jwt.GenToken(strconv.FormatUint(user.ID, 10), time.Now().Add(time.Hour*24*90))
	if err != nil {
		return "", nil, v1.ErrInternalServerError
	}
	return token, user, nil
}
