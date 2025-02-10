package service

import (
	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/repository"
	"app/pkg/constant"
	"context"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type UserService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (string, *model.User, error)
	GetLoginUser(ctx context.Context) error
	ListUserByPage(ctx context.Context, req *v1.UserQueryRequest) (v1.PageResult[v1.User], error)
	AddUser(ctx context.Context, req *v1.AddUserRequest) (uint64, error)
	DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (bool, error)
	UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (bool, error)
	AddUserSignIn(ctx context.Context, id uint64) (bool, error)
	GetUserSignIn(ctx context.Context, id uint64, year int) ([]int, error)
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

func (s *userService) GetUserSignIn(ctx context.Context, id uint64, year int) ([]int, error) {
	key := constant.GetUserSignInRedisKey(strconv.Itoa(year), strconv.FormatUint(id, 10))
	bitset, err := s.userRepo.GetUserSignIn(ctx, key)
	if err != nil {
		return nil, err
	}
	var dayList []int
	offset := 0

	for {
		// 使用 nextSetBit 获取下一个签到的天数
		nextOffset := constant.NextSetBit(bitset, offset)
		if nextOffset == -1 {
			break
		}

		// 将签到的天数添加到结果列表中
		dayList = append(dayList, nextOffset)

		// 更新偏移量，继续查找下一个签到的天数
		offset = nextOffset + 1
	}
	return dayList, nil
}

func (s *userService) AddUserSignIn(ctx context.Context, id uint64) (bool, error) {
	date := time.Now()
	year := date.Year()
	key := constant.GetUserSignInRedisKey(strconv.Itoa(year), strconv.FormatUint(id, 10))
	offset := date.YearDay()
	err := s.userRepo.AddUserSignIn(ctx, key, int64(offset))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *userService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (bool, error) {
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
	if req.UserProfile != nil && *req.UserProfile != "" {
		user.UserProfile = req.UserProfile
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *userService) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (bool, error) {
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

func (s *userService) AddUser(ctx context.Context, req *v1.AddUserRequest) (uint64, error) {
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user = &model.User{
		UserAccount:  *req.UserAccount,
		UserAvatar:   req.UserAvatar,
		UserName:     req.UserName,
		UserProfile:  req.UserProfile,
		UserRole:     *req.UserRole,
		UserPassword: string(hashedPassword), // 默认密码
	}
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return 0, err
	}
	var u *model.User
	u, err = s.userRepo.GetByAccount(ctx, *req.UserAccount)
	if err != nil {
		return 0, err
	}
	return u.ID, nil
}

func (s *userService) ListUserByPage(ctx context.Context, req *v1.UserQueryRequest) (v1.PageResult[v1.User], error) {
	current := req.Current
	size := req.PageSize
	users, total, err := s.userRepo.GetUser(ctx, req)
	if err != nil {
		return v1.PageResult[v1.User]{}, err
	}
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
	pages := total / *size + 1
	return v1.PageResult[v1.User]{
		Records: user,
		Total:   &total,
		Size:    size,
		Current: current,
		Pages:   &pages,
	}, nil
}

// GetLoginUser 获取当前登录用户
func (s *userService) GetLoginUser(ctx context.Context) error {
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
		UserName:     &req.UserAccount,
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
