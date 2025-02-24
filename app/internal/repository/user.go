package repository

import (
	v1 "app/api/v1"
	"app/internal/model"
	"app/pkg/utils"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint64) (*model.User, error)
	GetByAccount(ctx context.Context, account string) (*model.User, error)
	GetUser(ctx context.Context, req *v1.UserQueryRequest) ([]*model.User, int, error)
	DeleteById(ctx context.Context, user *model.User, id uint64) error
	GetCount(ctx context.Context) (int, error)
	AddUserSignIn(ctx context.Context, key string, offset int64) error
	GetUserSignIn(ctx context.Context, key string) ([]byte, error)
	GetTokenByDevice(ctx context.Context, id uint64, deviceType string) (string, error)
	AddTokenByDevice(ctx context.Context, id uint64, deviceType string, token string) error
	DeleteTokenByDevice(ctx context.Context, id uint64, deviceType string) error
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

func (r *userRepository) DeleteTokenByDevice(ctx context.Context, id uint64, deviceType string) error {
	idStr := utils.Uint64TOString(id)
	if err := r.rdb.HDel(ctx, "user_tokens:"+idStr, deviceType).Err(); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) AddTokenByDevice(ctx context.Context, id uint64, deviceType string, token string) error {
	idStr := utils.Uint64TOString(id)
	if err := r.rdb.HSet(ctx, "user_tokens:"+idStr, deviceType, token).Err(); err != nil {
		return err
	}
	// 设置过期时间 一天后过期
	if err := r.rdb.Expire(ctx, "user_tokens:"+idStr, time.Hour*24).Err(); err != nil {
		return err
	}
	return nil
}

// GetTokenByDevice 获取当前用户登录端的Token
func (r *userRepository) GetTokenByDevice(ctx context.Context, id uint64, deviceType string) (string, error) {
	idStr := utils.Uint64TOString(id)
	return r.rdb.HGet(ctx, "user_tokens:"+idStr, deviceType).Result()
}

func (r *userRepository) GetUserSignIn(ctx context.Context, key string) ([]byte, error) {
	bytes, err := r.rdb.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (r *userRepository) AddUserSignIn(ctx context.Context, key string, offset int64) error {
	return r.rdb.SetBit(ctx, key, offset, 1).Err()
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
	tx := r.DB(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("id = ?", id).Delete(user).Error; err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (r *userRepository) GetUser(ctx context.Context, req *v1.UserQueryRequest) ([]*model.User, int, error) {
	var users []*model.User
	var total int64
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
	// 拼接查询字符串
	// query 查询字符串
	var conditions []string
	var params []interface{}
	var query string
	var id, userAccount, userName, userRole, userProfile string
	if req.ID != nil {
		id = *req.ID
		conditions = append(conditions, "id LIKE ?")
		params = append(params, "%"+id+"%")
	}
	if req.UserName != nil {
		userName = *req.UserName
		conditions = append(conditions, "user_name LIKE ?")
		params = append(params, "%"+userName+"%")
	}
	if req.UserAccount != nil {
		userAccount = *req.UserAccount
		conditions = append(conditions, "user_account LIKE ?")
		params = append(params, "%"+userAccount+"%")
	}
	if req.UserRole != nil {
		userRole = *req.UserRole
		conditions = append(conditions, "user_role LIKE ?")
		params = append(params, "%"+userRole+"%")
	}
	if req.UserProfile != nil {
		userProfile = *req.UserProfile
		conditions = append(conditions, "user_profile LIKE ?")
		params = append(params, "%"+userProfile+"%")
	}

	// 构造完整的查询条件
	query = strings.Join(conditions, " AND ")

	if err := r.DB(ctx).Where(query, params...).Order(s).Find(&users).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return users, int(total), nil
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
