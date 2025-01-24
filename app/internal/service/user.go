package service

import (
	v1 "app/api/v1"
	"app/internal/model"
	"app/internal/repository"
	"context"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type UserService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (string, error)
	GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error)
	UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error
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

func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	// check username
	user, err := s.userRepo.GetByAccount(ctx, req.UserAccount)
	if err != nil {
		return v1.ErrInternalServerError
	}
	if user != nil {
		return v1.ErrAccountAlreadyUse
	}

	userPassword := req.UserPassword
	checkPassword := req.CheckPassword
	if userPassword != checkPassword {
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

func (s *userService) Login(ctx context.Context, req *v1.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByAccount(ctx, req.UserAccount)
	if err != nil || user == nil {
		return "", v1.ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(req.UserPassword))
	if err != nil {
		return "", err
	}
	token, err := s.jwt.GenToken(strconv.FormatUint(user.ID, 10), time.Now().Add(time.Hour*24*90))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error) {
	//user, err := s.userRepo.GetByID(ctx, userId)
	//if err != nil {
	//	return nil, err
	//}

	return &v1.GetProfileResponseData{
		//UserId:   user.UserId,
		//Nickname: user.Nickname,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error {
	//user, err := s.userRepo.GetByID(ctx, userId)
	//if err != nil {
	//	return err
	//}
	//
	//user.Email = req.Email
	//user.Nickname = req.Nickname
	//
	//if err = s.userRepo.Update(ctx, user); err != nil {
	//	return err
	//}

	return nil
}
