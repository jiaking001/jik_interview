package repository

import (
	v1 "app/api/v1"
	"app/internal/model"
	"context"
	"errors"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint64) (*model.User, error)
	GetByAccount(ctx context.Context, account string) (*model.User, error)
	GetUser(ctx context.Context, req *v1.UserQueryRequest) ([]*model.User, error)
	DeleteById(ctx context.Context, user *model.User, id uint64) error
	GetCount(ctx context.Context) (int, error)
}

func NewUserRepository(
	r *Repository,
) UserRepository {
	return &userRepository{
		Repository: r,
	}
}

type userRepository struct {
	*Repository
}

func (r *userRepository) GetCount(ctx context.Context) (int, error) {
	var total int64
	var user model.User
	if err := r.DB(ctx).Model(&user).Count(&total).Error; err != nil {
		return 0, err
	}
	return int(total), nil
}

func (r *userRepository) DeleteById(ctx context.Context, user *model.User, id uint64) error {
	if err := r.DB(ctx).Where("id = ?", id).Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetUser(ctx context.Context, req *v1.UserQueryRequest) ([]*model.User, error) {
	var users []*model.User
	var s string
	if req.SortOrder != nil && req.SortField != nil {
		var sortOrder string
		var sortField string
		if *req.SortField == "createTime" {
			sortField = "create_time"
		} else {
			sortField = "update_time"
		}
		if *req.SortOrder == "ascend" {
			sortOrder = "asc"
		} else {
			sortOrder = "desc"
		}
		s = sortField + " " + sortOrder
	}
	var id, userAccount, userName, userRole, userProfile string
	if req.ID != nil {
		id = *req.ID
	}
	if req.UserName != nil {
		userName = *req.UserName
	}
	if req.UserAccount != nil {
		userAccount = *req.UserAccount
	}
	if req.UserRole != nil {
		userRole = *req.UserRole
	}
	if req.UserProfile != nil {
		userProfile = *req.UserProfile
	}
	if err := r.DB(ctx).Where("id LIKE ? AND user_account LIKE ? AND user_name LIKE ? AND user_role LIKE ? AND user_profile LIKE ?",
		"%"+id+"%",
		"%"+userAccount+"%",
		"%"+userName+"%",
		"%"+userRole+"%",
		"%"+userProfile+"%",
	).Order(s).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.DB(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.DB(ctx).Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, userId uint64) (*model.User, error) {
	var user model.User
	if err := r.DB(ctx).Where("id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByAccount(ctx context.Context, account string) (*model.User, error) {
	var user model.User
	if err := r.DB(ctx).Where("user_account = ?", account).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
