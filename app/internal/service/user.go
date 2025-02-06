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
	ListUserByPage(ctx *gin.Context, req *v1.UserQueryRequest) (v1.PageResult[v1.User], error)
	AddUser(ctx *gin.Context, req *v1.AddUserRequest) (uint64, error)
	DeleteUser(ctx *gin.Context, req *v1.DeleteUserRequest) (bool, error)
	UpdateUser(ctx *gin.Context, req *v1.UpdateUserRequest) (bool, error)
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

func (s *userService) UpdateUser(ctx *gin.Context, req *v1.UpdateUserRequest) (bool, error) {
	if req == nil || req.Id == "" {
		return false, v1.ParamsError
	}

	id, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		return false, v1.ParamsError
	}
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return false, err
	}
	if req.UserAccount != nil && *req.UserAccount != "" {
		user.UserAccount = *req.UserAccount
	}
	if req.UserAvatar != nil && *req.UserAvatar != "" {
		user.UserAvatar = req.UserAvatar
	}
	if req.UserName != nil && *req.UserName != "" {
		user.UserName = req.UserName
	}
	if req.UserRole != nil && *req.UserRole != "" {
		user.UserRole = *req.UserRole
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *userService) DeleteUser(ctx *gin.Context, req *v1.DeleteUserRequest) (bool, error) {
	if req.Id <= "0" {
		return false, v1.ParamsError
	}
	id, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		return false, v1.ParamsError
	}
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return false, err
	}
	err = s.userRepo.DeleteById(ctx, user, id)
	if err != nil {
		return false, err
	}
	user.IsDelete = 1
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *userService) AddUser(ctx *gin.Context, req *v1.AddUserRequest) (uint64, error) {
	if *req.UserAccount == "" {
		return 0, v1.ErrIllegalAccount
	}
	user, err := s.userRepo.GetByAccount(ctx, *req.UserAccount)
	if err != nil {
		return 0, v1.ErrInternalServerError
	}
	if user != nil {
		return 0, v1.ErrAccountAlreadyUse
	}
	if len(*req.UserAccount) < 3 || len(*req.UserAccount) > 20 {
		return 0, v1.ErrIllegalAccount
	}

	user = &model.User{
		UserAccount:  *req.UserAccount,
		UserAvatar:   req.UserAvatar,
		UserName:     req.UserName,
		UserProfile:  req.UserProfile,
		UserRole:     *req.UserRole,
		UserPassword: "123456", // 默认密码
	}
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return 0, err
	}
	var u *model.User
	u, err = s.userRepo.GetByAccount(ctx, *req.UserAccount)
	return u.ID, nil
}

func (s *userService) ListUserByPage(ctx *gin.Context, req *v1.UserQueryRequest) (v1.PageResult[v1.User], error) {
	current := req.Current
	size := req.PageSize
	users, err := s.userRepo.GetUser(ctx)
	var user []v1.User
	for _, v := range users {
		var id string
		id = strconv.Itoa(int(v.ID))
		u := v1.User{
			ID:           &id,
			UserAccount:  &v.UserAccount,
			UserPassword: &v.UserPassword,
			UnionID:      v.UnionId,
			MpOpenID:     v.MpOpenId,
			UserName:     v.UserName,
			UserAvatar:   v.UserAvatar,
			UserProfile:  v.UserProfile,
			UserRole:     &v.UserRole,
			EditTime:     &v.EditTime,
			CreateTime:   &v.CreateTime,
			UpdateTime:   &v.UpdateTime,
			IsDelete:     &v.IsDelete,
		}
		user = append(user, u)
	}
	if err != nil {
		return v1.PageResult[v1.User]{}, err
	}
	total := 10
	pages := 10
	return v1.PageResult[v1.User]{
		Records: user,
		Total:   &total,
		Size:    size,
		Current: current,
		Pages:   &pages,
	}, nil
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
